package admin

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, rulesMiddleware gin.HandlerFunc) {
	adminGroup := r.Group("/admin",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//middleware.JwtIsAdminMiddleware()
	)

	adminGroup.GET("/roles", RolesRoute)
	adminGroup.GET("/groups", GroupsRoute)

	adminUsersGroup := adminGroup.Group("/users")

	adminUsersGroup.POST("", UsersRoute)
	adminUsersGroup.GET("/stats", UserStatsRoute)
	adminUsersGroup.POST("/update", UpdateUserRoute)
	adminUsersGroup.POST("/add", AddUserRoute)
	adminUsersGroup.DELETE("/delete/:id", DeleteUserRoute)
}
