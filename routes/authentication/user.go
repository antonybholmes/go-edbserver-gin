package authentication

import (
	email "github.com/antonybholmes/go-mail"
	"github.com/antonybholmes/go-mail/mailqueue"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/antonybholmes/go-web/middleware"

	"github.com/gin-gonic/gin"
)

type NameReq struct {
	Name string `json:"name"`
}

func UserUpdatedResp(c *gin.Context) {
	web.MakeOkResp(c, "account updated")
}

func UpdateUserRoute(c *gin.Context) {

	validator, err := middleware.NewValidator(c).ParseSignInRequestBody().LoadTokenClaims().Ok()

	userId := validator.Claims.Subject

	log.Debug().Msgf("UpdateUserRoute: publicId: %s ", userId)

	authUser, err := userdbcache.Instance().FindUserById(userId)

	if err != nil {
		c.Error(err)
		return
	}

	err = userdbcache.Instance().SetUserInfo(authUser,
		validator.UserBodyReq.Username,
		validator.UserBodyReq.Name,
		false)

	if err != nil {
		c.Error(err)
		return
	}

	//return SendUserInfoUpdatedEmail(c, authUser)

	// reload user details
	authUser, err = userdbcache.Instance().FindUserById(userId)

	if err != nil {
		c.Error(err)
		return
	}

	// send email notification of change
	email := email.MailItem{Name: authUser.Name,
		To:        authUser.Email,
		EmailType: email.EmailTypeAccountUpdated}
	mailqueue.Instance().SendMail(c.Request.Context(), &email)

	// send back updated user to having to do a separate call to get the new data
	web.MakeDataResp(c, "account updated confirmation email sent", authUser)

}

func UserRoute(c *gin.Context) {
	middleware.NewValidator(c).
		LoadAuthUserFromToken().
		Success(func(validator *middleware.Validator) {
			web.MakeDataResp(c, "", validator.AuthUser)
		})
}
