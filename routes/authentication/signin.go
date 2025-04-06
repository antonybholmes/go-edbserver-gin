package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"
	"github.com/gin-gonic/gin"
)

func UserSignedInResp(c *gin.Context) {
	web.MakeOkResp(c, "user signed in")
}

func PasswordlessEmailSentResp(c *gin.Context) {
	web.MakeOkResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c *gin.Context) {
	NewValidator(c).ParseLoginRequestBody().Success(func(validator *Validator) {

		if validator.LoginBodyReq.Password == "" {
			PasswordlessSigninEmailRoute(c, validator)
			return
		}

		authUser, err := userdbcache.FindUserByUsername(validator.LoginBodyReq.Username)

		if err != nil {
			web.UserDoesNotExistResp(c)
			return
		}

		if authUser.EmailVerifiedAt == 0 {
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

		refreshToken, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			web.TokenErrorResp(c)
			return
		}

		accessToken, err := tokengen.AccessToken(c, authUser.Uuid, roleClaim) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			web.TokenErrorResp(c)
			return
		}

		web.MakeDataResp(c, "", &web.LoginResp{
			RefreshToken: refreshToken,
			AccessToken:  accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessSigninEmailRoute(c *gin.Context, validator *Validator) {
	if validator == nil {
		validator = NewValidator(c)
	}

	validator.LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		authUser := validator.AuthUser

		passwordlessToken, err := tokengen.PasswordlessToken(c,
			authUser.Uuid,
			validator.LoginBodyReq.RedirectUrl)

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

		//log.Debug().Msgf("t %s ", passwordlessToken)

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Token:     passwordlessToken,
			EmailType: mailer.QUEUE_EMAIL_TYPE_PASSWORDLESS,
			Ttl:       fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		queue.PublishEmail(&email)

		//if err != nil {
		//	return web.ErrorReq(err)
		//}

		web.MakeOkResp(c, "check your email for a magic link to sign in")
	})
}

func PasswordlessSignInRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.PASSWORDLESS_TOKEN {
			web.WrongTokentTypeReq(c)
			return
		}

		authUser := validator.AuthUser

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

		t, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			web.TokenErrorResp(c)
			return
		}

		web.MakeDataResp(c, "", &web.RefreshTokenResp{RefreshToken: t})
	})
}
