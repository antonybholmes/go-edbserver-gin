package authentication

import (
	"fmt"

	"net/mail"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"

	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/gin-gonic/gin"
)

// Start passwordless login by sending an email
func SendResetEmailEmailRoute(c *gin.Context) {
	NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) {
		authUser := validator.AuthUser
		req := validator.LoginBodyReq

		newEmail, err := mail.ParseAddress(req.Email)

		if err != nil {
			c.Error(err)
			return
		}

		otpToken, err := tokengen.ResetEmailToken(c, authUser, newEmail)

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

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Token:     otpToken,
			EmailType: mailer.QUEUE_EMAIL_TYPE_EMAIL_RESET,
			Ttl:       fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			LinkUrl:   consts.URL_RESET_EMAIL,
		}
		queue.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		routes.MakeOkResp(c, "check your email for a change email link")
	})
}

func UpdateEmailRoute(c *gin.Context) {
	NewValidator(c).CheckEmailIsWellFormed().LoadAuthUserFromToken().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.CHANGE_EMAIL_TOKEN {
			routes.WrongTokentTypeReq(c)
		}

		err := auth.CheckOTPValid(validator.AuthUser, validator.Claims.OneTimePasscode)

		if err != nil {
			c.Error(err)
			return
		}

		authUser := validator.AuthUser
		uuid := authUser.Uuid

		err = userdbcache.SetEmailAddress(authUser, validator.Address, false)

		if err != nil {
			c.Error(err)
			return
		}

		authUser, err = userdbcache.FindUserByUuid(uuid)

		if err != nil {
			c.Error(err)
			return
		}

		//return SendEmailChangedEmail(c, authUser)

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.QUEUE_EMAIL_TYPE_EMAIL_UPDATED}
		queue.PublishEmail(&email)

		routes.MakeOkResp(c, "email updated confirmation email sent")
	})
}
