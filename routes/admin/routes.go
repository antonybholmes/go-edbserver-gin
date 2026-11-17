package admin

import (
	"github.com/antonybholmes/go-edbserver-gin/config"
	"github.com/antonybholmes/go-web/auth/userdb"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	rulesMiddleware gin.HandlerFunc,
	userDB userdb.UserDB,
	config *config.AppConfig) {
	adminGroup := r.Group("/admin",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//middleware.JwtIsAdminMiddleware()
	)

	adminHandlers := &AdminHandlers{
		UserDB: userDB,
		Config: config,
	}

	adminGroup.GET("/roles", adminHandlers.RolesRoute)
	adminGroup.GET("/groups", adminHandlers.GroupsRoute)

	adminUsersGroup := adminGroup.Group("/users")

	adminUsersGroup.POST("", adminHandlers.UsersRoute)
	adminUsersGroup.GET("/stats", adminHandlers.UserStatsRoute)
	adminUsersGroup.POST("/update", adminHandlers.UpdateUserRoute)
	adminUsersGroup.POST("/add", adminHandlers.AddUserRoute)
	adminUsersGroup.DELETE("/:id/delete", adminHandlers.DeleteUserRoute)
}
