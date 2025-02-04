package authorization

import (
	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/rdb"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/gin-gonic/gin"

	"github.com/antonybholmes/go-mailer"
)

type NameReq struct {
	Name string `json:"name"`
}

func UserUpdatedResp(c *gin.Context) {
	routes.MakeOkResp(c, "account updated")
}

func UpdateUserRoute(c *gin.Context) {

	authenticationroutes.NewValidator(c).ParseLoginRequestBody().LoadTokenClaims().Success(func(validator *authenticationroutes.Validator) {

		//db, err := userdbcache.AutoConn(nil) //not clear on what is needed for the user and password

		uuid := validator.Claims.Uuid

		authUser, err := userdbcache.FindUserByUuid(uuid)

		if err != nil {
			c.Error(err)
			return
		}

		err = userdbcache.SetUserInfo(authUser,
			validator.LoginBodyReq.Username,
			validator.LoginBodyReq.FirstName,
			validator.LoginBodyReq.LastName,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		//return SendUserInfoUpdatedEmail(c, authUser)

		// reload user details
		authUser, err = userdbcache.FindUserByUuid(uuid)

		if err != nil {
			c.Error(err)
			return
		}

		// send email notification of change
		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_ACCOUNT_UPDATED}
		rdb.PublishEmail(&email)

		// send back updated user to having to do a separate call to get the new data
		routes.MakeDataResp(c, "account updated confirmation email sent", authUser)
	})
}

func UserRoute(c *gin.Context) {
	authenticationroutes.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *authenticationroutes.Validator) {
			routes.MakeDataResp(c, "", validator.AuthUser)
		})
}

func SendUserInfoUpdatedEmail(c *gin.Context, authUser *auth.AuthUser) {

	file := "templates/email/account/updated.html"

	go authenticationroutes.SendEmailWithToken("Account Updated",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	UserUpdatedResp(c)

}
