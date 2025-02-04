package authenticationroutes

import (
	"fmt"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/rdb"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-mailer"
	"github.com/gin-gonic/gin"
)

func UserSignedInResp(c *gin.Context) {
	routes.MakeOkResp(c, "user signed in")
}

func PasswordlessEmailSentResp(c *gin.Context) {
	routes.MakeOkResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c *gin.Context) {
	NewValidator(c).ParseLoginRequestBody().Success(func(validator *Validator) {

		if validator.LoginBodyReq.Password == "" {
			PasswordlessSigninEmailRoute(c, validator)
		}

		authUser, err := userdbcache.FindUserByUsername(validator.LoginBodyReq.Username)

		if err != nil {
			routes.UserDoesNotExistResp(c)
		}

		if authUser.EmailVerifiedAt == 0 {
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

		refreshToken, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			routes.TokenErrorReq(c)
		}

		accessToken, err := tokengen.AccessToken(c, authUser.Uuid, roleClaim) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			routes.TokenErrorReq(c)
		}

		routes.MakeDataResp(c, "", &routes.LoginResp{RefreshToken: refreshToken, AccessToken: accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessSigninEmailRoute(c *gin.Context, validator *Validator) {
	if validator == nil {
		validator = NewValidator(c)
	}

	validator.LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		authUser := validator.AuthUser

		passwordlessToken, err := tokengen.PasswordlessToken(c, authUser.Uuid, validator.LoginBodyReq.VisitUrl)

		if err != nil {
			c.Error(err)
			return
		}

		// var file string

		// if validator.Req.CallbackUrl != "" {
		// 	file = "templates/email/passwordless/web.html"
		// } else {
		// 	file = "templates/email/passwordless/api.html"
		// }

		// go SendEmailWithToken("Passwordless Sign In",
		// 	authUser,
		// 	file,
		// 	passwordlessToken,
		// 	validator.Req.CallbackUrl,
		// 	validator.Req.VisitUrl)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       passwordlessToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_PASSWORDLESS,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			CallBackUrl: validator.LoginBodyReq.CallbackUrl,
			//VisitUrl:    validator.Req.VisitUrl
		}
		rdb.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		routes.MakeOkResp(c, "check your email for a passwordless sign in link")
	})
}

func PasswordlessSignInRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			routes.WrongTokentTypeReq(c)
			return
		}

		authUser := validator.AuthUser

		roles, err := userdbcache.UserRoleList(authUser)

		if err != nil {
			routes.AuthErrorReq(c, "could not get user roles")
			return
		}

		roleClaim := auth.MakeClaim(roles)

		if !auth.CanSignin(roleClaim) {
			routes.UserNotAllowedToSignIn(c)
			return
		}

		t, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			routes.TokenErrorReq(c)
			return
		}

		routes.MakeDataResp(c, "", &routes.RefreshTokenResp{RefreshToken: t})
	})
}
