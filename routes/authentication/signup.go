package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-edbserver-gin/config"
	email "github.com/antonybholmes/go-mail"

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

		email := email.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			Payload:   &email.Payload{DataType: "jwt", Data: token},
			EmailType: email.EmailTypeVerify,
			TTL:       fmt.Sprintf("%d minutes", h.Config.Auth.ShortTtlMins),
			LinkUrl:   h.Config.Account.VerifyEmailUrl,
			//VisitUrl:    req.VisitUrl
		}

		mailqueue.Instance().SendMail(c.Request.Context(), &email)

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

		email := email.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			EmailType: email.EmailTypeVerified}
		mailqueue.Instance().SendMail(c.Request.Context(), &email)

		web.MakeOkResp(c, "email address verified")
	})
}
