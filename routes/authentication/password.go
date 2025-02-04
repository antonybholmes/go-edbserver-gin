package authenticationroutes

import (
	"fmt"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/rdb"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-gonic/gin"

	"github.com/antonybholmes/go-mailer"
)

func PasswordUpdatedResp(c *gin.Context) {
	routes.MakeOkResp(c, "password updated")
}

// Start passwordless login by sending an email
func SendResetPasswordFromUsernameEmailRoute(c *gin.Context) {
	NewValidator(c).LoadAuthUserFromUsername().CheckUserHasVerifiedEmailAddress().Success(func(validator *Validator) {
		authUser := validator.AuthUser
		req := validator.LoginBodyReq

		otpToken, err := tokengen.ResetPasswordToken(c, authUser)

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

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       otpToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_PASSWORD_RESET,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			CallBackUrl: req.CallbackUrl,
			VisitUrl:    req.VisitUrl}
		rdb.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		routes.MakeOkResp(c, "check your email for a password reset link")
	})
}

func UpdatePasswordRoute(c *gin.Context) {
	NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) {

		if validator.Claims.Type != auth.RESET_PASSWORD_TOKEN {
			routes.ErrorResp(c, routes.ERROR_WRONG_TOKEN_TYPE)
			return
		}

		authUser := validator.AuthUser

		err := auth.CheckOTPValid(authUser, validator.Claims.OneTimePasscode)

		if err != nil {
			c.Error(err)
			return
		}

		err = userdbcache.SetPassword(authUser, validator.LoginBodyReq.Password)

		if err != nil {
			routes.ErrorResp(c, routes.ERROR_WRONG_TOKEN_TYPE)
			return
		}

		//return SendPasswordEmail(c, validator.AuthUser, validator.Req.Password)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_PASSWORD_UPDATED}
		rdb.PublishEmail(&email)

		routes.MakeOkResp(c, "password updated confirmation email sent")
	})
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
