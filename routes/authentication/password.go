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

func PasswordUpdatedResp(c *gin.Context) {
	web.MakeOkResp(c, "password updated")
}

// Start passwordless login by sending an email
func SendResetPasswordFromUsernameEmailRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {
		authUser := validator.AuthUser
		//req := validator.SignInBodyReq

		otpToken, err := tokengen.MakeResetPasswordToken(c, authUser)

		if err != nil {
			c.Error(err)
			return
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/password/reset/web.html"
		// } else {
		// 	file = "templates/email/password/reset/api.html"
		// }

		// go authentication.SendEmailWithToken("Password Reset",
		// 	authUser,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Payload:   &mailserver.Payload{DataType: "jwt", Data: otpToken},
			EmailType: edbmail.QUEUE_EMAIL_TYPE_PASSWORD_RESET,
			TTL:       fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			LinkUrl:   consts.URL_RESET_PASSWORD}
		mailqueue.SendMail(&email)

		//if err != nil {
		//	return web.ErrorReq(err)
		//}

		web.MakeOkResp(c, "check your email for a password reset link")
	})
}

func UpdatePasswordRoute(c *gin.Context) {
	NewValidator(c).ParseSignInRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.RESET_PASSWORD_TOKEN {
			web.BadReqResp(c, auth.ErrInvalidTokenType)
			return
		}

		authUser := validator.AuthUser

		err := auth.CheckOTPValid(authUser, validator.Claims.OneTimePasscode)

		if err != nil {
			c.Error(err)
			return
		}

		err = userdbcache.SetPassword(authUser, validator.UserBodyReq.Password)

		if err != nil {
			web.BadReqResp(c, err)
			return
		}

		//return SendPasswordEmail(c, validator.AuthUser, validator.Req.Password)

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: edbmail.QUEUE_EMAIL_TYPE_PASSWORD_UPDATED}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "password updated confirmation email sent")
	})
}

// func SendPasswordEmail(c *gin.Context, authUser *auth.AuthUser, password string) error {

// 	var file string

// 	if password == "" {
// 		file = "templates/email/password/switch-to-passwordless.html"
// 	} else {
// 		file = "templates/email/password/updated.html"
// 	}

// 	go SendEmailWithToken("Password Updated",
// 		authUser,
// 		file,
// 		"",
// 		"",
// 		"")

// 	return PasswordUpdatedResp(c)

// }
