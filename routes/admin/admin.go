package admin

import (
	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	"github.com/antonybholmes/go-edbserver-gin/config"

	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth/userdb"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
)

type (
	UserListReq struct {
		Offset  int
		Records int
	}

	UserStatResp struct {
		Users int `json:"users"`
	}

	AdminHandlers struct {
		UserDB userdb.UserDB
		Config *config.AppConfig
	}
)

func (h *AdminHandlers) UserStatsRoute(c *gin.Context) {

	var req UserListReq

	c.Bind(&req)

	users, err := h.UserDB.NumUsers()

	if err != nil {
		c.Error(err)
		return
	}

	resp := UserStatResp{Users: users}

	web.MakeDataResp(c, "", resp)

}

func (h *AdminHandlers) UsersRoute(c *gin.Context) {

	var req UserListReq

	c.Bind(&req)

	//log.Debug().Msgf("list users %v", req)

	users, err := h.UserDB.Users(req.Records, req.Offset)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", users)

}

func (h *AdminHandlers) RolesRoute(c *gin.Context) {

	roles, err := h.UserDB.Roles()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", roles)

}

func (h *AdminHandlers) GroupsRoute(c *gin.Context) {

	groups, err := h.UserDB.Groups()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", groups)

}

func (h *AdminHandlers) UpdateUserRoute(c *gin.Context) {

	middleware.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().LoadAuthUserFromId().Success(func(validator *middleware.Validator) {

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

		err := h.UserDB.SetUserInfo(authUser,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.Name,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		err = h.UserDB.SetEmailAddress(authUser,
			validator.Address,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		if validator.UserBodyReq.PictureUrl != "" {
			err = h.UserDB.SetPictureUrl(authUser,
				validator.UserBodyReq.PictureUrl,
				true)
		}

		if err != nil {
			c.Error(err)
			return
		}

		if validator.UserBodyReq.Password != "" {
			_, err = h.UserDB.SetPassword(authUser,
				validator.UserBodyReq.Password, true)

			if err != nil {
				c.Error(err)

				return
			}
		}

		// set roles
		err = h.UserDB.SetUserGroups(authUser,
			validator.UserBodyReq.Groups,
			true)

		if err != nil {
			c.Error(err)
			return
		}

		web.MakeOkResp(c, "user updated")
	})
}

func (h *AdminHandlers) AddUserRoute(c *gin.Context) {

	middleware.NewValidator(c).CheckUsernameIsWellFormed().CheckEmailIsWellFormed().Success(func(validator *middleware.Validator) {

		// assume email is not verified
		authUser, err := h.UserDB.CreateUser(
			validator.Address,
			validator.UserBodyReq.Username,
			validator.UserBodyReq.Password,
			validator.UserBodyReq.Name,
			validator.UserBodyReq.PictureUrl,
			validator.UserBodyReq.EmailIsVerified,
			"edb",
			false)

		if err != nil {
			c.Error(err)
			return
		}

		// tell user their account was created
		//go SendAccountCreatedEmail(authUser, validator.Address)

		email := mailserver.MailItem{
			Name:      authUser.Name,
			To:        authUser.Email,
			EmailType: edbmail.EmailQueueTypeAccountCreated,
			LinkUrl:   h.Config.App.URL}

		mailqueue.SendMail(&email)

		web.MakeOkResp(c, "account created email sent")
	})
}

func (h *AdminHandlers) DeleteUserRoute(c *gin.Context) {
	publicId := c.Param("id")

	err := h.UserDB.DeleteUser(publicId)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeOkResp(c, "user deleted")
}
