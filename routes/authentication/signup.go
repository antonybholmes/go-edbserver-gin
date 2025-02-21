package authenticationroutes

import (
	"fmt"

	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-mailer/queue"

	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-mailer"
	"github.com/gin-gonic/gin"
)

func SignupRoute(c *gin.Context) {
	NewValidator(c).CheckEmailIsWellFormed().Success(func(validator *Validator) {
		req := validator.LoginBodyReq

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
		//	return routes.ErrorReq(err)
		//}

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			Token:     token,
			EmailType: mailer.QUEUE_EMAIL_TYPE_VERIFY,
			Ttl:       fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			LinkUrl:   consts.URL_VERIFY_EMAIL,
			//VisitUrl:    req.VisitUrl
		}

		queue.PublishEmail(&email)

		routes.MakeOkResp(c, "check your email for a verification link")
	})
}

func EmailAddressVerifiedRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromToken().Success(func(validator *Validator) {

		authUser := validator.AuthUser

		// if verified, stop and just return true
		if authUser.EmailVerifiedAt == 0 {
			routes.MakeOkResp(c, "")
		}

		err := userdbcache.SetIsVerified(authUser.Uuid)

		if err != nil {
			routes.MakeSuccessResp(c, "unable to verify user", false)
		}

		// file := "templates/email/verify/verified.html"

		// go SendEmailWithToken("Email Address Verified",
		// 	authUser,
		// 	file,
		// 	"",
		// 	"",
		// 	"")

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.QUEUE_EMAIL_TYPE_VERIFIED}
		queue.PublishEmail(&email)

		routes.MakeOkResp(c, "email address verified")
	})
}
