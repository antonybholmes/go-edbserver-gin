package admin

import (
	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	"github.com/gin-gonic/gin"
)

type UserListReq struct {
	Offset  uint
	Records uint
}

type UserStatResp struct {
	Users uint `json:"users"`
}

func UserStatsRoute(c *gin.Context) {

	var req UserListReq

	c.Bind(&req)

	users, err := userdbcache.NumUsers()

	if err != nil {
		c.Error(err)
		return
	}

	resp := UserStatResp{Users: users}

	web.MakeDataResp(c, "", resp)

}

func UsersRoute(c *gin.Context) {

	var req UserListReq

	c.Bind(&req)

	//log.Debug().Msgf("list users %v", req)

	users, err := userdbcache.Users(req.Records, req.Offset)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", users)

}

func RolesRoute(c *gin.Context) {

	roles, err := userdbcache.Roles()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", roles)

}

func GroupsRoute(c *gin.Context) {

	groups, err := userdbcache.Groups()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", groups)

}

func UpdateUserRoute(c *gin.Context) {

	authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().LoadAuthUserFromId().Success(func(validator *authenticationroutes.Validator) {

		//log.Debug().Msgf("roles here")

		//db, err := userdbcache.NewConn()

		// if err != nil {
		// 	return web.ErrorReq(err)
		// }

		//defer db.Close()

		//authUser, err := userdbcache.FindUserByPublicId(validator.Req.PublicId)

		// if err != nil {
		// 	return web.ErrorReq(err)
		// }

		authUser := validator.AuthUser

		err := userdbcache.SetUserInfo(authUser,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.FirstName,
			validator.UserBodyReq.LastName,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		err = userdbcache.SetEmailAddress(authUser,
			validator.Address,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		if validator.UserBodyReq.Password != "" {
			err = userdbcache.SetPassword(authUser,
				validator.UserBodyReq.Password)

			if err != nil {
				c.Error(err)

				return
			}
		}

		// set roles
		err = userdbcache.SetUserGroups(authUser,
			validator.UserBodyReq.Groups,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		web.MakeOkResp(c, "user updated")
	})
}

func AddUserRoute(c *gin.Context) {

	authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *authenticationroutes.Validator) {

		// assume email is not verified
		authUser, err := userdbcache.Instance().CreateUser(
			validator.UserBodyReq.Username,
			validator.Address,
			validator.UserBodyReq.Password,
			validator.UserBodyReq.FirstName,
			validator.UserBodyReq.LastName,
			validator.UserBodyReq.EmailIsVerified)

		if err != nil {
			c.Error(err)
			return
		}

		// tell user their account was created
		//go SendAccountCreatedEmail(authUser, validator.Address)

		email := mailserver.MailItem{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: edbmail.EmailQueueTypeAccountCreated,
			LinkUrl:   consts.AppUrl}
		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "account created email sent")
	})
}

func DeleteUserRoute(c *gin.Context) {
	publicId := c.Param("id")

	err := userdbcache.DeleteUser(publicId)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeOkResp(c, "user deleted")
}
