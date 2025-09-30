package authentication

import (
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, jwtUserMiddleWare gin.HandlerFunc, updateTokenMiddleware gin.HandlerFunc) {
	// Allow users to sign up for an account
	r.POST("/signup", SignupRoute)

	authGroup := r.Group("/auth")

	// auth0Group := authGroup.Group("/auth0")
	// auth0Group.POST("/validate",
	// 	jwtAuth0UserMiddleware,
	// 	auth0routes.ValidateAuth0TokenRoute)

	authGroup.POST("/signin", UsernamePasswordSignInRoute)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/verified",
		jwtUserMiddleWare,
		middleware.JwtIsVerifyEmailTokenMiddleware(),
		EmailAddressVerifiedRoute,
	)

	// with the correct token, performs the update
	emailGroup.POST("/reset",
		jwtUserMiddleWare,
		SendResetEmailEmailRoute)

	// with the correct token, performs the update
	emailGroup.POST("/update",
		jwtUserMiddleWare,
		UpdateEmailRoute)

	passwordGroup := authGroup.Group("/passwords")

	// sends a reset link
	passwordGroup.POST("/reset",
		SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update",
		jwtUserMiddleWare,
		UpdatePasswordRoute)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c *gin.Context) {
		PasswordlessSignInEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin",
		jwtUserMiddleWare,
		PasswordlessSignInRoute,
	)

	tokenGroup := authGroup.Group("/tokens", jwtUserMiddleWare)
	tokenGroup.POST("/info", TokenInfoRoute)
	tokenGroup.POST("/access", NewAccessTokenRoute)

	usersGroup := authGroup.Group("/users",
		jwtUserMiddleWare)

	//usersGroup.POST("", authorizationroutes.UserRoute)

	// we do not use an id parameter here as the user is derived from the token
	// for security reasons. User can be updated using an update token
	usersGroup.POST("/update", updateTokenMiddleware, UpdateUserRoute)

	//usersGroup.POST("/passwords/update", authentication.UpdatePasswordRoute)
}
