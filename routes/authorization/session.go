package authorization

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/middleware"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"

	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/gin-gonic/gin"
)

func SessionUpdateUserRoute(c *gin.Context) {
	sessionData, err := middleware.ReadSessionInfo(c)

	if err != nil {
		c.Error(err)
		return
	}

	authUser := sessionData.AuthUser

	authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authenticationroutes.Validator) {

		err = userdbcache.SetUserInfo(authUser, validator.LoginBodyReq.Username, validator.LoginBodyReq.FirstName, validator.LoginBodyReq.LastName, false)

		if err != nil {
			c.Error(err)
			return
		}

		//SendUserInfoUpdatedEmail(c, authUser)

		email := mailer.QueueEmail{
			Name: authUser.FirstName,
			To:   authUser.Email,
			//Token:     passwordlessToken,
			EmailType: mailer.QUEUE_EMAIL_TYPE_EMAIL_UPDATED,
			//Ttl:       fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		queue.PublishEmail(&email)
	})
}
