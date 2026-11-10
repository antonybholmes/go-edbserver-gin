package main

import (
	"context"
	"time"

	"net/http"
	"runtime"

	"github.com/antonybholmes/go-beds/beddb"
	"github.com/antonybholmes/go-cytobands/cytobanddb"
	"github.com/antonybholmes/go-dna/dnadb"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	adminroutes "github.com/antonybholmes/go-edbserver-gin/routes/admin"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	sessionroutes "github.com/antonybholmes/go-edbserver-gin/routes/session"

	"github.com/antonybholmes/go-edbserver-gin/routes/modules"
	"github.com/antonybholmes/go-hubs/hubdb"
	mailserver "github.com/antonybholmes/go-mailserver"
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
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/antonybholmes/go-pathway/pathwaydb"
	"github.com/antonybholmes/go-scrna/scrnadb"
	"github.com/antonybholmes/go-seqs/seqdb"
	"github.com/antonybholmes/go-sys/env"
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
)

var (
	store cookie.Store
	rdb   *redis.Client
	re    *access.RuleEngine
)

func initDatabases() {
	dnadb.InitDnaDB(consts.DnaDir)
	genomedb.InitCache(consts.GenomesDB)
	//microarraydb.InitDB("data/microarray")

	gexdb.InitGexDB(consts.GexDB)

	scrnadb.InitScrnaDB(consts.ScrnaDir)

	wgsdb.InitDB(consts.WGSDB)

	geneconvdb.InitGeneConvDB(consts.GeneConvDB)

	motifsdb.InitMotifDB(consts.MotifsDB)

	pathwaydb.InitPathwayDB(consts.PathwayDB)

	seqdb.InitSeqDB(consts.SeqsDB)

	cytobanddb.InitCytobandDB(consts.CytobandsDir)

	beddb.InitBedDB(consts.BedsDB)

	hubdb.InitHubDB(consts.HubsDir)
}

func init() {
	//initLogger()

	log.SetAppName(consts.AppName)
	env.Ls()

	initDatabases()

	rdb = redis.NewClient(&redis.Options{
		Addr:     consts.RedisAddr,
		Username: "edb",
		Password: consts.RedisPassword,
		DB:       0,
	})

	mailqueue.InitMailQueue(mailserver.NewSqsEmailQueue(consts.SqsQueueUrl))

	re = access.NewRuleEngine()
	re.LoadRules("config/access-rules.json")
}

func aboutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, AboutResp{
		Name:      consts.AppName,
		Version:   consts.Version.Version,
		Build:     consts.Version.Build,
		Updated:   consts.Version.Updated,
		Copyright: consts.Copyright,
	})
}

func infoHandler(c *gin.Context) {
	web.MakeDataResp(c, "", InfoResp{
		Arch:   runtime.GOARCH,
		IpAddr: c.ClientIP(),
	})
}

func setupMiddlewares(r *gin.Engine) {
	r.Use(otelgin.Middleware("edb-server"))
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     consts.AllowOrigins,
		AllowMethods:     consts.AllowMethods,
		AllowHeaders:     consts.AllowHeaders,
		AllowCredentials: consts.AllowCredentials,
		MaxAge:           consts.CorsMaxAge,
	}))

	store = cookie.NewStore([]byte(consts.SessionKey),
		[]byte(consts.SessionEncryptionKey))
	r.Use(sessions.Sessions(consts.SessionName, store))
}

func registerAppRoutes(r *gin.Engine,
	jwtUserMiddleWare gin.HandlerFunc,
	rulesMiddleware gin.HandlerFunc,
	updateTokenMiddleware gin.HandlerFunc,
	otp *auth.OTP) {
	r.GET("/about", aboutHandler)
	r.GET("/info", infoHandler)

	adminroutes.RegisterRoutes(r, rulesMiddleware)
	authenticationroutes.RegisterRoutes(r, jwtUserMiddleWare, updateTokenMiddleware)
	sessionroutes.RegisterRoutes(r, otp, jwtUserMiddleWare)
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
	tokengen.Init(token.NewES256TokenSigner(consts.JwtES256PrivateKey))

	claimsParser := middleware.NewUserJWTParser(middleware.NewJwtClaimsES256Parser(consts.JwtES256PublicKey))
	jwtUserMiddleWare := middleware.UserJWTMiddleware(claimsParser)
	rulesMiddleware := middleware.RulesMiddleware(claimsParser, re)
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

	r := gin.New()

	setupMiddlewares(r)

	registerAppRoutes(r, jwtUserMiddleWare, rulesMiddleware, updateTokenMiddleware, otp)

	log.Info().Msgf("starting server on %s", consts.AppPort)
	r.Run(consts.AppPort)
}
