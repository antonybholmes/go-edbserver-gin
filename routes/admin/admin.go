package admin

import (
	"github.com/antonybholmes/go-edb-server-gin/consts"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/userdbcache"
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

func UpdateUserRoute(c *gin.Context) {

	authenticationroutes.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().LoadAuthUserFromUuid().Success(func(validator *authenticationroutes.Validator) {

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
			validator.LoginBodyReq.Username,
			validator.LoginBodyReq.FirstName,
			validator.LoginBodyReq.LastName,
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

		if validator.LoginBodyReq.Password != "" {
			err = userdbcache.SetPassword(authUser,
				validator.LoginBodyReq.Password)

			if err != nil {
				c.Error(err)

				return
			}
		}

		// set roles
		err = userdbcache.SetUserRoles(authUser,
			validator.LoginBodyReq.Roles,
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
			validator.LoginBodyReq.Username,
			validator.Address,
			validator.LoginBodyReq.Password,
			validator.LoginBodyReq.FirstName,
			validator.LoginBodyReq.LastName,
			validator.LoginBodyReq.EmailIsVerified)

		if err != nil {
			c.Error(err)
			return
		}

		// tell user their account was created
		//go SendAccountCreatedEmail(authUser, validator.Address)

		email := mailer.QueueEmail{
			Name:      authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.QUEUE_EMAIL_TYPE_ACCOUNT_CREATED,
			LinkUrl:   consts.APP_URL}
		queue.PublishEmail(&email)

		web.MakeOkResp(c, "account created email sent")
	})
}

func DeleteUserRoute(c *gin.Context) {
	uuid := c.Param("uuid")

	err := userdbcache.DeleteUser(uuid)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeOkResp(c, "user deleted")
}
