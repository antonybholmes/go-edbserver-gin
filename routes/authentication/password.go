package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-edbserver-gin/config"
	email "github.com/antonybholmes/go-mail"

	"github.com/antonybholmes/go-mail/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/token"
	"github.com/antonybholmes/go-web/auth/token/tokengen"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PasswordHandler struct {
	Config *config.AppConfig
}

func (h *PasswordHandler) PasswordUpdatedResp(c *gin.Context) {
	web.MakeOkResp(c, "password updated")
}

// Start passwordless login by sending an email
func (h *PasswordHandler) SendResetPasswordFromUsernameEmailRoute(c *gin.Context) {
	middleware.NewValidator(c).LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *middleware.Validator) {
		authUser := validator.AuthUser
		//req := validator.SignInBodyReq

		otpToken, err := tokengen.Instance().MakeResetPasswordToken(c, authUser, jwt.ClaimStrings{"reset-password"})

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

		email := email.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			Payload:   &email.Payload{DataType: "jwt", Data: otpToken},
			EmailType: email.EmailTypePasswordReset,
			TTL:       fmt.Sprintf("%d minutes", int(h.Config.Auth.ShortTtlMins)),
			LinkUrl:   h.Config.Account.ResetPasswordUrl}

		mailqueue.Instance().SendMail(c.Request.Context(), &email)

		web.MakeOkResp(c, "check your email for a password reset link")
	})
}

func UpdatePasswordRoute(c *gin.Context) {
	validator, err := middleware.NewValidator(c).ParseSignInRequestBody().LoadAuthUserFromToken().Ok()

	if err != nil {
		c.Error(err)
		return
	}

	if validator.Claims.Type != token.TokenTypeResetPassword {
		web.BadReqResp(c, token.ErrInvalidTokenType)
		return
	}

	authUser := validator.AuthUser

	err = auth.CheckOTPValid(authUser, validator.Claims.OneTimePasscode)

	if err != nil {
		c.Error(err)
		return
	}

	_, err = userdbcache.Instance().SetPassword(authUser, validator.UserBodyReq.Password, false)

	if err != nil {
		web.BadReqResp(c, err)
		return
	}

	//return SendPasswordEmail(c, validator.AuthUser, validator.Req.Password)

	email := email.MailItem{
		Name:      authUser.Name,
		To:        authUser.Email,
		EmailType: email.EmailTypePasswordUpdated}
	mailqueue.Instance().SendMail(c.Request.Context(), &email)

	web.MakeOkResp(c, "password updated confirmation email sent")

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
