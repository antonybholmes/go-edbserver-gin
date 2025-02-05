package authenticationroutes

import (
	"fmt"

	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/rdb"
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

		otpToken, err := tokengen.VerifyEmailToken(c, authUser, req.VisitUrl)

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

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       otpToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_VERIFY,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			CallBackUrl: req.CallbackUrl,
			//VisitUrl:    req.VisitUrl
		}

		rdb.PublishEmail(&email)

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

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_VERIFIED}
		rdb.PublishEmail(&email)

		routes.MakeOkResp(c, "email address verified")
	})
}
