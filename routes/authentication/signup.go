package authentication

import (
	"fmt"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/tokengen"
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

		token, err := tokengen.MakeVerifyEmailToken(c, authUser, req.RedirectUrl)

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
			Payload:   &mailserver.Payload{DataType: "jwt", Data: token},
			EmailType: edbmail.EmailQueueTypeVerify,
			TTL:       fmt.Sprintf("%d minutes", int(consts.ShortTtlMins.Minutes())),
			LinkUrl:   consts.UrlVerifyEmail,
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

		err := userdbcache.SetIsVerified(authUser.Id)

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
			EmailType: edbmail.EmailQueueTypeVerified}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "email address verified")
	})
}
