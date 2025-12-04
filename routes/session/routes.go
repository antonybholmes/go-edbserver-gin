package session

import (
	"github.com/antonybholmes/go-edbserver-gin/consts"
	"github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/oauth2"
	"github.com/antonybholmes/go-web/middleware"
	csrfmiddleware "github.com/antonybholmes/go-web/middleware/csrf"
	omw "github.com/antonybholmes/go-web/middleware/oauth2"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	otp *auth.OTP,
	oidcVerifer *oauth2.OIDCVerifier,
	jwtUserMiddleWare gin.HandlerFunc) {

	otpRoutes := authentication.NewOTPRoutes(otp)

	sessionRoutes := NewSessionRoutes(otpRoutes)

	sessionMiddleware := middleware.SessionIsValidMiddleware()

	jwtAuth0Middleware := omw.JwtAuth0Middleware(consts.JwtAuth0RsaPublicKey)

	jwtClerkMiddleware := omw.JwtClerkMiddleware(consts.JwtClerkRsaPublicKey)

	jwtSupabaseMiddleware := omw.JwtSupabaseMiddleware(consts.JwtSupabaseSecretKey)

	jwtCognitoMiddleware := omw.JwtOIDCMiddleware(oidcVerifer)

	csrfMiddleware := csrfmiddleware.CSRFValidateMiddleware()

	sessionGroup := r.Group("/sessions")

	sessionAuthGroup := sessionGroup.Group("/auth")

	sessionOAuth2Group := sessionAuthGroup.Group("/oauth2")

	sessionOAuth2Group.POST("/auth0/signin",
		jwtAuth0Middleware,
		sessionRoutes.SessionSignInUsingAuth0Route)

	sessionOAuth2Group.POST("/cognito/signin",
		jwtCognitoMiddleware,
		sessionRoutes.SessionSignInUsingCognitoRoute)

	sessionOAuth2Group.POST("/clerk/signin",
		jwtClerkMiddleware,
		sessionRoutes.SessionSignInUsingClerkRoute)

	sessionOAuth2Group.POST("/supabase/signin",
		jwtSupabaseMiddleware,
		sessionRoutes.SessionSignInUsingSupabaseRoute)

	sessionAuthGroup.POST("/signin",
		sessionRoutes.SessionUsernamePasswordSignInRoute)

	sessionOtpGroup := sessionAuthGroup.Group("/otp")

	sessionOtpGroup.POST("/send", sessionRoutes.SessionEmailOTPRoute)

	sessionOtpGroup.POST("/signin",
		sessionRoutes.SessionSignInUsingEmailAndOTPRoute)

	sessionAuthGroup.POST("/passwordless/validate",
		jwtUserMiddleWare,
		sessionRoutes.SessionPasswordlessValidateSignInRoute)

	sessionGroup.POST("/api/keys/signin",
		sessionRoutes.SessionApiKeySignInRoute)

	sessionGroup.GET("/info",
		sessionMiddleware,
		sessionRoutes.SessionInfoRoute)

	sessionGroup.POST("/csrf-token/refresh",
		sessionRoutes.SessionNewCSRFTokenRoute)

	sessionGroup.POST("/signout",
		//csrfMiddleware,
		sessionMiddleware,
		SessionSignOutRoute)

	sessionTokensGroup := sessionGroup.Group("/tokens",
		csrfMiddleware,
		sessionMiddleware)

	//sessionTokensGroup.POST("/access",
	//		authenticationroutes.NewAccessTokenFromSessionRoute)

	// issue tokens
	sessionTokensGroup.POST("/create/:type",
		CreateTokenFromSessionRoute)

	// update session info
	sessionGroup.POST("/refresh",
		csrfMiddleware,
		sessionMiddleware,
		sessionRoutes.SessionRefreshRoute)

	sessionUserGroup := sessionGroup.Group("/user",
		csrfMiddleware,
		sessionMiddleware)
	sessionUserGroup.GET("", UserFromSessionRoute)
	sessionUserGroup.POST("/update",
		SessionUpdateUserRoute)

	sessionUserGroup.POST("/passwords/update",
		SessionUpdatePasswordRoute)
}
