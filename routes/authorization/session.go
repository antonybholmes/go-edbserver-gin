package authorization

import (
	"github.com/antonybholmes/go-auth/userdbcache"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/gin-gonic/gin"
)

func SessionUpdateUserRoute(c *gin.Context) {
	sessionData, err := authenticationroutes.ReadSessionInfo(c)

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

		SendUserInfoUpdatedEmail(c, authUser)
	})
}
