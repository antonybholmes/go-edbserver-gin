package authentication

import (
	"fmt"

	"net/mail"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"

	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/gin-gonic/gin"
)

// Start passwordless login by sending an email
func SendResetEmailEmailRoute(c *gin.Context) {
	NewValidator(c).ParseSignInRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) {
		authUser := validator.AuthUser
		req := validator.UserBodyReq

		newEmail, err := mail.ParseAddress(req.Email)

		if err != nil {
			c.Error(err)
			return
		}

		otpToken, err := tokengen.MakeResetEmailToken(c, authUser, newEmail)

		if err != nil {
			c.Error(err)
			return
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/email/reset/web.html"
		// } else {
		// 	file = "templates/email/email/reset/api.html"
		// }

		// go BaseSendEmailWithToken("Update Email",
		// 	authUser,
		// 	newEmail,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Payload:   &mailserver.Payload{DataType: "jwt", Data: otpToken},
			EmailType: edbmail.EmailQueueTypeEmailReset,
			TTL:       fmt.Sprintf("%d minutes", int(consts.ShortTtlMins.Minutes())),
			LinkUrl:   consts.UrlResetEmail,
		}
		mailqueue.SendMail(&email)

		//if err != nil {
		//	return web.ErrorReq(err)
		//}

		web.MakeOkResp(c, "check your email for a change email link")
	})
}

func UpdateEmailRoute(c *gin.Context) {
	NewValidator(c).CheckEmailIsWellFormed().LoadAuthUserFromToken().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.TokenTypeChangeEmail {
			auth.WrongTokenTypeReq(c)
		}

		err := auth.CheckOTPValid(validator.AuthUser,
			validator.Claims.OneTimePasscode)

		if err != nil {
			c.Error(err)
			return
		}

		authUser := validator.AuthUser
		publicId := authUser.PublicId

		err = userdbcache.SetEmailAddress(authUser,
			validator.Address,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		authUser, err = userdbcache.FindUserByPublicId(publicId)

		if err != nil {
			c.Error(err)
			return
		}

		//return SendEmailChangedEmail(c, authUser)

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: edbmail.EmailQueueTypeEmailUpdated}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "email updated confirmation email sent")
	})
}
