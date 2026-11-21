package main

import (
	"context"
	"time"

	"net/http"
	"runtime"

	"github.com/antonybholmes/go-beds/beddb"
	"github.com/antonybholmes/go-cytobands/cytobanddb"
	"github.com/antonybholmes/go-dna/dnadb"
	"github.com/antonybholmes/go-edbserver-gin/config"
	adminroutes "github.com/antonybholmes/go-edbserver-gin/routes/admin"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	sessionroutes "github.com/antonybholmes/go-edbserver-gin/routes/session"
	userdbcache "github.com/antonybholmes/go-web/auth/userdb/cache"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/antonybholmes/go-edbserver-gin/routes/modules"
	"github.com/antonybholmes/go-hubs/hubdb"
	awsmail "github.com/antonybholmes/go-mail/aws"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/access"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/token"
	"github.com/antonybholmes/go-web/auth/token/tokengen"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/antonybholmes/go-web/middleware"

	utilsroutes "github.com/antonybholmes/go-edbserver-gin/routes/utils"
	"github.com/antonybholmes/go-geneconv/geneconvdb"
	"github.com/antonybholmes/go-genome/genomedb"
	"github.com/antonybholmes/go-gex/gexdb"
	"github.com/antonybholmes/go-mail/mailqueue"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/antonybholmes/go-pathway/pathwaydb"
	"github.com/antonybholmes/go-scrna/scrnadb"
	"github.com/antonybholmes/go-seqs/seqdb"
	"github.com/antonybholmes/go-wgs/wgsdb"
	_ "github.com/mattn/go-sqlite3"
)

type (
	AboutResp struct {
		Name      string `json:"name"`
		Copyright string `json:"copyright"`
		Version   string `json:"version"`
		Updated   string `json:"updated"`
		Build     int    `json:"build"`
	}

	InfoResp struct {
		IpAddr string `json:"ipAddr"`
		Arch   string `json:"arch"`
	}

	AboutHandler struct {
		Config *config.AppConfig
	}
)

var (
	store cookie.Store
	rdb   *redis.Client
)

func initDatabases(cfg *config.AppConfig) {
	dnadb.InitDnaDB(cfg.Databases.DnaDir)
	genomedb.InitCache(cfg.Databases.Genomes)
	//microarraydb.InitDB("data/microarray")

	gexdb.InitGexDB(cfg.Databases.Gex)

	scrnadb.InitScrnaDB(cfg.Databases.ScrnaDir)

	wgsdb.InitDB(cfg.Databases.WGS)

	geneconvdb.InitGeneConvDB(cfg.Databases.GeneConv)

	motifsdb.InitMotifDB(cfg.Databases.Motifs)

	pathwaydb.InitPathwayDB(cfg.Databases.Pathway)

	seqdb.InitSeqDB(cfg.Databases.Seqs)

	cytobanddb.InitCytobandDB(cfg.Databases.CytobandsDir)

	beddb.InitBedDB(cfg.Databases.Beds)

	hubdb.InitHubDB(cfg.Databases.HubsDir)
}

func (h *AboutHandler) aboutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, AboutResp{
		Name:      config.AppName,
		Version:   h.Config.App.Version.Version,
		Build:     h.Config.App.Version.Build,
		Updated:   h.Config.App.Version.Updated,
		Copyright: config.Copyright,
	})
}

func (h *AboutHandler) infoHandler(c *gin.Context) {
	web.MakeDataResp(c, "", InfoResp{
		Arch:   runtime.GOARCH,
		IpAddr: c.ClientIP(),
	})
}

func setupMiddlewares(r *gin.Engine, cfg *config.AppConfig) {
	// r.Use(otelgin.Middleware("edb-server"))
	// r.Use(gin.Recovery())
	// r.Use(middleware.LoggingMiddleware())
	// r.Use(middleware.ErrorHandlerMiddleware())

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     consts.AllowOrigins,
	// 	AllowMethods:     consts.AllowMethods,
	// 	AllowHeaders:     consts.AllowHeaders,
	// 	AllowCredentials: consts.AllowCredentials,
	// 	MaxAge:           consts.CorsMaxAge,
	// }))

	// r.Use(sessions.Sessions(consts.SessionName, store))

	corsConfig := cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowHeaders:     cfg.Cors.AllowHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		MaxAge:           time.Duration(cfg.Cors.MaxAge) * time.Second,
	}

	r.Use(cors.New(corsConfig))                       // Check CORS first
	r.Use(middleware.ErrorHandlerMiddleware())        // Handle errors early
	r.Use(middleware.LoggingMiddleware())             // Log relevant requests only
	r.Use(otelgin.Middleware("edb-server"))           // Trace valid requests
	r.Use(sessions.Sessions(cfg.Session.Name, store)) // Load sessions for routes that need them
	r.Use(gin.Recovery())                             // Recover from panics
}

func registerAppRoutes(r *gin.Engine,
	jwtUserMiddleWare gin.HandlerFunc,
	rulesMiddleware gin.HandlerFunc,
	updateTokenMiddleware gin.HandlerFunc,
	otp *auth.OTP,
	cfg *config.AppConfig) {

	aboutHandler := &AboutHandler{
		Config: cfg,
	}

	r.GET("/about", aboutHandler.aboutHandler)
	r.GET("/info", aboutHandler.infoHandler)

	adminroutes.RegisterRoutes(r, rulesMiddleware, userdbcache.Instance(), cfg)
	authenticationroutes.RegisterRoutes(r, jwtUserMiddleWare, updateTokenMiddleware, cfg)
	sessionroutes.RegisterRoutes(r, otp, jwtUserMiddleWare, cfg)
	modules.RegisterRoutes(r, rulesMiddleware)

	utilsGroup := r.Group("/utils")
	xlsxGroup := utilsGroup.Group("/xlsx")
	xlsxGroup.POST("/sheets", utilsroutes.XlsxSheetsRoute)
	xlsxGroup.POST("/to/:format", utilsroutes.XlsxToRoute)

	utilsGroup.GET("/passwords/hash", utilsroutes.HashedPasswordRoute)
	utilsGroup.GET("/randkey", utilsroutes.RandomKeyRoute)
	utilsGroup.GET("/uuidv7", utilsroutes.UUIDv7Route)
}

func main() {
	cfg := config.Instance()

	log.SetAppName(config.AppName)
	//env.Ls()

	store = cookie.NewStore([]byte(cfg.Session.Key),
		[]byte(cfg.Session.EncryptionKey))

	initDatabases(cfg)

	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: "edb",
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msgf("load AWS configuration: %v", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	mailqueue.InitMailQueue(awsmail.NewSQSEmailQueue(sqsClient, cfg.AWS.SQS.QueueUrl))

	re := access.NewRuleEngine()
	re.LoadRules("config/access-rules.json")

	tokengen.Init(token.NewES256TokenSigner(cfg.Auth.Keys.JwtES256PrivateKey))

	claimsParser := middleware.NewUserJWTParser(middleware.NewJwtClaimsES256Parser(cfg.Auth.Keys.JwtES256PublicKey))
	jwtUserMiddleWare := middleware.UserJWTMiddleware(claimsParser)
	rulesMiddleware := middleware.AccessTokenRulesMiddleware(claimsParser, re)

	updateTokenMiddleware := middleware.JwtIsUpdateTokenMiddleware()

	otp := auth.NewDefaultOTP(rdb)

	defer rdb.Close()

	tp, err := initTracerProvider()
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer provider: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal().Msgf("failed to shutdown tracer provider: %v", err)
		}
	}()

	router := gin.New()
	router.SetTrustedProxies(cfg.App.TrustedProxies)

	setupMiddlewares(router, cfg)

	registerAppRoutes(router,
		jwtUserMiddleWare,
		rulesMiddleware,
		updateTokenMiddleware,
		otp,
		cfg)

	log.Info().Msgf("starting server on %s", cfg.App.Port)
	router.Run(cfg.App.Port)
}
