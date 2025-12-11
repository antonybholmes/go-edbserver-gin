package session

import (
	"context"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	"github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/oauth2"
	"github.com/antonybholmes/go-web/middleware"
	csrfmiddleware "github.com/antonybholmes/go-web/middleware/csrf"
	omw "github.com/antonybholmes/go-web/middleware/oauth2"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	otp *auth.OTP,

	jwtUserMiddleWare gin.HandlerFunc) {

	ctx := context.Background()

	auth0OIDCVerifer, err := oauth2.NewOIDCVerifier(ctx,
		consts.Auth0Domain,
		consts.Auth0Audience,
		consts.Auth0EmailClaim,
		consts.Auth0NameClaim)

	if err != nil {
		log.Fatal().Msgf("failed to create auth0 oidc verifier: %v", err)
	}

	// so we can verify cognito tokens
	congnitoOIDCVerifer, err := oauth2.NewStandardOIDCVerifier(ctx,
		consts.CognitoDomain,
		consts.CognitoAudience,
	)

	if err != nil {
		log.Fatal().Msgf("failed to create cognito oidc verifier: %v", err)
	}

	clerkOIDCVerifer, err := oauth2.NewStandardOIDCVerifier(ctx,
		consts.ClerkDomain,
		consts.ClerkAudience,
	)

	if err != nil {
		log.Fatal().Msgf("failed to create clerk oidc verifier: %v", err)
	}

	// supabaseOIDCVerifer, err := oauth2.NewStandardOIDCVerifier(ctx,
	// 	consts.SupabaseDomain,
	// 	consts.SupabaseAudience,
	// )

	// if err != nil {
	// 	log.Fatal().Msgf("failed to create supabase oidc verifier: %v", err)
	// }

	otpRoutes := authentication.NewOTPRoutes(otp)

	sessionRoutes := NewSessionRoutes(otpRoutes)

	sessionMiddleware := middleware.SessionIsValidMiddleware()

	//jwtAuth0Middleware2 := omw.JwtAuth0Middleware(consts.JwtAuth0RsaPublicKey)
	jwtAuth0Middleware := omw.JwtOIDCMiddleware(auth0OIDCVerifer)

	jwtCognitoMiddleware := omw.JwtOIDCMiddleware(congnitoOIDCVerifer)

	//jwtClerkMiddleware := omw.JwtClerkMiddleware(consts.JwtClerkRsaPublicKey)
	jwtClerkMiddleware := omw.JwtOIDCMiddleware(clerkOIDCVerifer)

	//jwtSupabaseMiddleware := omw.JwtOIDCMiddleware(supabaseOIDCVerifer)
	jwtSupabaseMiddleware := omw.JwtSupabaseMiddleware(consts.SupabaseJwtSecretKey)

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

	// support both spellings
	sessionGroup.POST("/signout",
		//csrfMiddleware,
		sessionMiddleware,
		SessionSignOutRoute)
	sessionGroup.POST("/sign-out",
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
