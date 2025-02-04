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
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

var SESSION_OPT_ZERO *sessions.Options

//var SESSION_OPT_24H *sessions.Options
//var SESSION_OPT_30_DAYS *sessions.Options
//var SESSION_OPT_7_DAYS *sessions.Options

func init() {

	// HttpOnly and Secure are disabled so we can use them
	// cross domain for testing
	// http only false to allow js to delete etc on the client side

	// For sessions that should end when browser closes
	SESSION_OPT_ZERO = &sessions.Options{
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
	sessionOptions *sessions.Options
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

	return &SessionRoutes{sessionOptions: &options}
}

func (sr *SessionRoutes) SessionUsernamePasswordSignInRoute(c *gin.Context) {
	validator, err := NewValidator(c).ParseLoginRequestBody().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.LoginBodyReq.Password == "" {
		return PasswordlessSigninEmailRoute(c, validator)
	}

	user := validator.LoginBodyReq.Username

	authUser, err := userdbcache.FindUserByUsername(user)

	if err != nil {
		return routes.UserDoesNotExistReq()
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		return routes.EmailNotVerifiedReq()
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		return routes.AuthErrorReq("could not get user roles")
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		return routes.UserNotAllowedToSignIn()
	}

	err = authUser.CheckPasswordsMatch(validator.LoginBodyReq.Password)

	if err != nil {
		c.Error(err)
		return
	}

	sess, err := session.Get(consts.SESSION_NAME, c)

	if err != nil {
		c.Error(err)
		return
	}

	// set session options
	if validator.LoginBodyReq.StaySignedIn {
		sess.Options = sr.sessionOptions
	} else {
		sess.Options = SESSION_OPT_ZERO
	}

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roleClaim //auth.MakeClaim(authUser.Roles)

	sess.Save(c.Request(), c.Response())

	return UserSignedInResp(c)
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
		return routes.UserDoesNotExistReq()
	}

	if authUser.EmailVerifiedAt == auth.EMAIL_NOT_VERIFIED_TIME_S {
		return routes.EmailNotVerifiedReq()
	}

	roles, err := userdbcache.UserRoleList(authUser)

	if err != nil {
		return routes.AuthErrorReq("could not get user roles")
	}

	roleClaim := auth.MakeClaim(roles)

	if !auth.CanSignin(roleClaim) {
		return routes.UserNotAllowedToSignIn()
	}

	err = sr.initSession(c, authUser) //, roleClaim)

	if err != nil {
		return fmt.Errorf("%s", ERROR_CREATING_SESSION)
	}

	return c.String(http.StatusOK, "Session set")

	// resp, err := readSession(c)

	// if err != nil {
	// 	return routes.ErrorReq(ERROR_CREATING_SESSION)
	// }

	// routes.MakeDataResp(c, "", resp)
}

func (sr *SessionRoutes) SessionSignInUsingAuth0Route(c *gin.Context) {
	user := c.Get("user").(*jwt.Token)
	tokenClaims := user.Claims.(*auth.Auth0TokenClaims)

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
		return routes.ErrorResp("user roles not found")
	}

	roleClaim := auth.MakeClaim(roles)

	//log.Debug().Msgf("user %v", authUser)

	if !auth.CanSignin(roleClaim) {
		return routes.UserNotAllowedToSignIn()
	}

	err = sr.initSession(c, authUser) // roleClaim)

	if err != nil {
		c.Error(err)
		return
	}

	return UserSignedInResp(c)
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
		c.Error(err)
		return
	}

	sess, err := session.Get(consts.SESSION_NAME, c)

	if err != nil {
		return fmt.Errorf("%s", ERROR_CREATING_SESSION)
	}

	// set session options
	sess.Options = sr.sessionOptions

	//sess.Values[SESSION_PUBLICID] = authUser.PublicId
	//sess.Values[SESSION_ROLES] = roles //auth.MakeClaim(authUser.Roles)
	sess.Values[SESSION_USER] = string(userData)

	now := time.Now().UTC()
	sess.Values[SESSION_CREATED_AT] = now.Format(time.RFC3339)
	sess.Values[SESSION_EXPIRES_AT] = now.Add(time.Duration(sess.Options.MaxAge) * time.Second).Format(time.RFC3339)

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
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

	return c.NoContent(http.StatusOK)
}

func ReadSessionInfo(c *gin.Context) (*SessionInfo, error) {
	sess, err := session.Get(consts.SESSION_NAME, c)

	if err != nil {
		return nil, fmt.Errorf(ERROR_CREATING_SESSION)
	}

	userData, _ := sess.Values[SESSION_USER].(string)

	var user auth.AuthUser

	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	//publicId, _ := sess.Values[SESSION_PUBLICID].(string)
	//roles, _ := sess.Values[SESSION_ROLES].(string)
	createdAt, _ := sess.Values[SESSION_CREATED_AT].(string)
	expires, _ := sess.Values[SESSION_EXPIRES_AT].(string)
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
	authUser := c.Get("authUser").(*auth.AuthUser)

	// refresh user
	authUser, err := userdbcache.FindUserById(authUser.Id)

	if err != nil {
		c.Error(err)
		return
	}

	//
	// For the moment just update the user info

	//err = sr.initSession(c, authUser)

	sess, err := session.Get(consts.SESSION_NAME, c)

	if err != nil {
		return fmt.Errorf("%s", ERROR_CREATING_SESSION)
	}

	userData, err := json.Marshal(authUser)

	if err != nil {
		c.Error(err)
		return
	}

	log.Debug().Msgf("saving %s", string(userData))

	sess.Values[SESSION_USER] = string(userData)

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		c.Error(err)
		return
	}

	return c.NoContent(http.StatusOK)
}

// Validate the passwordless token we generated and create
// a user session. The session acts as a refresh token and
// can be used to generate access tokens to use resources
func (sr *SessionRoutes) SessionPasswordlessValidateSignInRoute(c *gin.Context) {

	NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			return routes.WrongTokentTypeReq()
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
			return routes.UserNotAllowedToSignIn()
		}

		err = sr.initSession(c, authUser) //, roleClaim)

		if err != nil {
			c.Error(err)
			return
		}

		return UserSignedInResp(c)
	})
}

func SessionSignOutRoute(c *gin.Context) {
	sess, err := session.Get(consts.SESSION_NAME, c)

	if err != nil {
		return routes.ErrorResp(ERROR_CREATING_SESSION)
	}

	log.Debug().Msgf("invalidate session")

	// invalidate by time
	sess.Values[SESSION_USER] = ""
	//sess.Values[SESSION_ROLES] = ""
	sess.Values[SESSION_CREATED_AT] = ""
	sess.Values[SESSION_EXPIRES_AT] = ""
	sess.Options.MaxAge = -1

	sess.Save(c.Request(), c.Response())

	routes.MakeOkResp(c, "user has been signed out")
}

func NewAccessTokenFromSessionRoute(c *gin.Context) {
	// sess, _ := session.Get(consts.SESSION_NAME, c)

	// userData, ok := sess.Values["user"].(string)

	// if !ok {
	// 	return routes.ErrorReq(fmt.Errorf("malformed user info"))
	// }

	// var user auth.AuthUser
	// if err := json.Unmarshal([]byte(userData), &user); err != nil {
	// 	return routes.ErrorReq(err)
	// }

	user := c.Get("authUser").(*auth.AuthUser)

	//publicId, _ := sess.Values[SESSION_PUBLICID].(string)
	//r//oles, _ := sess.Values[SESSION_ROLES].(string)

	// if publicId == "" {
	// 	return routes.ErrorReq(fmt.Errorf("public id cannot be empty"))
	// }

	// generate a new token from what is stored in the sesssion
	t, err := tokengen.AccessToken(c, user.Uuid, auth.MakeClaim(user.Roles))

	if err != nil {
		return routes.TokenErrorReq()
	}

	routes.MakeDataResp(c, "", &routes.AccessTokenResp{AccessToken: t})
}

func UserFromSessionRoute(c *gin.Context) {
	user := c.Get("authUser").(*auth.AuthUser)

	routes.MakeDataResp(c, "", user)
}
