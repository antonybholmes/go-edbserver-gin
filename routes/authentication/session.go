package authenticationroutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

const (
	//SESSION_PUBLICID   string = "publicId"
	//SESSION_ROLES      string = "roles"
	SESSION_USER       string = "user"
	SESSION_CREATED_AT string = "createdAt"
	SESSION_EXPIRES_AT string = "expiresAt"
)

const (
	ERROR_CREATING_SESSION string = "error creating session"
)

var SESSION_OPT_ZERO sessions.Options

//var SESSION_OPT_24H *sessions.Options
//var SESSION_OPT_30_DAYS *sessions.Options
//var SESSION_OPT_7_DAYS *sessions.Options

func init() {

	// HttpOnly and Secure are disabled so we can use them
	// cross domain for testing
	// http only false to allow js to delete etc on the client side

	// For sessions that should end when browser closes
	SESSION_OPT_ZERO = sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	// SESSION_OPT_24H = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_DAY_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// SESSION_OPT_30_DAYS = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_30_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// SESSION_OPT_7_DAYS = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_7_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// }
}

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
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	return &SessionRoutes{sessionOptions: options}
}

func (sr *SessionRoutes) SessionUsernamePasswordSignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.LoginBodyReq.Password == "" {
		PasswordlessSigninEmailRoute(c, validator)
	}

	user := validator.LoginBodyReq.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	if err != nil {
		routes.UserDoesNotExistResp(c)
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		routes.EmailNotVerifiedReq(c)
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		routes.AuthErrorReq(c, "could not get user roles")
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		routes.UserNotAllowedToSignIn(c)
	}

	err = authUser.CheckPasswordsMatch(validator.LoginBodyReq.Password)

	if err != nil {
		c.Error(err)
		return
	}

	sess := sessions.Default(c) //Key(consts.SESSION_NAME)

	if err != nil {
		c.Error(err)
		return
	}

	// set session options
	if validator.LoginBodyReq.StaySignedIn {
		sess.Options(sr.sessionOptions)
	} else {
		sess.Options(SESSION_OPT_ZERO)
	}

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	sess.Save() //c.Request(), c.Response())

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
		routes.UserDoesNotExistResp(c)
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		routes.EmailNotVerifiedReq(c)
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		routes.AuthErrorReq(c, "could not get user roles")
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		routes.UserNotAllowedToSignIn(c)
	}

	err = sr.initSession(c, authUser) //, roleClaim)

	if err != nil {
		routes.ErrorResp(c, ERROR_CREATING_SESSION)
		return
	}

	c.String(http.StatusOK, "Session set")

	// resp, err := readSession(c)

	// if err != nil {
	// 	routes.ErrorReq(ERROR_CREATING_SESSION)
	// }

	// routes.MakeDataResp(c, "", resp)
}

func (sr *SessionRoutes) SessionSignInUsingAuth0Route(c *gin.Context) {
	user, ok := c.Get("user")

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		routes.TokenErrorReq(c)
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

	authUser, err := userdbcache.CreateUserFromAuth0(tokenClaims.Name, email)

	if err != nil {
		c.Error(err)
		return
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		routes.ErrorResp(c, "user roles not found")
	}

	roleClaim := auth.MakeClaim(roles)

	//log.Debug().Msgf("user %v", authUser)

	if !auth.CanSignin(roleClaim) {
		routes.UserNotAllowedToSignIn(c)
	}

	err = sr.initSession(c, authUser) // roleClaim)

	if err != nil {
		c.Error(err)
		return
	}

	UserSignedInResp(c)
}

type SessionInfo struct {
	AuthUser  *auth.AuthUser `json:"user"`
	IsValid   bool           `json:"valid"`
	CreatedAt string         `json:"createdAt"`
	ExpiresAt string         `json:"expiresAt"`
}

// initialize a session with default age and ids
func (sr *SessionRoutes) initSession(c *gin.Context, authUser *auth.AuthUser) error {

	userData, err := json.Marshal(authUser)

	if err != nil {
		return err
	}

	sess := sessions.Default(c) // .Get(consts.SESSION_NAME, c)

	if err != nil {
		return fmt.Errorf("%s", ERROR_CREATING_SESSION)
	}

	// set session options
	sess.Options(sr.sessionOptions)

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roles //auth.MakeClaim(authUser.Roles)
	sess.Set(SESSION_USER, string(userData))

	now := time.Now().UTC()
	sess.Set(SESSION_CREATED_AT, now.Format(time.RFC3339))
	sess.Set(SESSION_EXPIRES_AT, now.Add(time.Duration(sr.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
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

func ReadSessionInfo(c *gin.Context) (*SessionInfo, error) {
	sess := sessions.Default(c) //.Get(consts.SESSION_NAME, c)

	userData, _ := sess.Get(SESSION_USER).(string)

	var user auth.AuthUser

	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	//publicId, _ := sess.Values[SESSION_PUBLICID].(string)
	//roles, _ := sess.Values[SESSION_ROLES].(string)
	createdAt, _ := sess.Get(SESSION_CREATED_AT).(string)
	expires, _ := sess.Get(SESSION_EXPIRES_AT).(string)
	//isValid := publicId != ""

	return &SessionInfo{AuthUser: &user,
			CreatedAt: createdAt,
			ExpiresAt: expires},
		nil
}

// Read the user session. Can also be used to determin if session is valid
func (sr *SessionRoutes) SessionInfoRoute(c *gin.Context) {
	sessionInfo, err := ReadSessionInfo(c)

	if err != nil {
		c.Error(err)
		return
	}

	routes.MakeDataResp(c, "", sessionInfo)
}

func (sr *SessionRoutes) SessionRenewRoute(c *gin.Context) {
	user, ok := c.Get("authUser")

	if !ok {
		routes.ErrorResp(c, "no auth user")
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

	if err != nil {
		c.Error(fmt.Errorf("%s", ERROR_CREATING_SESSION))
		return
	}

	userData, err := json.Marshal(authUser)

	if err != nil {
		c.Error(err)
		return
	}

	log.Debug().Msgf("saving %s", string(userData))

	sess.Set(SESSION_USER, string(userData))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sr *SessionRoutes) SessionPasswordlessValidateSignInRoute(c *gin.Context) {

	NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			routes.WrongTokentTypeReq(c)
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
			routes.UserNotAllowedToSignIn(c)
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
	sess.Set(SESSION_USER, "")
	//sess.Values[SESSION_ROLES] = ""
	sess.Set(SESSION_CREATED_AT, "")
	sess.Set(SESSION_EXPIRES_AT, "")
	sess.Options(SESSION_OPT_ZERO)

	sess.Save() //c.Request(), c.Response())

	routes.MakeOkResp(c, "user has been signed out")
}

func NewAccessTokenFromSessionRoute(c *gin.Context) {
	// sess, _ := session.Get(consts.SESSION_NAME, c)

	// userData, ok := sess.Values["user"].(string)

	// if !ok {
	// 	routes.ErrorReq(fmt.Errorf("malformed user info"))
	// }

	// var user auth.AuthUser
	// if err := json.Unmarshal([]byte(userData), &user); err != nil {
	// 	routes.ErrorReq(err)
	// }

	user, ok := c.Get("authUser")

	if !ok {
		routes.ErrorResp(c, "no auth user")
		return
	}

	authUser := user.(*auth.AuthUser)
	//publicId, _ := sess.Values[SESSION_PUBLICID].(string)
	//r//oles, _ := sess.Values[SESSION_ROLES].(string)

	// if publicId == "" {
	// 	routes.ErrorReq(fmt.Errorf("public id cannot be empty"))
	// }

	// generate a new token from what is stored in the sesssion
	t, err := tokengen.AccessToken(c, authUser.Uuid, auth.MakeClaim(authUser.Roles))

	if err != nil {
		routes.TokenErrorReq(c)
	}

	routes.MakeDataResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func UserFromSessionRoute(c *gin.Context) {
	user, ok := c.Get("authUser")

	if !ok {
		routes.ErrorResp(c, "no auth user")
		return
	}

	routes.MakeDataResp(c, "", user)
}
