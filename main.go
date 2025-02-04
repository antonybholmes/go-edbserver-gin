package main

import (
	"net/http"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

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
	"github.com/antonybholmes/go-edb-server-gin/routes/authorization"
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
	toolsroutes "github.com/antonybholmes/go-edb-server-gin/routes/tools"
	utilroutes "github.com/antonybholmes/go-edb-server-gin/routes/util"
	"github.com/antonybholmes/go-geneconv/geneconvdbcache"
	"github.com/antonybholmes/go-genes/genedbcache"
	"github.com/antonybholmes/go-gex/gexdbcache"
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
	//email := gomailer.RedisQueueEmail{To: "antony@antonybholmes.dev"}
	//rdb.PublishEmail(&email)

	//
	// Set logging to file
	//

	// fileLogger := &lumberjack.Logger{
	// 	Filename:   fmt.Sprintf("logs/%s.log", consts.APP_NAME),
	// 	MaxSize:    5, //
	// 	MaxBackups: 10,
	// 	MaxAge:     14,
	// 	Compress:   true,
	// }

	//logger := zerolog.New(io.MultiWriter(os.Stderr, fileLogger)).With().Timestamp().Logger()

	// we use != development because it means we need to set the env variable in order
	// to see debugging work. The default is to assume production, in which case we use
	// lumberjack
	// if os.Getenv("APP_ENV") != "development" {
	// 	logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
	// }

	sessionRoutes := authenticationroutes.NewSessionRoutes()
	rdfMiddlware := RDFMiddleware()

	r := gin.Default()

	//r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8000",
			"https://edb.rdf-lab.org",
			"https://dev.edb-app-astro.pages.dev",
			"https://edb-client-astro.pages.dev"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Set-Cookie"},
		// for sharing session cookie for validating logins etc
		AllowCredentials: true,      // Allow credentials (cookies, HTTP authentication)
		MaxAge:           12 * 3600, // Cache preflight response for 12 hours
	}))

	r.Use(ErrorHandler())

	store = cookie.NewStore([]byte(consts.SESSION_KEY), []byte(consts.SESSION_ENCRYPTION_KEY))
	r.Use(sessions.Sessions(consts.SESSION_NAME, store))

	r.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			AboutResp{Name: consts.NAME,
				Version:   consts.VERSION,
				Updated:   consts.UPDATED,
				Copyright: consts.COPYRIGHT})
	})

	r.GET("/info", func(c *gin.Context) {
		routes.MakeDataResp(c, "", InfoResp{Arch: runtime.GOARCH, IpAddr: c.ClientIP()})
	})

	toolsGroup := r.Group("/tools")
	toolsGroup.GET("/passwords/hash", toolsroutes.HashedPasswordRoute)
	toolsGroup.GET("/key", toolsroutes.RandomKeyRoute)

	//
	// Routes
	//

	adminGroup := r.Group("/admin",
		JwtMiddleware(),
		JwtIsAccessTokenMiddleware(),
		JwtHasAdminPermissionMiddleware())

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

	emailGroup.POST("/verify",
		authenticationroutes.EmailAddressVerifiedRoute,
		JwtMiddleware())

	// with the correct token, performs the update
	emailGroup.POST("/reset", JwtMiddleware(), authenticationroutes.SendResetEmailEmailRoute)
	// with the correct token, performs the update
	emailGroup.POST("/update", JwtMiddleware(), authenticationroutes.UpdateEmailRoute)

	passwordGroup := authGroup.Group("/passwords")

	// sends a reset link
	passwordGroup.POST("/reset", authenticationroutes.SendResetPasswordFromUsernameEmailRoute)

	// with the correct token, updates a password
	passwordGroup.POST("/update", JwtMiddleware(), authenticationroutes.UpdatePasswordRoute)

	passwordlessGroup := authGroup.Group("/passwordless")

	passwordlessGroup.POST("/email", func(c *gin.Context) {
		authenticationroutes.PasswordlessSigninEmailRoute(c, nil)
	})

	passwordlessGroup.POST("/signin",
		authenticationroutes.PasswordlessSignInRoute,
		JwtMiddleware())

	tokenGroup := authGroup.Group("/tokens", JwtMiddleware())
	tokenGroup.POST("/info", authorization.TokenInfoRoute)
	tokenGroup.POST("/access", authorization.NewAccessTokenRoute)

	usersGroup := authGroup.Group("/users", JwtMiddleware(),
		JwtIsAccessTokenMiddleware())

	usersGroup.POST("", authorization.UserRoute)

	usersGroup.POST("/update", authorization.UpdateUserRoute)

	//usersGroup.POST("/passwords/update", authentication.UpdatePasswordRoute)

	//
	// Deal with logins where we want a session
	//

	sessionGroup := r.Group("/sessions")

	//sessionAuthGroup := sessionGroup.Group("/auth")

	sessionGroup.POST("/signin", sessionRoutes.SessionUsernamePasswordSignInRoute)
	sessionGroup.POST("/auth0/signin", JwtAuth0Middleware(), sessionRoutes.SessionSignInUsingAuth0Route)

	sessionGroup.POST("/passwordless/signin",
		JwtMiddleware(),
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

	sessionUserGroup.POST("/update", authorization.SessionUpdateUserRoute)

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
		JwtMiddleware(),
		JwtIsAccessTokenMiddleware(),
		rdfMiddlware,
		mutationroutes.PileupRoute,
	)

	gexGroup := moduleGroup.Group("/gex")
	gexGroup.GET("/platforms", gexroutes.PlaformsRoute)
	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)
	gexGroup.POST("/datasets", gexroutes.GexDatasetsRoute)
	gexGroup.POST("/exp",
		JwtMiddleware(),
		JwtIsAccessTokenMiddleware(),
		rdfMiddlware,
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
		JwtMiddleware(),
		JwtIsAccessTokenMiddleware(),
		rdfMiddlware)

	seqsGroup.GET("/genomes", seqroutes.GenomeRoute)
	seqsGroup.GET("/platforms/:assembly", seqroutes.PlatformRoute)
	//tracksGroup.GET("/:platform/:assembly/tracks", seqroutes.TracksRoute)
	seqsGroup.GET("/search/:assembly", seqroutes.SearchSeqRoute)
	seqsGroup.POST("/bins", seqroutes.BinsRoute)

	cytobandsGroup := moduleGroup.Group("/cytobands")
	cytobandsGroup.GET("/:assembly/:chr", cytobandroutes.CytobandsRoute)

	bedsGroup := moduleGroup.Group("/beds", JwtMiddleware(),
		JwtIsAccessTokenMiddleware(),
		rdfMiddlware)
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
	xlsxGroup.POST("/sheets", utilroutes.XlsxSheetsRoute)
	xlsxGroup.POST("/to/:format", utilroutes.XlsxToRoute)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
