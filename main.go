package main

import (
	"context"
	"time"

	"net/http"
	"runtime"

	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	adminroutes "github.com/antonybholmes/go-edbserver-gin/routes/admin"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	sessionroutes "github.com/antonybholmes/go-edbserver-gin/routes/session"

	"github.com/antonybholmes/go-edbserver-gin/routes/modules"
	"github.com/antonybholmes/go-hubs/hubsdbcache"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/access"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/antonybholmes/go-web/middleware"

	utilsroutes "github.com/antonybholmes/go-edbserver-gin/routes/utils"
	"github.com/antonybholmes/go-geneconv/geneconvdbcache"
	"github.com/antonybholmes/go-genome/genomedbcache"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-pathway/pathwaydbcache"
	scrnadbcache "github.com/antonybholmes/go-scrna/cache"
	"github.com/antonybholmes/go-seqs/seqsdbcache"
	"github.com/antonybholmes/go-sys/env"
	_ "github.com/mattn/go-sqlite3"
)

type (
	AboutResp struct {
		Name      string `json:"name"`
		Copyright string `json:"copyright"`
		Version   string `json:"version"`
		Updated   string `json:"updated"`
	}

	InfoResp struct {
		IpAddr string `json:"ipAddr"`
		Arch   string `json:"arch"`
	}
)

const PreflightMaxAge = 12 * 3600 // 12 hours

// var store *sqlitestorr.SqliteStore
var (
	store cookie.Store

	rdb *redis.Client

	re *access.RuleEngine
)

// func initLogger() {

// 	//multiWriter := io.MultiWriter(os.Stderr, fileLogger)

// 	//logger := zerolog.New(multiWriter).With().Timestamp().Logger()
// 	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

// 	// we use != development because it means we need to set the env variable in order
// 	// to see debugging work. The default is to assume production, in which case we use
// 	// lumberjack
// 	if os.Getenv("APP_ENV") != "development" {
// 		// reduce log events in production
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)

// 		// log events to log files with rotation
// 		fileLogger := &lumberjack.Logger{
// 			Filename:   fmt.Sprintf("logs/%s.log", consts.AppName),
// 			MaxSize:    10,   // Max size in MB before rotating
// 			MaxBackups: 3,    // Keep 3 backup files
// 			MaxAge:     7,    // Retain files for 7 days
// 			Compress:   true, // Compress old log files
// 		}

// 		logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, fileLogger)).With().Timestamp().Logger()
// 	}

// 	log.Logger = logger
// }

func init() {
	//initLogger()

	log.SetAppName(consts.AppName)

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

	//userdbcache.InitCache() //"data/users.db")

	//mailserver.Init()

	dnadbcache.InitCache("data/modules/dna")
	genomedbcache.InitCache("data/modules/genome")
	//microarraydb.InitDB("data/microarray")

	gexdbcache.InitCache("data/modules/gex")

	scrnadbcache.InitCache("data/modules/scrna")

	mutationdbcache.InitCache("data/modules/mutations")

	geneconvdbcache.InitCache("data/modules/geneconv/geneconv.db")

	motifsdb.InitCache("data/modules/motifs/motifs.db")

	pathwaydbcache.InitCache("data/modules/pathway/pathway-v3.db")

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

	mailqueue.InitMailQueue(mailserver.NewSqsEmailQueue(consts.SqsQueueUrl))

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

	//accessTokenMiddleware := middleware.JwtIsAccessTokenMiddleware()

	rulesMiddleware := middleware.RulesMiddleware(claimsParser, re)

	updateTokenMiddleware := middleware.JwtIsUpdateTokenMiddleware()

	//rdfRoleMiddleware := middleware.JwtHasRDFRoleMiddleware()

	otp := auth.NewDefaultOTP(rdb)

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
		AllowCredentials: true,            // Allow credentials (cookies, HTTP authentication)
		MaxAge:           PreflightMaxAge, // Cache preflight response for 12 hours
	}))

	store = cookie.NewStore([]byte(consts.SessionKey),
		[]byte(consts.SessionEncryptionKey))
	r.Use(sessions.Sessions(consts.SessionName, store))

	r.GET("/about", func(c *gin.Context) {

		c.JSON(http.StatusOK,
			AboutResp{
				Name:      consts.AppName,
				Version:   consts.Version.Version,
				Updated:   consts.Version.Updated,
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

	adminroutes.RegisterRoutes(r, rulesMiddleware)

	authenticationroutes.RegisterRoutes(r, jwtUserMiddleWare, updateTokenMiddleware)

	sessionroutes.RegisterRoutes(r,
		otp,
		jwtUserMiddleWare)

	//
	// Deal with logins where we want a session
	//

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
