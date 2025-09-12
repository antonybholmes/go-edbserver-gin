package authentication

import (
	"fmt"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"
	"github.com/gin-gonic/gin"
)

func SignupRoute(c *gin.Context) {
	NewValidator(c).CheckEmailIsWellFormed().Success(func(validator *Validator) {
		req := validator.UserBodyReq

		authUser, err := userdbcache.CreateUserFromSignup(req)

		if err != nil {
			c.Error(err)
			return
		}

		token, err := tokengen.VerifyEmailToken(c, authUser, req.RedirectUrl)

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
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Token:     token,
			EmailType: mailserver.QUEUE_EMAIL_TYPE_VERIFY,
			TTL:       fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			LinkUrl:   consts.URL_VERIFY_EMAIL,
			//VisitUrl:    req.VisitUrl
		}

		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "check your email for a verification link")
	})
}

func EmailAddressVerifiedRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromToken().Success(func(validator *Validator) {

		authUser := validator.AuthUser

		// if verified, stop and just return true
		if authUser.EmailVerifiedAt == 0 {
			web.MakeOkResp(c, "")
		}

		err := userdbcache.SetIsVerified(authUser.PublicId)

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
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailserver.QUEUE_EMAIL_TYPE_VERIFIED}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "email address verified")
	})
}
