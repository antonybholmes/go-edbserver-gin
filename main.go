package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"net/http"
	"os"
	"runtime"

	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	adminroutes "github.com/antonybholmes/go-edbserver-gin/routes/admin"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	authorizationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authorization"
	"github.com/antonybholmes/go-edbserver-gin/routes/modules"
	"github.com/antonybholmes/go-hubs/hubsdbcache"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/access"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/antonybholmes/go-web/userdbcache"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/antonybholmes/go-web/middleware"

	utilsroutes "github.com/antonybholmes/go-edbserver-gin/routes/utils"
	"github.com/antonybholmes/go-geneconv/geneconvdbcache"
	"github.com/antonybholmes/go-genome/genomedbcache"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
	"github.com/antonybholmes/go-scrna/scrnadbcache"
	"github.com/antonybholmes/go-seqs/seqsdbcache"
	"github.com/antonybholmes/go-sys/env"
	_ "github.com/mattn/go-sqlite3"
)

type AboutResp struct {
	Name      string `json:"name"`
	Copyright string `json:"copyright"`
	Version   string `json:"version"`
	Updated   string `json:"updated"`
}

type InfoResp struct {
	IpAddr string `json:"ipAddr"`
	Arch   string `json:"arch"`
}

// var store *sqlitestorr.SqliteStore
var store cookie.Store

var rdb *redis.Client

var re *access.RuleEngine

const PREFLIGHT_MAX_AGE = 12 * 3600 // 12 hours

func initLogger() {
	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", consts.AppName),
		MaxSize:    10,   // Max size in MB before rotating
		MaxBackups: 3,    // Keep 3 backup files
		MaxAge:     7,    // Retain files for 7 days
		Compress:   true, // Compress old log files
	}

	multiWriter := io.MultiWriter(os.Stderr, fileLogger)

	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

	// we use != development because it means we need to set the env variable in order
	// to see debugging work. The default is to assume production, in which case we use
	// lumberjack
	if os.Getenv("APP_ENV") != "development" {
		//zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
	}

	log.Logger = logger
}

func init() {
	initLogger()

	env.Ls()

	// store = sys.Must(sqlitestorr.NewSqliteStore("data/users.db",
	// 	"sessions",
	// 	"/",
	// 	auth.MAX_AGE_7_DAYS_SECS,
	// 	[]byte(consts.SESSION_SECRET)))

	// follow https://github.com/gorilla/sessions/blob/main/storr.go#L55
	// SESSION_KEY should be 64 bytes/chars and SESSION_ENCRYPTION_KEY should be 32 bytes/chars

	// storr.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   auth.MAX_AGE_7_DAYS_SECS,
	// 	HttpOnly: false,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode}

	userdbcache.InitCache() //"data/users.db")

	//mailserver.Init()

	dnadbcache.InitCache("data/modules/dna")
	genomedbcache.InitCache("data/modules/genome")
	//microarraydb.InitDB("data/microarray")

	gexdbcache.InitCache("data/modules/gex")

	scrnadbcache.InitCache("data/modules/scrna")

	mutationdbcache.InitCache("data/modules/mutations")

	geneconvdbcache.InitCache("data/modules/geneconv/geneconv.db")

	motifsdb.InitCache("data/modules/motifs/motifs.db")

	pathwaydbcache.InitCache("data/modules/pathway/pathway-v2.db")

	seqsdbcache.InitCache("data/modules/seqs/")

	cytobandsdbcache.InitCache("data/modules/cytobands/")

	bedsdbcache.InitCache("data/modules/beds/")

	hubsdbcache.InitCache("data/modules/hubs/")

	rdb = redis.NewClient(&redis.Options{
		Addr:     consts.RedisAddr,
		Username: "edb",
		Password: consts.RedisPassword,
		DB:       0, // use default DB
	})

	//queue.Init(mailserver.NewRedisEmailQueue(rdb))

	mailqueue.Init(mailserver.NewSQSEmailQueue(consts.SqsQueueUrl))

	re = access.NewRuleEngine()

	re.LoadRules("config/access-rules.json")

	// writer := kafka.NewWriter(kafka.WriterConfig{
	// 	Brokers:  []string{"localhost:9094"}, // Kafka broker
	// 	Topic:    mailserver.QUEUE_EMAIL_CHANNEL, // Topic name
	// 	Balancer: &kafka.LeastBytes{},        // Balancer (optional)
	// })

	// queue.Init(mailserver.NewKafkaEmailPublisher(writer))

}

func main() {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	//consts.Init()

	tokengen.Init(consts.JwtRsaPrivateKey)

	//env.Load()

	// list env to see what is loaded
	//env.Ls()

	//initCache()

	// test redis
	//email := gomailserver.QueueEmail{To: "antony@antonybholmes.dev"}
	//rdb.PublishEmail(&email)

	//
	// Set logging to file
	//

	// all subsequent middleware is reliant on this to function
	claimsParser := middleware.JwtClaimsRSAParser(consts.JwtRsaPublicKey)
	jwtUserMiddleWare := middleware.UserJWTMiddleware(claimsParser)

	jwtAuth0Middleware := middleware.JwtAuth0Middleware(consts.JwtAuth0RsaPublicKey)

	jwtClerkMiddleware := middleware.JwtClerkMiddleware(consts.JwtClerkRsaPublicKey)

	jwtSupabaseMiddleware := middleware.JwtSupabaseMiddleware(consts.JwtSupabaseSecretKey)

	csrfMiddleware := middleware.CSRFValidateMiddleware()

	sessionMiddleware := middleware.SessionIsValidMiddleware()

	//accessTokenMiddleware := middleware.JwtIsAccessTokenMiddleware()

	rulesMiddleware := middleware.RulesMiddleware(claimsParser, re)

	updateTokenMiddleware := middleware.JwtIsUpdateTokenMiddleware()

	//rdfRoleMiddleware := middleware.JwtHasRDFRoleMiddleware()

	otp := auth.NewDefaultOTP(rdb)

	otpRoutes := authenticationroutes.NewOTPRoutes(otp)

	sessionRoutes := authenticationroutes.NewSessionRoutes(otpRoutes)

	// Setup tracer provider
	tp, err := initTracerProvider()
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer provider: %v", err)
	}
	defer func() {
		// Shutdown flushes any remaining spans
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal().Msgf("failed to shutdown tracer provider: %v", err)
		}
	}()

	//r := gin.Default()
	r := gin.New()

	// Add OpenTelemetry middleware to Gin router
	r.Use(otelgin.Middleware("edb-server"))

	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	//r.Use(middleware.CSRFCookieMiddleware())

	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8000",
			"https://edb.rdf-lab.org",
			"https://edb-client-astro.pages.dev",
			"https://edb-client-next.pages.dev",
			"https://edb-client-next.vercel.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Set-Cookie"},
		// for sharing session cookie for validating logins etc
		AllowCredentials: true,              // Allow credentials (cookies, HTTP authentication)
		MaxAge:           PREFLIGHT_MAX_AGE, // Cache preflight response for 12 hours
	}))

	store = cookie.NewStore([]byte(consts.SessionKey),
		[]byte(consts.SessionEncryptionKey))
	r.Use(sessions.Sessions(consts.SessionName, store))

	r.GET("/about", func(c *gin.Context) {
		fmt.Println("Handler:", c.FullPath())
		c.JSON(http.StatusOK,
			AboutResp{
				Name:      consts.AppName,
				Version:   consts.Version,
				Updated:   consts.Updated,
				Copyright: consts.Copyright})
	})

	r.GET("/info", func(c *gin.Context) {
		web.MakeDataResp(c, "", InfoResp{
			Arch:   runtime.GOARCH,
			IpAddr: c.ClientIP()})
	})

	//
	// Routes
	//

	adminGroup := r.Group("/admin",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//middleware.JwtIsAdminMiddleware()
	)

	adminGroup.GET("/roles", adminroutes.RolesRoute)

	adminUsersGroup := adminGroup.Group("/users")

	adminUsersGroup.POST("", adminroutes.UsersRoute)
	adminUsersGroup.GET("/stats", adminroutes.UserStatsRoute)
	adminUsersGroup.POST("/update", adminroutes.UpdateUserRoute)
	adminUsersGroup.POST("/add", adminroutes.AddUserRoute)
	adminUsersGroup.DELETE("/delete/:id", adminroutes.DeleteUserRoute)

	// Allow users to sign up for an account
	r.POST("/signup", authenticationroutes.SignupRoute)

	//
	// user groups: start
	//

	authGroup := r.Group("/auth")

	// auth0Group := authGroup.Group("/auth0")
	// auth0Group.POST("/validate",
	// 	jwtAuth0UserMiddleware,
	// 	auth0routes.ValidateAuth0TokenRoute)

	authGroup.POST("/signin", authenticationroutes.UsernamePasswordSignInRoute)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/verified",
		jwtUserMiddleWare,
		middleware.JwtIsVerifyEmailTokenMiddleware(),
		authenticationroutes.EmailAddressVerifiedRoute,
	)

	// with the correct token, performs the update
	emailGroup.POST("/reset",
		jwtUserMiddleWare,
		authenticationroutes.SendResetEmailEmailRoute)

	// with the correct token, performs the update
	emailGroup.POST("/update",
		jwtUserMiddleWare,
		authenticationroutes.UpdateEmailRoute)

	passwordGroup := authGroup.Group("/passwords")

	// sends a reset link
	passwordGroup.POST("/reset",
		authenticationroutes.SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update",
		jwtUserMiddleWare,
		authenticationroutes.UpdatePasswordRoute)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c *gin.Context) {
		authenticationroutes.PasswordlessSignInEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin",
		jwtUserMiddleWare,
		authenticationroutes.PasswordlessSignInRoute,
	)

	tokenGroup := authGroup.Group("/tokens", jwtUserMiddleWare)
	tokenGroup.POST("/info", authorizationroutes.TokenInfoRoute)
	tokenGroup.POST("/access", authorizationroutes.NewAccessTokenRoute)

	usersGroup := authGroup.Group("/users",
		jwtUserMiddleWare)

	//usersGroup.POST("", authorizationroutes.UserRoute)

	// we do not use an id parameter here as the user is derived from the token
	// for security reasons. User can be updated using an update token
	usersGroup.POST("/update", updateTokenMiddleware, authorizationroutes.UpdateUserRoute)

	//usersGroup.POST("/passwords/update", authentication.UpdatePasswordRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := r.Group("/sessions")

	sessionAuthGroup := sessionGroup.Group("/auth")

	sessionOAuth2Group := sessionAuthGroup.Group("/oauth2")

	sessionOAuth2Group.POST("/auth0/signin",
		jwtAuth0Middleware,
		sessionRoutes.SessionSignInUsingAuth0Route)

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
		authenticationroutes.SessionSignOutRoute)

	sessionTokensGroup := sessionGroup.Group("/tokens",
		csrfMiddleware,
		sessionMiddleware)

	//sessionTokensGroup.POST("/access",
	//		authenticationroutes.NewAccessTokenFromSessionRoute)

	// issue tokens
	sessionTokensGroup.POST("/create/:type",
		authenticationroutes.CreateTokenFromSessionRoute)

	// update session info
	sessionGroup.POST("/refresh",
		csrfMiddleware,
		sessionMiddleware,
		sessionRoutes.SessionRefreshRoute)

	sessionUserGroup := sessionGroup.Group("/user",
		csrfMiddleware,
		sessionMiddleware)
	sessionUserGroup.GET("", authenticationroutes.UserFromSessionRoute)
	sessionUserGroup.POST("/update",
		authorizationroutes.SessionUpdateUserRoute)

	sessionUserGroup.POST("/passwords/update",
		authorizationroutes.SessionUpdatePasswordRoute)

	modules.RegisterRoutes(r, rulesMiddleware)

	//
	// Util routes
	//

	utilsGroup := r.Group("/utils")
	//moduleGroup.Use(jwtMiddleWare,JwtIsAccessTokenMiddleware)

	xlsxGroup := utilsGroup.Group("/xlsx")
	xlsxGroup.POST("/sheets", utilsroutes.XlsxSheetsRoute)
	xlsxGroup.POST("/to/:format", utilsroutes.XlsxToRoute)

	utilsGroup.GET("/passwords/hash", utilsroutes.HashedPasswordRoute)
	utilsGroup.GET("/randkey", utilsroutes.RandomKeyRoute)
	utilsGroup.GET("/uuidv7", utilsroutes.UUIDv7Route)

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// required so it can listen externally within docker container
	r.Run("0.0.0.0:8080") //"localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
