package main

import (
	"fmt"
	"io"

	"net/http"
	"os"
	"runtime"

	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edbserver-gin/consts"
	adminroutes "github.com/antonybholmes/go-edbserver-gin/routes/admin"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	authorizationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authorization"
	bedroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/beds"
	cytobandroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/cytobands"
	dnaroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/dna"
	geneconvroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/geneconv"
	genomeroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/genome"
	gexroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/gex"
	hubroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/hubs"
	motifroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/motifs"
	mutationroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/mutation"
	pathwayroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/pathway"
	scrnaroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/scrna"
	seqroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/seqs"
	"github.com/antonybholmes/go-hubs/hubsdbcache"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-web"
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
	"github.com/antonybholmes/go-mailserver/queue"
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

func initLogger() {
	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", consts.APP_NAME),
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
		Addr:     consts.REDIS_ADDR,
		Username: "edb",
		Password: consts.REDIS_PASSWORD,
		DB:       0, // use default DB
	})

	//queue.Init(mailserver.NewRedisEmailQueue(rdb))

	queue.Init(mailserver.NewSQSEmailQueue(consts.SQS_QUEUE_URL))

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

	tokengen.Init(consts.JWT_RSA_PRIVATE_KEY)

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
	jwtUserMiddleWare := middleware.JwtUserMiddleware(
		middleware.JwtClaimsRSAParser(consts.JWT_RSA_PUBLIC_KEY))

	jwtAuth0Middleware := middleware.JwtAuth0Middleware(consts.JWT_AUTH0_RSA_PUBLIC_KEY)

	jwtClerkMiddleware := middleware.JwtClerkMiddleware(consts.JWT_CLERK_RSA_PUBLIC_KEY)

	jwtSupabaseMiddleware := middleware.JwtSupabaseMiddleware(consts.JWT_SUPABASE_SECRET_KEY)

	csrfMiddleware := middleware.CSRFValidateMiddleware()

	sessionMiddleware := middleware.SessionIsValidMiddleware()

	accessTokenMiddleware := middleware.JwtIsAccessTokenMiddleware()

	updateTokenMiddleware := middleware.JwtIsUpdateTokenMiddleware()

	rdfRoleMiddleware := middleware.JwtHasRDFRoleMiddleware()

	otp := authenticationroutes.NewOTP(rdb)

	otpRoutes := authenticationroutes.NewOTPRoutes(otp)

	sessionRoutes := authenticationroutes.NewSessionRoutes(otpRoutes)

	//r := gin.Default()
	r := gin.New()

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
		AllowCredentials: true,      // Allow credentials (cookies, HTTP authentication)
		MaxAge:           12 * 3600, // Cache preflight response for 12 hours
	}))

	store = cookie.NewStore([]byte(consts.SESSION_KEY),
		[]byte(consts.SESSION_ENCRYPTION_KEY))
	r.Use(sessions.Sessions(consts.SESSION_NAME, store))

	r.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			AboutResp{
				Name:      consts.APP_NAME,
				Version:   consts.VERSION,
				Updated:   consts.UPDATED,
				Copyright: consts.COPYRIGHT})
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
		jwtUserMiddleWare,
		accessTokenMiddleware,
		middleware.JwtIsAdminMiddleware())

	adminGroup.GET("/roles", adminroutes.RolesRoute)

	adminUsersGroup := adminGroup.Group("/users")

	adminUsersGroup.POST("", adminroutes.UsersRoute)
	adminUsersGroup.GET("/stats", adminroutes.UserStatsRoute)
	adminUsersGroup.POST("/update", adminroutes.UpdateUserRoute)
	adminUsersGroup.POST("/add", adminroutes.AddUserRoute)
	adminUsersGroup.DELETE("/delete/:uuid", adminroutes.DeleteUserRoute)

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
		authenticationroutes.PasswordlessSigninEmailRoute(c, nil)
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

	sessionGroup.GET("/csrf-token",
		sessionRoutes.SessionCsrfTokenRoute)

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

	//
	// module groups: start
	//

	moduleGroup := r.Group("/modules")
	//moduleGroup.Use(jwtMiddleWare,JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")
	dnaGroup.POST("/:assembly", dnaroutes.DNARoute)
	dnaGroup.GET("/genomes", dnaroutes.GenomesRoute)

	genomeGroup := moduleGroup.Group("/genome")
	genomeGroup.GET("/genomes", genomeroutes.GenomesRoute)
	genomeGroup.POST("/within/:assembly", genomeroutes.WithinGenesRoute)
	genomeGroup.POST("/closest/:assembly", genomeroutes.ClosestGeneRoute)
	genomeGroup.POST("/annotate/:assembly", genomeroutes.AnnotateRoute)
	genomeGroup.POST("/overlap/:assembly", genomeroutes.OverlappingGenesRoute)
	genomeGroup.GET("/info/:assembly", genomeroutes.SearchForGeneByNameRoute)

	// mutationsGroup := moduleGroup.Group("/mutations",
	// 	jwtMiddleWare,
	// 	JwtIsAccessTokenMiddleware,
	// 	NewJwtPermissionsMiddleware("rdf"))

	mutationsGroup := moduleGroup.Group("/mutations")
	mutationsGroup.GET("/datasets/:assembly",
		mutationroutes.MutationDatasetsRoute)
	mutationsGroup.POST("/:assembly/:name",
		mutationroutes.MutationsRoute)
	mutationsGroup.POST("/maf/:assembly",
		mutationroutes.PileupRoute)

	mutationsGroup.POST("/pileup/:assembly",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		mutationroutes.PileupRoute,
	)

	gexGroup := moduleGroup.Group("/gex")
	gexGroup.GET("/species", gexroutes.SpeciesRoute)
	gexGroup.GET("/technologies", gexroutes.TechnologiesRoute)
	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)

	gexGroup.GET("/datasets/:species/:technology",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		gexroutes.GexDatasetsRoute)

	gexGroup.POST("/exp",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		gexroutes.GexGeneExpRoute,
	)

	scrnaGroup := moduleGroup.Group("/scrna")
	scrnaGroup.GET("/species", scrnaroutes.ScrnaSpeciesRoute)
	scrnaGroup.GET("/assemblies/:species", scrnaroutes.ScrnaAssembliesRoute)
	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)

	scrnaGroup.GET("/datasets/:species/:assembly",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		scrnaroutes.ScrnaDatasetsRoute)

	// scrnaGroup.GET("/clusters/:id",
	// 	jwtUserMiddleWare,
	// 	accessTokenMiddleware,
	// 	rdfRoleMiddleware,
	// 	scrnaroutes.ScrnaClustersRoute,
	// )

	scrnaGroup.GET("/metadata/:id",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		scrnaroutes.ScrnaMetadataRoute,
	)

	scrnaGroup.GET("/genes/:id",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		scrnaroutes.ScrnaGenesRoute,
	)

	scrnaGroup.GET("/genes/search/:id",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		scrnaroutes.ScrnaSearchGenesRoute,
	)

	scrnaGroup.POST("/gex/:id",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		scrnaroutes.ScrnaGexRoute,
	)

	hubsGroup := moduleGroup.Group("/hubs")
	hubsGroup.GET("/:assembly",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware,
		hubroutes.HubsRoute,
	)

	geneConvGroup := moduleGroup.Group("/geneconv")
	geneConvGroup.POST("/convert/:from/:to", geneconvroutes.ConvertRoute)

	// geneConvGroup.POST("/:species", func(c *gin.Context) {
	// 	return geneconvroutes.GeneInfoRoute(c, "")
	// })

	motifsGroup := moduleGroup.Group("/motifs")
	motifsGroup.GET("/datasets", motifroutes.DatasetsRoute)
	motifsGroup.POST("/search", motifroutes.SearchRoute)

	pathwayGroup := moduleGroup.Group("/pathway")
	pathwayGroup.GET("/genes", pathwayroutes.GenesRoute)
	pathwayGroup.POST("/dataset", pathwayroutes.DatasetRoute)
	pathwayGroup.GET("/datasets", pathwayroutes.DatasetsRoute)
	pathwayGroup.POST("/overlap", pathwayroutes.PathwayOverlapRoute)

	seqsGroup := moduleGroup.Group("/seqs",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware)

	seqsGroup.GET("/genomes", seqroutes.GenomeRoute)
	seqsGroup.GET("/platforms/:assembly", seqroutes.PlatformRoute)
	//tracksGroup.GET("/:platform/:assembly/tracks", seqroutes.TracksRoute)
	seqsGroup.GET("/search/:assembly", seqroutes.SearchSeqRoute)
	seqsGroup.POST("/bins", seqroutes.BinsRoute)

	cytobandsGroup := moduleGroup.Group("/cytobands")
	cytobandsGroup.GET("/:assembly/:chr", cytobandroutes.CytobandsRoute)

	bedsGroup := moduleGroup.Group("/beds",
		jwtUserMiddleWare,
		accessTokenMiddleware,
		rdfRoleMiddleware)
	bedsGroup.GET("/genomes", bedroutes.GenomeRoute)
	bedsGroup.GET("/platforms/:assembly", bedroutes.PlatformRoute)
	bedsGroup.GET("/search/:assembly", bedroutes.SearchBedsRoute)
	bedsGroup.POST("/regions", bedroutes.BedRegionsRoute)

	//
	// module groups: end
	//

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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// required so it can listen externally within docker container
	r.Run("0.0.0.0:8080") //"localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
