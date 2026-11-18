package authentication

import (
	"fmt"

	edbmailconfig "github.com/antonybholmes/go-edb-mailserver/config"
	edbmail "github.com/antonybholmes/go-edb-mailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/config"

	mailserver "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-mail/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth/token/tokengen"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignupHandler struct {
	Config *config.AppConfig
}

func (h *SignupHandler) SignupRoute(c *gin.Context) {
	middleware.NewValidator(c).CheckEmailIsWellFormed().Success(func(validator *middleware.Validator) {
		req := validator.UserBodyReq

		authUser, err := userdbcache.Instance().CreateUserFromSignup(req)

		if err != nil {
			c.Error(err)
			return
		}

		token, err := tokengen.Instance().MakeVerifyEmailToken(c, authUser, jwt.ClaimStrings{"verify-email"}, req.RedirectUrl)

		//log.Debug().Msgf("%s", otpJwt)

		if err != nil {
			c.Error(err)
			return
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/verify/web.html"
		// } else {
		// 	file = "templates/email/verify/api.html"
		// }

		// go SendEmailWithToken("Email Verification",
		// 	authUser,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		//if err != nil {
		//	return web.ErrorReq(err)
		//}

		email := mailserver.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			Payload:   &mailserver.Payload{DataType: "jwt", Data: token},
			EmailType: edbmail.EmailTypeVerify,
			TTL:       fmt.Sprintf("%d minutes", h.Config.Auth.ShortTtlMins),
			LinkUrl:   edbmailconfig.UrlVerifyEmail,
			//VisitUrl:    req.VisitUrl
		}

		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "check your email for a verification link")
	})
}

func EmailAddressVerifiedRoute(c *gin.Context) {
	middleware.NewValidator(c).LoadAuthUserFromToken().Success(func(validator *middleware.Validator) {

		authUser := validator.AuthUser

		// if verified, stop and just return true
		if authUser.EmailVerifiedAt != nil {
			web.MakeOkResp(c, "")
		}

		_, err := userdbcache.Instance().SetEmailIsVerified(authUser)

		if err != nil {
			web.MakeSuccessResp(c, "unable to verify user", false)
		}

		// file := "templates/email/verify/verified.html"

		// go SendEmailWithToken("Email Address Verified",
		// 	authUser,
		// 	file,
		// 	"",
		// 	"",
		// 	"")

		email := mailserver.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			EmailType: edbmail.EmailTypeVerified}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "email address verified")
	})
}
