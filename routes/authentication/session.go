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

const MAX_AGE_ONE_YEAR_SECS = 31536000 // 60 * 60 * 24 * 365

type SessionRoutes struct {
	sessionOptions sessions.Options
	OTPRoutes      *OTPRoutes
}

func NewSessionRoutes(otpRoutes *OTPRoutes) *SessionRoutes {
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

	return &SessionRoutes{sessionOptions: options, OTPRoutes: otpRoutes}
}

// initialize a session with default age and ids
func (sessionRoutes *SessionRoutes) initSession(c *gin.Context, authUser *auth.AuthUser) error {

	userData, err := json.Marshal(authUser)

	if err != nil {
		return err
	}

	//csrfToken, err := web.GenerateCSRFToken()

	// if err != nil {
	// 	return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	// }

	// log.Debug().Msgf("GenerateCSRFToken %s", csrfToken)

	sess := sessions.Default(c) // .Get(consts.SESSION_NAME, c)

	// set session options
	sess.Options(sessionRoutes.sessionOptions)

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roles //auth.MakeClaim(authUser.Roles)
	sess.Set(web.SESSION_USER, string(userData))
	//sess.Set(web.SESSION_CSRF_TOKEN, csrfToken)

	now := time.Now().UTC()
	sess.Set(web.SESSION_CREATED_AT, now.Format(time.RFC3339))
	sess.Set(web.SESSION_EXPIRES_AT, now.Add(time.Duration(sessionRoutes.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return err
	}

	// csrfToken, err := middleware.CreateCSRFTokenCookie(c)

	// if err != nil {
	// 	//c.Error(err)
	// 	return "", err
	// }

	// Also send it to the client in a readable cookie
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:  consts.CSRF_COOKIE_NAME,
	// 	Value: csrfToken,
	// 	Path:  "/",
	// 	//Domain:   "ed.site.com", // or leave empty if called via ed.site.com
	// 	MaxAge:   MAX_AGE_ONE_YEAR_SECS, // 0 means until browser closes
	// 	Secure:   true,
	// 	HttpOnly: false, // must be readable from JS!
	// 	SameSite: http.SameSiteNoneMode,
	// })

	return nil
}

// create empty session for testing
// func (sr *SessionRoutes) InitSessionRoute(c *gin.Context) {

// 	  err := sr.initSession(c, &auth.AuthUser{})

// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	MakeCsrfTokenResp(c, csrfToken)

// }

func (sessionRoutes *SessionRoutes) SessionUsernamePasswordSignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).LoadAuthUserFromUsername().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.UserBodyReq.Password == "" {
		PasswordlessSigninEmailRoute(c, validator)
		return
	}

	user := validator.UserBodyReq.Username

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
		web.ForbiddenResp(c, "could not get user roles")
		return
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
		return
	}

	err = authUser.CheckPasswordsMatch(validator.UserBodyReq.Password)

	if err != nil {
		c.Error(err)
		return
	}

	userData, err := json.Marshal(authUser)

	if err != nil {
		c.Error(err)
		return
	}

	sess := sessions.Default(c) //Key(consts.SESSION_NAME)

	// set session options
	if validator.UserBodyReq.StaySignedIn {
		sess.Options(sessionRoutes.sessionOptions)
	} else {
		sess.Options(middleware.SESSION_OPT_ZERO)
	}

	sess.Set(web.SESSION_USER, string(userData))

	now := time.Now().UTC()
	sess.Set(web.SESSION_CREATED_AT, now.Format(time.RFC3339))
	sess.Set(web.SESSION_EXPIRES_AT, now.Add(time.Duration(sessionRoutes.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeCsrfTokenResp(c)
	//return c.NoContent(http.StatusOK)
}

func (sessionRoutes *SessionRoutes) SessionApiKeySignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.FindUserByApiKey(validator.UserBodyReq.ApiKey)

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
		web.ForbiddenResp(c, "could not get user roles")
		return
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
		return
	}

	err = sessionRoutes.initSession(c, authUser) //, roleClaim)

	if err != nil {
		web.BadReqResp(c, middleware.ERROR_CREATING_SESSION)
		return
	}

	//MakeCsrfTokenResp(c, token)

	web.MakeOkResp(c, "user has been signed in")

	// resp, err := readSession(c)

	// if err != nil {
	// 	web.ErrorReq(ERROR_CREATING_SESSION)
	// }

	// web.MakeDataResp(c, "", resp)
}

func (sessionRoutes *SessionRoutes) SessionSignInUsingAuth0Route(c *gin.Context) {
	user, ok := c.Get(web.SESSION_USER)

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

	sessionRoutes.sessionSignInUsingOAuth2(c, authUser)
}

func (sessionRoutes *SessionRoutes) SessionSignInUsingClerkRoute(c *gin.Context) {
	user, ok := c.Get(web.SESSION_USER)

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

	sessionRoutes.sessionSignInUsingOAuth2(c, authUser)
}

func (sessionRoutes *SessionRoutes) SessionSignInUsingSupabaseRoute(c *gin.Context) {
	user, ok := c.Get(web.SESSION_USER)

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

	sessionRoutes.sessionSignInUsingOAuth2(c, authUser)
}

func (sessionRoutes *SessionRoutes) SessionEmailOTPRoute(c *gin.Context) {
	sessionRoutes.OTPRoutes.EmailOTPRoute(c)
}

func (sessionRoutes *SessionRoutes) SessionSignInUsingEmailAndOTPRoute(c *gin.Context) {
	validator, err := NewValidator(c).CheckEmailIsWellFormed().Ok()

	if err != nil {
		web.BaseBadReqResp(c, err)
		return
	}

	username := validator.Address.Address

	otpValid, err := sessionRoutes.OTPRoutes.OTP.validateOTP(username, validator.UserBodyReq.OTP)

	if !otpValid || err != nil {
		web.BadReqResp(c, "invalid one time passcode")
		return
	}

	authUser, err := userdbcache.CreateUserFromOAuth2(username, validator.Address)

	if err != nil {
		c.Error(err)
		return
	}

	sessionRoutes.sessionSignInUsingOAuth2(c, authUser)
}

func (sessionRoutes *SessionRoutes) sessionSignInUsingOAuth2(c *gin.Context, authUser *auth.AuthUser) {

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		web.BadReqResp(c, "user roles not found")
	}

	roleClaim := auth.MakeClaim(roles)

	//log.Debug().Msgf("user %v", authUser)

	if !auth.CanSignin(roleClaim) {
		web.UserNotAllowedToSignInErrorResp(c)
	}

	err = sessionRoutes.initSession(c, authUser) // roleClaim)

	if err != nil {
		web.BaseUnauthorizedResp(c, err)
		return
	}

	//log.Debug().Msgf("token %s", token)

	web.MakeCsrfTokenResp(c)

	//web.MakeOkResp(c, "user has been signed in")
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sessionRoutes *SessionRoutes) SessionPasswordlessValidateSignInRoute(c *gin.Context) {

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

		err = sessionRoutes.initSession(c, authUser) //, roleClaim)

		if err != nil {
			web.BaseInternalErrorResp(c, err)
			return
		}

		web.MakeOkResp(c, "user has signed in") //MakeCsrfTokenResp(c, token)
	})
}

func SessionSignOutRoute(c *gin.Context) {
	sess := sessions.Default(c) //.Get(consts.SESSION_NAME, c)

	//log.Debug().Msgf("invalidate session")

	// invalidate by time
	//sess.Set(web.SESSION_USER, "")
	//sess.Values[SESSION_ROLES] = ""
	//sess.Set(web.SESSION_CREATED_AT, "")
	//sess.Set(web.SESSION_EXPIRES_AT, "")
	sess.Clear()
	sess.Options(middleware.SESSION_OPT_CLEAR) //.SESSION_OPT_ZERO)
	sess.Save()                                //c.Request(), c.Response())

	// Also send it to the client in a readable cookie
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:     web.CSRF_COOKIE_NAME,
	// 	Value:    "",
	// 	Path:     "/",
	// 	MaxAge:   -1, // 0 means until browser closes
	// 	Secure:   true,
	// 	HttpOnly: false, // must be readable from JS!
	// 	SameSite: http.SameSiteNoneMode,
	// })

	web.MakeOkResp(c, "user has been signed out")
}

// Read the user session. Can also be used to determin if session is valid
func (sessionRoutes *SessionRoutes) SessionInfoRoute(c *gin.Context) {
	session := sessions.Default(c)

	sessionInfo, err := middleware.ReadSessionInfo(c, session)

	if err != nil {
		web.UnauthorizedResp(c, "session not found or expired")
		return
	}

	web.MakeDataResp(c, "", sessionInfo)
}

func (sessionRoutes *SessionRoutes) SessionCsrfTokenRoute(c *gin.Context) {
	web.MakeCsrfTokenResp(c)
}

func (sessionRoutes *SessionRoutes) SessionRefreshRoute(c *gin.Context) {
	user, ok := c.Get(web.SESSION_USER)

	if !ok {
		web.UnauthorizedResp(c, "no auth user")
		return
	}

	// refresh user
	authUser, err := userdbcache.FindUserById(user.(*auth.AuthUser).Id)

	if err != nil {
		web.UnauthorizedResp(c, "user not found")
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

	sess.Set(web.SESSION_USER, string(userData))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		web.InternalErrorResp(c, "error saving session")
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
	user, _ := c.Get(web.SESSION_USER)

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
		web.BaseInternalErrorResp(c, err)
		return
	}

	web.MakeDataResp(c, "", &web.TokenResp{Token: token})

}

func UserFromSessionRoute(c *gin.Context) {
	user, ok := c.Get(web.SESSION_USER)

	if !ok {
		web.BadReqResp(c, "no auth user")
		return
	}

	web.MakeDataResp(c, "", user)
}
