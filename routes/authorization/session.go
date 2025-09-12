package authorization

import (
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/antonybholmes/go-web/userdbcache"
	"github.com/gin-contrib/sessions"

	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/gin-gonic/gin"
)

// type UserUpdateReq struct {
// 	Username  string `json:"password"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// }

type PasswordUpdateReq struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func SessionUpdateUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	// Read the session info, which includes the AuthUser

	sessionData, err := middleware.ReadSessionInfo(c, session)

	if err != nil {
		c.Error(err)
		return
	}

	authUser := sessionData.AuthUser

	authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().Success(func(validator *authenticationroutes.Validator) {

		err = userdbcache.SetUserInfo(authUser,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.FirstName,
			validator.UserBodyReq.LastName,
			false)

		if err != nil {
			c.Error(err)
			return
		}

		//SendUserInfoUpdatedEmail(c, authUser)

		email := mailserver.MailItem{
			Name: authUser.FirstName,
			To:   authUser.Email,
			//Token:     passwordlessToken,
			EmailType: mailserver.QUEUE_EMAIL_TYPE_EMAIL_UPDATED,
			//Ttl:       fmt.Sprintf("%d minutes", int(consts.PASSWORDLESS_TOKEN_TTL_MINS.Minutes())),
			//LinkUrl:   consts.URL_SIGN_IN,
			//VisitUrl:    validator.Req.VisitUrl
		}

		mailqueue.SendMail(&email)
	})
}

func SessionUpdatePasswordRoute(c *gin.Context) {
	user, _ := c.Get(web.SESSION_USER)

	authUser := user.(*auth.AuthUser)

	// use current session user
	authUser, err := userdbcache.FindUserByPublicId(authUser.PublicId)

	if err != nil {
		web.BadReqResp(c, "user not found")
		return
	}

	var req PasswordUpdateReq

	err = c.ShouldBindJSON(&req)

	if err != nil {
		web.BadReqResp(c, "invalid request body")
		return
	}

	err = auth.CheckPassword(req.Password)

	if err != nil {
		web.BadReqResp(c, "password does not meet requirements")
		return
	}

	err = authUser.CheckPasswordsMatch(req.Password)

	if err != nil {
		web.BadReqResp(c, "current and new password do not match")
		return
	}

	err = userdbcache.SetPassword(authUser, req.NewPassword)

	if err != nil {
		web.BadReqResp(c, "could not update password")
		return
	}

	web.MakeOkResp(c, "password updated successfully")
}
