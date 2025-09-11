package authorization

import (
	"github.com/antonybholmes/go-web/userdbcache"
	mailserver "github.com/antonybholmes/go_mailserver"
	"github.com/antonybholmes/go_mailserver/queue"
	"github.com/rs/zerolog/log"

	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

type NameReq struct {
	Name string `json:"name"`
}

func UserUpdatedResp(c *gin.Context) {
	web.MakeOkResp(c, "account updated")
}

func UpdateUserRoute(c *gin.Context) {

	authenticationroutes.NewValidator(c).ParseLoginRequestBody().LoadTokenClaims().Success(func(validator *authenticationroutes.Validator) {

		//db, err := userdbcache.AutoConn(nil) //not clear on what is needed for the user and password

		publicId := validator.Claims.UserId

		log.Debug().Msgf("UpdateUserRoute: publicId: %s ", publicId)

		authUser, err := userdbcache.FindUserByPublicId(publicId)

		if err != nil {
			c.Error(err)
			return
		}

		err = userdbcache.SetUserInfo(authUser,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.FirstName,
			validator.UserBodyReq.LastName,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		//return SendUserInfoUpdatedEmail(c, authUser)

		// reload user details
		authUser, err = userdbcache.FindUserByPublicId(publicId)

		if err != nil {
			c.Error(err)
			return
		}

		// send email notification of change
		email := mailserver.QueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailserver.QUEUE_EMAIL_TYPE_ACCOUNT_UPDATED}
		queue.PublishEmail(&email)

		// send back updated user to having to do a separate call to get the new data
		web.MakeDataResp(c, "account updated confirmation email sent", authUser)
	})
}

func UserRoute(c *gin.Context) {
	authenticationroutes.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *authenticationroutes.Validator) {
			web.MakeDataResp(c, "", validator.AuthUser)
		})
}
