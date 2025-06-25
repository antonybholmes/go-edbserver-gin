package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SessionRoutes struct {
	sessionOptions sessions.Options
}

func NewSessionRoutes() *SessionRoutes {
	maxAge := auth.MAX_AGE_7_DAYS_SECS

	t := os.Getenv("SESSION_TTL_HOURS")

	if t != "" {
		v, err := strconv.ParseUint(t, 10, 32)

		if err == nil {
			maxAge = int((time.Duration(v) * time.Hour).Seconds())
		}
	}

	options := sessions.Options{
		Path: "/",
		//Domain:   consts.APP_DOMAIN,
		MaxAge:   maxAge,
		HttpOnly: true, //false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	return &SessionRoutes{sessionOptions: options}
}

// initialize a session with default age and ids
func (sr *SessionRoutes) initSession(c *gin.Context, authUser *auth.AuthUser) error {

	userData, err := json.Marshal(authUser)

	if err != nil {
		return err
	}

	sess := sessions.Default(c) // .Get(consts.SESSION_NAME, c)

	// set session options
	sess.Options(sr.sessionOptions)

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roles //auth.MakeClaim(authUser.Roles)
	sess.Set(middleware.SESSION_USER, string(userData))

	now := time.Now().UTC()
	sess.Set(middleware.SESSION_CREATED_AT, now.Format(time.RFC3339))
	sess.Set(middleware.SESSION_EXPIRES_AT, now.Add(time.Duration(sr.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return err
	}

	return nil
}

// create empty session for testing
func (sr *SessionRoutes) InitSessionRoute(c *gin.Context) {

	err := sr.initSession(c, &auth.AuthUser{})

	if err != nil {
		c.Error(err)
		return
	}

}

func (sr *SessionRoutes) SessionUsernamePasswordSignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.LoginBodyReq.Password == "" {
		PasswordlessSigninEmailRoute(c, validator)
		return
	}

	user := validator.LoginBodyReq.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	if err != nil {
		web.UserDoesNotExistResp(c)
		return
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		web.EmailNotVerifiedReq(c)
		return
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		web.AuthErrorResp(c, "could not get user roles")
		return
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
		return
	}

	err = authUser.CheckPasswordsMatch(validator.LoginBodyReq.Password)

	if err != nil {
		c.Error(err)
		return
	}

	sess := sessions.Default(c) //Key(consts.SESSION_NAME)

	// set session options
	if validator.LoginBodyReq.StaySignedIn {
		sess.Options(sr.sessionOptions)
	} else {
		sess.Options(middleware.SESSION_OPT_ZERO)
	}

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

	UserSignedInResp(c)
	//return c.NoContent(http.StatusOK)
}

func (sr *SessionRoutes) SessionApiKeySignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.FindUserByApiKey(validator.LoginBodyReq.ApiKey)

	if err != nil {
		web.UserDoesNotExistResp(c)
		return
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		web.EmailNotVerifiedReq(c)
		return
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		web.AuthErrorResp(c, "could not get user roles")
		return
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
		return
	}

	err = sr.initSession(c, authUser) //, roleClaim)

	if err != nil {
		web.ErrorResp(c, middleware.ERROR_CREATING_SESSION)
		return
	}

	c.String(http.StatusOK, "session created")

	// resp, err := readSession(c)

	// if err != nil {
	// 	web.ErrorReq(ERROR_CREATING_SESSION)
	// }

	// web.MakeDataResp(c, "", resp)
}

func (sr *SessionRoutes) SessionSignInUsingAuth0Route(c *gin.Context) {
	user, ok := c.Get(middleware.SESSION_USER)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		web.TokenErrorResp(c)

		return
	}

	tokenClaims := user.(*auth.Auth0TokenClaims)

	//myClaims := user.Claims.(*auth.TokenClaims) //hmm.Claims.(*TokenClaims)

	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*TokenClaims)

	//log.Debug().Msgf("auth0 claims %v", tokenClaims)
	//log.Debug().Msgf("auth0 claims %v", tokenClaims.Email)

	email, err := mail.ParseAddress(tokenClaims.Email)

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.CreateUserFromOAuth2(tokenClaims.Name, email)

	if err != nil {
		c.Error(err)
		return
	}

	sr.sessionSignInUsingOAuth2(c, authUser)
}

func (sr *SessionRoutes) SessionSignInUsingClerkRoute(c *gin.Context) {
	user, ok := c.Get(middleware.SESSION_USER)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		web.TokenErrorResp(c)

		return
	}

	tokenClaims := user.(*auth.ClerkTokenClaims)

	email, err := mail.ParseAddress(tokenClaims.Email)

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.CreateUserFromOAuth2(tokenClaims.Name, email)

	if err != nil {
		c.Error(err)
		return
	}

	sr.sessionSignInUsingOAuth2(c, authUser)
}

func (sr *SessionRoutes) SessionSignInUsingSupabaseRoute(c *gin.Context) {
	user, ok := c.Get(middleware.SESSION_USER)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		web.TokenErrorResp(c)

		return
	}

	tokenClaims := user.(*auth.SupabaseTokenClaims)

	email, err := mail.ParseAddress(tokenClaims.Email)

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.CreateUserFromOAuth2(tokenClaims.Email, email)

	if err != nil {
		c.Error(err)
		return
	}

	sr.sessionSignInUsingOAuth2(c, authUser)
}

func (sr *SessionRoutes) sessionSignInUsingOAuth2(c *gin.Context, authUser *auth.AuthUser) {

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		web.ErrorResp(c, "user roles not found")
	}

	roleClaim := auth.MakeClaim(roles)

	//log.Debug().Msgf("user %v", authUser)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
	}

	err = sr.initSession(c, authUser) // roleClaim)

	if err != nil {
		c.Error(err)
		return
	}

	UserSignedInResp(c)
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sr *SessionRoutes) SessionPasswordlessValidateSignInRoute(c *gin.Context) {

	NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			web.WrongTokentTypeReq(c)
			return
		}

		authUser := validator.AuthUser

		roles, err := userdbcache.UserRoleList(authUser)

		if err != nil {
			c.Error(err)
			return
		}

		roleClaim := auth.MakeClaim(roles)

		//log.Debug().Msgf("user %v", authUser)

		if !auth.CanSignin(roleClaim) {
			web.UserNotAllowedToSignInErrorResp(c)
			return
		}

		err = sr.initSession(c, authUser) //, roleClaim)

		if err != nil {
			c.Error(err)
			return
		}

		UserSignedInResp(c)
	})
}

func SessionSignOutRoute(c *gin.Context) {
	sess := sessions.Default(c) //.Get(consts.SESSION_NAME, c)

	log.Debug().Msgf("invalidate session")

	// invalidate by time
	sess.Set(middleware.SESSION_USER, "")
	//sess.Values[SESSION_ROLES] = ""
	sess.Set(middleware.SESSION_CREATED_AT, "")
	sess.Set(middleware.SESSION_EXPIRES_AT, "")
	sess.Options(middleware.SESSION_OPT_ZERO)

	sess.Save() //c.Request(), c.Response())

	web.MakeOkResp(c, "user has been signed out")
}

// Read the user session. Can also be used to determin if session is valid
func (sr *SessionRoutes) SessionInfoRoute(c *gin.Context) {
	sessionInfo, err := middleware.ReadSessionInfo(c)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", sessionInfo)
}

func (sr *SessionRoutes) SessionRefreshRoute(c *gin.Context) {
	user, ok := c.Get(middleware.SESSION_USER)

	if !ok {
		web.ErrorResp(c, "no auth user")
		return
	}

	// refresh user
	authUser, err := userdbcache.FindUserById(user.(*auth.AuthUser).Id)

	if err != nil {
		c.Error(err)
		return
	}

	//
	// For the moment just update the user info

	//err = sr.initSession(c, authUser)

	sess := sessions.Default(c) // .Get(consts.SESSION_NAME, c)

	userData, err := json.Marshal(authUser)

	if err != nil {
		c.Error(err)
		return
	}

	log.Debug().Msgf("saving %s", string(userData))

	sess.Set(middleware.SESSION_USER, string(userData))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

}

// func NewAccessTokenFromSessionRoute(c *gin.Context) {

// 	user, _ := c.Get(middleware.SESSION_USER)

// 	authUser := user.(*auth.AuthUser)

// 	t, err := tokengen.AccessToken(c, authUser.PublicId, auth.MakeClaim(authUser.Roles))

// 	if err != nil {
// 		web.TokenErrorResp(c)
// 		return
// 	}

// 	web.MakeDataResp(c, "", &web.AccessTokenResp{AccessToken: t})
// }

func CreateTokenFromSessionRoute(c *gin.Context) {

	tokenType := c.Param("type")

	// user must exist or middleware would have failed
	user, _ := c.Get(middleware.SESSION_USER)

	authUser := user.(*auth.AuthUser)

	var token string
	var err error

	switch tokenType {
	case "access":
		// Generate encoded token and send it as response.
		token, err = tokengen.AccessToken(c, authUser.PublicId, auth.MakeClaim(authUser.Roles))

		if err != nil {
			err = fmt.Errorf("error creating access token: %w", err)
		}
	case "update":
		// Generate encoded token and send it as response.
		token, err = tokengen.UpdateToken(c, authUser.PublicId, auth.MakeClaim(authUser.Roles))

		if err != nil {
			err = fmt.Errorf("error creating update token: %w", err)
		}
	default:
		err = fmt.Errorf("unknown token type")
	}

	if err != nil {
		web.BaseErrorResp(c, err)
		return
	}

	web.MakeDataResp(c, "", &web.TokenResp{Token: token})

}

func UserFromSessionRoute(c *gin.Context) {
	user, ok := c.Get(middleware.SESSION_USER)

	if !ok {
		web.ErrorResp(c, "no auth user")
		return
	}

	web.MakeDataResp(c, "", user)
}
