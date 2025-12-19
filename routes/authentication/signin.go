package authentication

import (
	"errors"
	"fmt"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserRoles = errors.New("could not get user roles")
)

func PasswordlessEmailSentResp(c *gin.Context) {
	web.MakeOkResp(c, "passwordless email sent")
}

func UsernamePasswordSignInRoute(c *gin.Context) {
	middleware.NewValidator(c).ParseSignInRequestBody().Success(func(validator *middleware.Validator) {

		if validator.UserBodyReq.Password == "" {
			PasswordlessSignInEmailRoute(c, validator)
			return
		}

		authUser, err := userdbcache.FindUserByUsername(validator.UserBodyReq.Username)

		if err != nil {
			web.UserDoesNotExistResp(c)
			return
		}

		if authUser.EmailVerifiedAt == nil {
			web.EmailNotVerifiedReq(c)
			return
		}

		// roles, err := userdbcache.UserRoleSet(authUser)

		// if err != nil {
		// 	web.ForbiddenResp(c, ErrUserRoles)
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

		refreshToken, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			auth.TokenErrorResp(c)
			return
		}

		accessToken, err := tokengen.AccessToken(c, authUser.Id, auth.GetRolesFromUser(authUser)) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			auth.TokenErrorResp(c)
			return
		}

		web.MakeDataResp(c, "", &web.SignInResp{
			RefreshToken: refreshToken,
			AccessToken:  accessToken})
	})
}

// Start passwordless login by sending an email
func PasswordlessSignInEmailRoute(c *gin.Context, validator *middleware.Validator) {
	if validator == nil {
		validator = middleware.NewValidator(c)
	}

	validator.LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *middleware.Validator) {

		authUser := validator.AuthUser

		passwordlessToken, err := tokengen.MakePasswordlessToken(c,
			authUser.Id,
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
			Name:      authUser.Name,
			To:        authUser.Email,
			Payload:   &mailserver.Payload{DataType: "code", Data: passwordlessToken},
			EmailType: edbmail.EmailQueueTypePasswordless,
			TTL:       fmt.Sprintf("%d minutes", int(consts.PasswordlessTokenTtlMins.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		mailqueue.SendMail(&email)

		//if err != nil {
		//	return web.ErrorReq(err)
		//}

		web.MakeOkResp(c, "check your email for a magic link to login")
	})
}

func PasswordlessSignInRoute(c *gin.Context) {
	middleware.NewValidator(c).LoadAuthUserFromToken().CheckUserHasVerifiedEmailAddress().Success(func(validator *middleware.Validator) {

		if validator.Claims.Type != auth.TokenTypePasswordless {
			auth.WrongTokenTypeReq(c)
			return
		}

		authUser := validator.AuthUser

		// roles, err := userdbcache.UserRoleSet(authUser)

		// if err != nil {
		// 	web.ForbiddenResp(c, ErrUserRoles)
		// 	return
		// }

		//roleClaim := auth.MakeClaim(roles)

		if !auth.UserHasWebLoginInRole(authUser) {
			web.UserNotAllowedToSignInErrorResp(c)
			return
		}

		t, err := tokengen.RefreshToken(c, authUser) //auth.MakeClaim(authUser.Roles))

		if err != nil {
			auth.TokenErrorResp(c)
			return
		}

		web.MakeDataResp(c, "", &web.RefreshTokenResp{RefreshToken: t})
	})
}
