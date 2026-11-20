package authentication

import (
	"fmt"

	"net/mail"

	edbmail "github.com/antonybholmes/go-edb-mailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/config"

	mailserver "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/token"
	"github.com/antonybholmes/go-web/auth/token/tokengen"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/golang-jwt/jwt/v5"

	"github.com/antonybholmes/go-mail/mailqueue"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	Config *config.AppConfig
}

// Start passwordless login by sending an email
func (h *EmailHandler) SendResetEmailEmailRoute(c *gin.Context) {
	validator, err := middleware.NewValidator(c).ParseSignInRequestBody().LoadAuthUserFromToken().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	authUser := validator.AuthUser
	req := validator.UserBodyReq

	newEmail, err := mail.ParseAddress(req.Email)

	if err != nil {
		c.Error(err)
		return
	}

	otpToken, err := tokengen.Instance().MakeResetEmailToken(c, authUser, jwt.ClaimStrings{"reset-password"}, newEmail)

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
		Name:      authUser.Name,
		To:        authUser.Email,
		Payload:   &mailserver.Payload{DataType: "jwt", Data: otpToken},
		EmailType: edbmail.EmailTypeEmailReset,
		TTL:       fmt.Sprintf("%d minutes", int(h.Config.Auth.ShortTtlMins)),
		LinkUrl:   h.Config.Account.ResetEmailUrl,
	}

	mailqueue.SendMail(&email)

	//if err != nil {
	//	return web.ErrorReq(err)
	//}

	web.MakeOkResp(c, "check your email for a change email link")

}

func (h *EmailHandler) UpdateEmailRoute(c *gin.Context) {
	middleware.NewValidator(c).CheckEmailIsWellFormed().LoadAuthUserFromToken().Success(func(validator *middleware.Validator) {

		if validator.Claims.Type != token.TokenTypeChangeEmail {
			auth.WrongTokenTypeReq(c)
		}

		err := auth.CheckOTPValid(validator.AuthUser,
			validator.Claims.OneTimePasscode)

		if err != nil {
			c.Error(err)
			return
		}

		authUser := validator.AuthUser
		id := authUser.Id

		err = userdbcache.Instance().SetEmailAddress(authUser,
			validator.Address,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		authUser, err = userdbcache.Instance().FindUserById(id)

		if err != nil {
			c.Error(err)
			return
		}

		//return SendEmailChangedEmail(c, authUser)

		email := mailserver.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			EmailType: edbmail.EmailTypeEmailUpdated}

		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "email updated confirmation email sent")
	})
}
