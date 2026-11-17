package authentication

import (
	"github.com/antonybholmes/go-edbserver-gin/config"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	jwtUserMiddleWare gin.HandlerFunc,
	updateTokenMiddleware gin.HandlerFunc,
	cfg *config.AppConfig) {

	signupHandler := &SignupHandler{
		Config: cfg,
	}

	// Allow users to sign up for an account
	r.POST("/signup", signupHandler.SignupRoute)

	authGroup := r.Group("/auth")

	// auth0Group := authGroup.Group("/auth0")
	// auth0Group.POST("/validate",
	// 	jwtAuth0UserMiddleware,
	// 	auth0routes.ValidateAuth0TokenRoute)

	signinHandler := &SignInHandler{
		Config: cfg,
	}
	authGroup.POST("/signin", signinHandler.UsernamePasswordSignInRoute)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/verified",
		jwtUserMiddleWare,
		middleware.JwtIsVerifyEmailTokenMiddleware(),
		EmailAddressVerifiedRoute,
	)

	emailHandler := &EmailHandler{
		Config: cfg,
	}

	// with the correct token, performs the update
	emailGroup.POST("/reset",
		jwtUserMiddleWare,
		emailHandler.SendResetEmailEmailRoute)

	// with the correct token, performs the update
	emailGroup.POST("/update",
		jwtUserMiddleWare,
		emailHandler.UpdateEmailRoute)

	passwordGroup := authGroup.Group("/passwords")

	passwordHandler := &PasswordHandler{
		Config: cfg,
	}

	// sends a reset link
	passwordGroup.POST("/reset",
		passwordHandler.SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update",
		jwtUserMiddleWare,
		UpdatePasswordRoute)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c *gin.Context) {
		PasswordlessSignInEmailRoute(c, nil, cfg.Auth.ShortTtlMins)
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
