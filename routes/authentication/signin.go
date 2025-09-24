package authentication

import (
	"fmt"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"
	"github.com/gin-gonic/gin"
)

func PasswordlessEmailSentResp(c *gin.Context) {
	web.MakeOkResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c *gin.Context) {
	NewValidator(c).ParseLoginRequestBody().Success(func(validator *Validator) {

		if validator.UserBodyReq.Password == "" {
			PasswordlessSigninEmailRoute(c, validator)
			return
		}

		authUser, err := userdbcache.FindUserByUsername(validator.UserBodyReq.Username)

		if err != nil {
			web.UserDoesNotExistResp(c)
			return
		}

		if authUser.EmailVerifiedAt == 0 {
			web.EmailNotVerifiedReq(c)
			return
		}

		roles, err := userdbcache.UserRoleSet(authUser)

		if err != nil {
			web.ForbiddenResp(c, "could not get user roles")
			return
		}

		//roleClaim := auth.MakeClaim(roles)

		if !auth.HasSignInRole(roles) {
			web.UserNotAllowedToSignInErrorResp(c)
			return
		}

		err = authUser.CheckPasswordsMatch(validator.UserBodyReq.Password)

		if err != nil {
			c.Error(err)
			return
		}

		refreshToken, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			web.TokenErrorResp(c)
			return
		}

		accessToken, err := tokengen.AccessToken(c, authUser.PublicId, roles.Keys()) //auth.MakeClaim(authUser.Roles))

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

		passwordlessToken, err := tokengen.MakePasswordlessToken(c,
			authUser.PublicId,
			validator.UserBodyReq.RedirectUrl)

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

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Payload:   &mailserver.Payload{DataType: "code", Data: passwordlessToken},
			EmailType: edbmail.QUEUE_EMAIL_TYPE_PASSWORDLESS,
			TTL:       fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		mailqueue.SendMail(&email)

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

		roles, err := userdbcache.UserRoleSet(authUser)

		if err != nil {
			web.ForbiddenResp(c, "could not get user roles")
			return
		}

		//roleClaim := auth.MakeClaim(roles)

		if !auth.HasSignInRole(roles) {
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
