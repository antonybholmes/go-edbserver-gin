package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	adminroutes "github.com/antonybholmes/go-edb-server-gin/routes/admin"
	auth0routes "github.com/antonybholmes/go-edb-server-gin/routes/auth0"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	authorizationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authorization"
	bedroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/beds"
	cytobandroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/cytobands"
	dnaroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/dna"
	geneconvroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/geneconv"
	generoutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/genes"
	gexroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/gex"
	motifroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/motifs"
	mutationroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/mutation"
	pathwayroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/pathway"
	seqroutes "github.com/antonybholmes/go-edb-server-gin/routes/modules/seqs"

	utilsroutes "github.com/antonybholmes/go-edb-server-gin/routes/utils"
	"github.com/antonybholmes/go-geneconv/geneconvdbcache"
	"github.com/antonybholmes/go-genes/genedbcache"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
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

func init() {

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
	genedbcache.InitCache("data/modules/genes")
	//microarraydb.InitDB("data/microarray")

	gexdbcache.InitCache("data/modules/gex")

	mutationdbcache.InitCache("data/modules/mutations")

	geneconvdbcache.InitCache("data/modules/geneconv/geneconv.db")

	motifsdb.InitCache("data/modules/motifs/motifs.db")

	pathwaydbcache.InitCache("data/modules/pathway/pathway-v2.db")

	seqsdbcache.InitCache("data/modules/seqs/")

	cytobandsdbcache.InitCache("data/modules/cytobands/")

	bedsdbcache.InitCache("data/modules/beds/")

	rdb := redis.NewClient(&redis.Options{
		Addr:     consts.REDIS_ADDR,
		Username: "edb",
		Password: consts.REDIS_PASSWORD,
		DB:       0, // use default DB
	})

	queue.Init(mailer.NewRedisEmailPublisher(rdb))

	// writer := kafka.NewWriter(kafka.WriterConfig{
	// 	Brokers:  []string{"localhost:9094"}, // Kafka broker
	// 	Topic:    mailer.QUEUE_EMAIL_CHANNEL, // Topic name
	// 	Balancer: &kafka.LeastBytes{},        // Balancer (optional)
	// })

	// queue.Init(mailer.NewKafkaEmailPublisher(writer))
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
	//email := gomailer.QueueEmail{To: "antony@antonybholmes.dev"}
	//rdb.PublishEmail(&email)

	//
	// Set logging to file
	//

	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", consts.APP_NAME),
		MaxSize:    10,   // Max size in MB before rotating
		MaxBackups: 3,    // Keep 3 backup files
		MaxAge:     7,    // Retain files for 7 days
		Compress:   true, // Compress old log files
	}

	logger := zerolog.New(io.MultiWriter(os.Stderr, fileLogger)).With().Timestamp().Logger()

	// we use != development because it means we need to set the env variable in order
	// to see debugging work. The default is to assume production, in which case we use
	// lumberjack
	if os.Getenv("APP_ENV") != "development" {
		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
	}

	sessionRoutes := authenticationroutes.NewSessionRoutes()

	//r := gin.Default()
	r := gin.New()

	r.Use(LoggingMiddleware(logger))
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8000",
			"https://edb.rdf-lab.org",
			"https://dev.edb-app-astro.pages.dev",
			"https://edb-client-astro.pages.dev",
			"https://edb-client-next.pages.dev"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Set-Cookie"},
		// for sharing session cookie for validating logins etc
		AllowCredentials: true,      // Allow credentials (cookies, HTTP authentication)
		MaxAge:           12 * 3600, // Cache preflight response for 12 hours
	}))

	r.Use(ErrorHandlerMiddleware())

	store = cookie.NewStore([]byte(consts.SESSION_KEY), []byte(consts.SESSION_ENCRYPTION_KEY))
	r.Use(sessions.Sessions(consts.SESSION_NAME, store))

	r.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			AboutResp{Name: consts.APP_NAME,
				Version:   consts.VERSION,
				Updated:   consts.UPDATED,
				Copyright: consts.COPYRIGHT})
	})

	r.GET("/info", func(c *gin.Context) {
		routes.MakeDataResp(c, "", InfoResp{Arch: runtime.GOARCH, IpAddr: c.ClientIP()})
	})

	//
	// Routes
	//

	adminGroup := r.Group("/admin",
		JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware(),
		JwtIsAdminMiddleware())

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
	auth0Group := authGroup.Group("/auth0")
	auth0Group.POST("/validate", JwtAuth0Middleware(), auth0routes.ValidateAuth0TokenRoute)

	authGroup.POST("/signin", authenticationroutes.UsernamePasswordSignInRoute)

	emailGroup := authGroup.Group("/email")

	emailGroup.POST("/verified",
		JwtParseMiddleware(),
		JwtIsVerifyEmailTokenMiddleware(),
		authenticationroutes.EmailAddressVerifiedRoute,
	)

	// with the correct token, performs the update
	emailGroup.POST("/reset", JwtParseMiddleware(), authenticationroutes.SendResetEmailEmailRoute)
	// with the correct token, performs the update
	emailGroup.POST("/update", JwtParseMiddleware(), authenticationroutes.UpdateEmailRoute)

	passwordGroup := authGroup.Group("/passwords")

	// sends a reset link
	passwordGroup.POST("/reset", authenticationroutes.SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update", JwtParseMiddleware(), authenticationroutes.UpdatePasswordRoute)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c *gin.Context) {
		authenticationroutes.PasswordlessSigninEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin",
		JwtParseMiddleware(),
		authenticationroutes.PasswordlessSignInRoute,
	)

	tokenGroup := authGroup.Group("/tokens", JwtParseMiddleware())
	tokenGroup.POST("/info", authorizationroutes.TokenInfoRoute)
	tokenGroup.POST("/access", authorizationroutes.NewAccessTokenRoute)

	usersGroup := authGroup.Group("/users", JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware())

	usersGroup.POST("", authorizationroutes.UserRoute)

	usersGroup.POST("/update", authorizationroutes.UpdateUserRoute)

	//usersGroup.POST("/passwords/update", authentication.UpdatePasswordRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := r.Group("/sessions")

	//sessionAuthGroup := sessionGroup.Group("/auth")

	sessionGroup.POST("/auth0/signin", JwtAuth0Middleware(), sessionRoutes.SessionSignInUsingAuth0Route)

	sessionGroup.POST("/auth/signin", sessionRoutes.SessionUsernamePasswordSignInRoute)
	sessionGroup.POST("/auth/passwordless/validate",
		JwtParseMiddleware(),
		sessionRoutes.SessionPasswordlessValidateSignInRoute)

	sessionGroup.POST("/api/keys/signin", sessionRoutes.SessionApiKeySignInRoute)

	//sessionGroup.POST("/init", sessionRoutes.InitSessionRoute)
	sessionGroup.GET("/info", sessionRoutes.SessionInfoRoute)

	sessionGroup.POST("/signout", authenticationroutes.SessionSignOutRoute)

	//sessionGroup.POST("/email/reset", authentication.SessionSendResetEmailEmailRoute)

	//sessionGroup.POST("/password/reset", authentication.SessionSendResetPasswordEmailRoute)

	sessionGroup.POST("/tokens/access",
		SessionIsValidMiddleware(),
		authenticationroutes.NewAccessTokenFromSessionRoute)

	sessionGroup.POST("/refresh",
		SessionIsValidMiddleware(),
		sessionRoutes.SessionRenewRoute)

	sessionUserGroup := sessionGroup.Group("/user", SessionIsValidMiddleware())
	sessionUserGroup.GET("", authenticationroutes.UserFromSessionRoute)
	sessionUserGroup.POST("/update", authorizationroutes.SessionUpdateUserRoute)

	// sessionPasswordGroup := sessionAuthGroup.Group("/passwords")
	// sessionPasswordGroup.Use(SessionIsValidMiddleware)

	// sessionPasswordGroup.POST("/update", func(c *gin.Context) {
	// 	return authentication.SessionUpdatePasswordRoute(c)
	// })

	//
	// sessions: end
	//

	//
	// passwordless groups: end
	//

	//
	// module groups: start
	//

	moduleGroup := r.Group("/modules")
	//moduleGroup.Use(jwtMiddleWare,JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")

	dnaGroup.POST("/:assembly", dnaroutes.DNARoute)

	dnaGroup.GET("/genomes", dnaroutes.GenomesRoute)

	genesGroup := moduleGroup.Group("/genes")

	genesGroup.GET("/genomes", generoutes.GenomesRoute)
	genesGroup.POST("/within/:assembly", generoutes.WithinGenesRoute)
	genesGroup.POST("/closest/:assembly", generoutes.ClosestGeneRoute)
	genesGroup.POST("/annotate/:assembly", generoutes.AnnotateRoute)
	genesGroup.POST("/overlap/:assembly", generoutes.OverlappingGenesRoute)
	genesGroup.GET("/info/:assembly", generoutes.GeneInfoRoute)
	// get version info about the database itself
	//genesGroup.GET("/db/:assembly", generoutes.GeneDBInfoRoute)

	// mutationsGroup := moduleGroup.Group("/mutations",
	// 	jwtMiddleWare,
	// 	JwtIsAccessTokenMiddleware,
	// 	NewJwtPermissionsMiddleware("rdf"))

	mutationsGroup := moduleGroup.Group("/mutations")
	mutationsGroup.GET("/datasets/:assembly", mutationroutes.MutationDatasetsRoute)
	mutationsGroup.POST("/:assembly/:name", mutationroutes.MutationsRoute)
	mutationsGroup.POST("/maf/:assembly", mutationroutes.PileupRoute)

	mutationsGroup.POST("/pileup/:assembly",
		JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware(),
		RDFMiddleware(),
		mutationroutes.PileupRoute,
	)

	gexGroup := moduleGroup.Group("/gex")
	gexGroup.GET("/platforms", gexroutes.PlaformsRoute)
	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)
	gexGroup.POST("/datasets", gexroutes.GexDatasetsRoute)
	gexGroup.POST("/exp",
		JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware(),
		RDFMiddleware(),
		gexroutes.GexGeneExpRoute,
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
		JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware(),
		RDFMiddleware())

	seqsGroup.GET("/genomes", seqroutes.GenomeRoute)
	seqsGroup.GET("/platforms/:assembly", seqroutes.PlatformRoute)
	//tracksGroup.GET("/:platform/:assembly/tracks", seqroutes.TracksRoute)
	seqsGroup.GET("/search/:assembly", seqroutes.SearchSeqRoute)
	seqsGroup.POST("/bins", seqroutes.BinsRoute)

	cytobandsGroup := moduleGroup.Group("/cytobands")
	cytobandsGroup.GET("/:assembly/:chr", cytobandroutes.CytobandsRoute)

	bedsGroup := moduleGroup.Group("/beds", JwtParseMiddleware(),
		JwtIsAccessTokenMiddleware(),
		RDFMiddleware())
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
