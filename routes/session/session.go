package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/userdb"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"
	csrfmiddleware "github.com/antonybholmes/go-web/middleware/csrf"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const MaxAgeOneYearSecs = 31536000 // 60 * 60 * 24 * 365

var (
	ErrNoSessionUser  = errors.New("no auth user")
	ErrSavingSession  = errors.New("error saving session")
	ErrSessionExpired = errors.New("session not found or expired")
)

type SessionRoutes struct {
	sessionOptions sessions.Options
	OTPRoutes      *authentication.OTPRoutes
}

func NewSessionRoutes(otpRoutes *authentication.OTPRoutes) *SessionRoutes {
	maxAge := auth.MaxAge7DaysSecs

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
	sess.Set(web.SessionUser, string(userData))
	//sess.Set(web.SESSION_CSRF_TOKEN, csrfToken)

	now := time.Now().UTC()
	sess.Set(web.SessionCreatedAt, now.Format(time.RFC3339))
	sess.Set(web.SessionExpiresAt, now.Add(time.Duration(sessionRoutes.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

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
	validator, err := authentication.NewValidator(c).LoadAuthUserFromUsername().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.UserBodyReq.Password == "" {
		authentication.PasswordlessSignInEmailRoute(c, validator)
		return
	}

	user := validator.UserBodyReq.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	if err != nil {
		web.UserDoesNotExistResp(c)
		return
	}

	if authUser.EmailVerifiedAt == userdb.EmailNotVerifiedDate {
		web.EmailNotVerifiedReq(c)
		return
	}

	// roles, err := userdbcache.UserRoleSet(authUser)

	// if err != nil {
	// 	web.ForbiddenResp(c, authentication.ErrUserRoles)
	// 	return
	// }

	//roleClaim := auth.MakeClaim(roles)

	if !auth.UserHasWebLoginInRole(authUser) {
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
		sess.Options(middleware.SessionOptsZero)
	}

	sess.Set(web.SessionUser, string(userData))

	now := time.Now().UTC()
	sess.Set(web.SessionCreatedAt, now.Format(time.RFC3339))
	sess.Set(web.SessionExpiresAt, now.Add(time.Duration(sessionRoutes.sessionOptions.MaxAge)*time.Second).Format(time.RFC3339))

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

	csrfmiddleware.MakeNewCSRFTokenResp(c)
	//return c.NoContent(http.StatusOK)
}

func (sessionRoutes *SessionRoutes) SessionApiKeySignInRoute(c *gin.Context) {
	validator, err := authentication.NewValidator(c).ParseSignInRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	authUser, err := userdbcache.FindUserByApiKey(validator.UserBodyReq.ApiKey)

	if err != nil {
		web.UserDoesNotExistResp(c)
		return
	}

	if authUser.EmailVerifiedAt == userdb.EmailNotVerifiedDate {
		web.EmailNotVerifiedReq(c)
		return
	}

	// roles, err := userdbcache.UserRoleSet(authUser)

	// if err != nil {
	// 	web.ForbiddenResp(c, authentication.ErrUserRoles)
	// 	return
	// }

	//roleClaim := auth.MakeClaim(roles)

	if !auth.UserHasWebLoginInRole(authUser) {
		web.UserNotAllowedToSignInErrorResp(c)
		return
	}

	err = sessionRoutes.initSession(c, authUser) //, roleClaim)

	if err != nil {
		web.BadReqResp(c, auth.ErrCreatingSession)
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
	user, ok := c.Get(web.SessionUser)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		auth.TokenErrorResp(c)

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
	user, ok := c.Get(web.SessionUser)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		auth.TokenErrorResp(c)

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
	user, ok := c.Get(web.SessionUser)

	for key := range c.Keys {
		log.Debug().Msgf("key %s", key)
	}

	if !ok {
		auth.TokenErrorResp(c)

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
	sessionRoutes.OTPRoutes.Email6DigitOTPRoute(c)
}

func (sessionRoutes *SessionRoutes) SessionSignInUsingEmailAndOTPRoute(c *gin.Context) {

	validator, err := authentication.NewValidator(c).CheckEmailIsWellFormed().Ok()

	if err != nil {
		web.BadReqResp(c, err)
		return
	}

	username := validator.Address.Address

	err = sessionRoutes.OTPRoutes.OTP.ValidateOTP(username, validator.UserBodyReq.OTP)

	if err != nil {
		if auth.IsRateLimitError(err) {
			web.TooManyRequestsResp(c, err)
		} else {
			web.BadReqResp(c, err)
		}

		return
	}

	authUser, err := userdbcache.CreateUserFromOAuth2(username, validator.Address)

	if err != nil {
		log.Debug().Msgf("error creating user from otp sign in: %v", err)
		c.Error(err)
		return
	}

	sessionRoutes.sessionSignInUsingOAuth2(c, authUser)
}

func (sessionRoutes *SessionRoutes) sessionSignInUsingOAuth2(c *gin.Context, authUser *auth.AuthUser) {

	// roles, err := userdbcache.UserRoleSet(authUser)

	// if err != nil {
	// 	web.BadReqResp(c, authentication.ErrUserRoles)
	// }

	//roleClaim := auth.MakeClaim(roles)

	//log.Debug().Msgf("user %v", authUser)

	if !auth.UserHasWebLoginInRole(authUser) {
		web.UserNotAllowedToSignInErrorResp(c)
	}

	err := sessionRoutes.initSession(c, authUser) // roleClaim)

	if err != nil {
		web.UnauthorizedResp(c, err)
		return
	}

	//log.Debug().Msgf("token %s", token)

	csrfmiddleware.MakeNewCSRFTokenResp(c)

	//web.MakeOkResp(c, "user has been signed in")
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sessionRoutes *SessionRoutes) SessionPasswordlessValidateSignInRoute(c *gin.Context) {

	authentication.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *authentication.Validator) {

		if validator.Claims.Type != auth.TokenTypePasswordless {
			auth.WrongTokenTypeReq(c)
			return
		}

		authUser := validator.AuthUser

		// roles, err := userdbcache.UserRoleSet(authUser)

		// if err != nil {
		// 	c.Error(err)
		// 	return
		// }

		//roleClaim := auth.MakeClaim(roles)

		//log.Debug().Msgf("user %v", authUser)

		if !auth.UserHasWebLoginInRole(authUser) {
			web.UserNotAllowedToSignInErrorResp(c)
			return
		}

		err := sessionRoutes.initSession(c, authUser) //, roleClaim)

		if err != nil {
			web.InternalErrorResp(c, err)
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
	sess.Options(middleware.SessionOptsClear) //.SESSION_OPT_ZERO)
	sess.Save()                               //c.Request(), c.Response())

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
		web.UnauthorizedResp(c, ErrSessionExpired)
		return
	}

	web.MakeDataResp(c, "", sessionInfo)
}

func (sessionRoutes *SessionRoutes) SessionNewCSRFTokenRoute(c *gin.Context) {
	csrfmiddleware.MakeNewCSRFTokenResp(c)
}

func (sessionRoutes *SessionRoutes) SessionRefreshRoute(c *gin.Context) {
	user, ok := c.Get(web.SessionUser)

	if !ok {
		web.UnauthorizedResp(c, ErrNoSessionUser)
		return
	}

	// refresh user
	authUser, err := userdbcache.FindUserById(user.(*auth.AuthUser).Id)

	if err != nil {
		web.UnauthorizedResp(c, auth.ErrUserDoesNotExist)
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

	sess.Set(web.SessionUser, string(userData))

	err = sess.Save() //c.Request(), c.Response())

	if err != nil {
		web.InternalErrorResp(c, ErrSavingSession)
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
	user, _ := c.Get(web.SessionUser)

	authUser := user.(*auth.AuthUser)

	var token string
	var err error

	roles := auth.GetRolesFromUser(authUser)

	switch tokenType {
	case "access":
		// Generate encoded token and send it as response.
		token, err = tokengen.AccessToken(c, authUser.Id, roles)

		if err != nil {
			err = fmt.Errorf("error creating access token: %w", err)
		}
	case "update":
		// Generate encoded token and send it as response.
		token, err = tokengen.UpdateToken(c, authUser.Id, roles)

		if err != nil {
			err = fmt.Errorf("error creating update token: %w", err)
		}
	default:
		err = fmt.Errorf("unknown token type")
	}

	if err != nil {
		web.InternalErrorResp(c, err)
		return
	}

	web.MakeDataResp(c, "", &web.TokenResp{Token: token})

}

func UserFromSessionRoute(c *gin.Context) {
	user, ok := c.Get(web.SessionUser)

	if !ok {
		web.BadReqResp(c, ErrNoSessionUser)
		return
	}

	web.MakeDataResp(c, "", user)
}

type PasswordUpdateReq struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func SessionUpdateUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	// Read the session info, which includes the AuthUser

	sessionData, err := middleware.ReadSessionInfo(c, session)

	if err != nil {
		c.Error(err)
		return
	}

	authUser := sessionData.AuthUser

	authentication.NewValidator(c).CheckUsernameIsWellFormed().Success(func(validator *authentication.Validator) {

		err = userdbcache.SetUserInfo(authUser,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.FirstName,
			validator.UserBodyReq.LastName,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		//SendUserInfoUpdatedEmail(c, authUser)

		email := mailserver.MailItem{
			Name: authUser.FirstName,
			To:   authUser.Email,
			//Token:     passwordlessToken,
			EmailType: edbmail.EmailQueueTypeAccountUpdated,
			//Ttl:       fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		mailqueue.SendMail(&email)
	})
}

func SessionUpdatePasswordRoute(c *gin.Context) {
	user, _ := c.Get(web.SessionUser)

	authUser := user.(*auth.AuthUser)

	// use current session user
	authUser, err := userdbcache.FindUserById(authUser.Id)

	if err != nil {
		web.BadReqResp(c, auth.ErrUserDoesNotExist)
		return
	}

	var req PasswordUpdateReq

	err = c.ShouldBindJSON(&req)

	if err != nil {
		web.BadReqResp(c, web.ErrInvalidBody)
		return
	}

	err = userdb.CheckPassword(req.Password)

	if err != nil {
		web.BadReqResp(c, auth.ErrPasswordDoesNotMeetCriteria)
		return
	}

	err = authUser.CheckPasswordsMatch(req.Password)

	if err != nil {
		web.BadReqResp(c, auth.ErrPasswordsDoNotMatch)
		return
	}

	err = userdbcache.SetPassword(authUser, req.NewPassword)

	if err != nil {
		web.BadReqResp(c, auth.ErrCouldNotUpdatePassword)
		return
	}

	web.MakeOkResp(c, "password updated successfully")
}
