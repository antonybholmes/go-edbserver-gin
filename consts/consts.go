package consts

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"os"
	"time"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/antonybholmes/go-sys/log"

	"github.com/golang-jwt/jwt/v5"
	"go.yaml.in/yaml/v4"
)

type AppConfig struct {
	App struct {
		URL     string `yaml:"url"`
		Domain  string `yaml:"domain"`
		LogFile string `yaml:"logFile"`
		Port    string `yaml:"port"`
	} `yaml:"app"`
	Cors struct {
		AllowOrigins     []string `yaml:"allowOrigins"`
		AllowMethods     []string `yaml:"allowMethods"`
		AllowHeaders     []string `yaml:"allowHeaders"`
		AllowCredentials bool     `yaml:"allowCredentials"`
		MaxAge           int      `yaml:"maxAge"`
	} `yaml:"cors"`
	Session struct {
		Name     string `yaml:"name"`
		TTLHours int    `yaml:"ttlHours"`
	} `yaml:"session"`
	Auth struct {
		PasswordlessTokenTtlMins int `yaml:"passwordlessTokenTtlMins"`
		AccessTokenTtlMins       int `yaml:"accessTokenTtlMins"`
		OtpTokenTtlMins          int `yaml:"otpTokenTtlMins"`
		ShortTtlMins             int `yaml:"shortTtlMins"`
	} `yaml:"auth"`
	Databases struct {
		Motifs       string `yaml:"motifs"`
		WGS          string `yaml:"wgs"`
		GeneConv     string `yaml:"geneconv"`
		Pathway      string `yaml:"pathway"`
		Gex          string `yaml:"gex"`
		Seqs         string `yaml:"seqs"`
		Beds         string `yaml:"beds"`
		Genomes      string `yaml:"genomes"`
		ScrnaDir     string `yaml:"scrnaDir"`
		CytobandsDir string `yaml:"cytobandsDir"`
		DnaDir       string `yaml:"dnaDir"`
		HubsDir      string `yaml:"hubsDir"`
	} `yaml:"databases"`
}

const (
	Name       = "Experiments Server"
	AppName    = "edbserver"
	Copyright  = "Copyright (C) 2024-2025 Antony Holmes"
	ConfigFile = "config/app.yaml"
)

var (
	AppUrl    string
	AppDomain string
	AppPort   string
	Version   sys.VersionInfo

	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	CorsMaxAge       time.Duration

	PasswordlessTokenTtlMins time.Duration
	AccessTokenTtlMins       time.Duration
	OtpTokenTtlMins          time.Duration
	ShortTtlMins             time.Duration

	//JwtRsaPrivateKey *rsa.PrivateKey //[]byte
	//JwtRsaPublicKey  *rsa.PublicKey  //[]byte

	JwtES256PrivateKey *ecdsa.PrivateKey
	JwtES256PublicKey  *ecdsa.PublicKey

	JwtAuth0RsaPublicKey *rsa.PublicKey
	JwtClerkRsaPublicKey *rsa.PublicKey
	SupabaseJwtSecretKey string
	SessionName          string
	SessionKey           string
	SessionEncryptionKey string
	//Updated              string

	RedisAddr     string
	RedisPassword string

	Auth0Audience   string
	Auth0Domain     string
	Auth0EmailClaim string
	Auth0NameClaim  string

	CognitoAudience string
	CognitoDomain   string

	ClerkDomain   string
	ClerkAudience string

	SupabaseDomain   string
	SupabaseAudience string

	UrlResetEmail    string
	UrlResetPassword string
	UrlVerifyEmail   string

	SqsQueueUrl string

	MotifsDB     string
	WGSDB        string
	GeneConvDB   string
	PathwayDB    string
	CytobandsDir string
	GexDB        string
	SeqsDB       string
	HubsDir      string
	BedsDB       string
	DnaDir       string
	GenomesDB    string
	ScrnaDir     string
)

func loadConfig() AppConfig {
	data, err := os.ReadFile(ConfigFile)

	if err != nil {
		log.Fatal().Msgf("failed to read config file: %v", err)
	}

	var cfg AppConfig

	if err := yaml.Load(data, &cfg); err != nil {
		log.Fatal().Msgf("failed to parse config file: %v", err)
	}

	return cfg
}

func init() {
	// val, exists := os.LookupEnv("APP_ENV")

	// if !exists || val != "production" {
	// 	env.Load("consts.env")
	// }

	env.Load("version.env")

	cfg := loadConfig()

	AppUrl = cfg.App.URL
	AppDomain = cfg.App.Domain
	AppPort = cfg.App.Port
	AllowOrigins = cfg.Cors.AllowOrigins
	AllowMethods = cfg.Cors.AllowMethods
	AllowHeaders = cfg.Cors.AllowHeaders
	AllowCredentials = cfg.Cors.AllowCredentials
	CorsMaxAge = time.Duration(cfg.Cors.MaxAge) * time.Second

	SessionName = cfg.Session.Name

	PasswordlessTokenTtlMins = time.Duration(cfg.Auth.PasswordlessTokenTtlMins) * time.Minute
	AccessTokenTtlMins = time.Duration(cfg.Auth.AccessTokenTtlMins) * time.Minute
	OtpTokenTtlMins = time.Duration(cfg.Auth.OtpTokenTtlMins) * time.Minute
	ShortTtlMins = time.Duration(cfg.Auth.ShortTtlMins) * time.Minute

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	UrlResetEmail = os.Getenv("URL_RESET_EMAIL")
	UrlResetPassword = os.Getenv("URL_RESET_PASSWORD")
	UrlVerifyEmail = os.Getenv("URL_VERIFY_EMAIL")

	SessionKey = os.Getenv("SESSION_KEY")
	SessionEncryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")

	Auth0Audience = os.Getenv("AUTH0_AUDIENCE")
	Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	Auth0EmailClaim = os.Getenv("AUTH0_EMAIL_CLAIM")
	Auth0NameClaim = os.Getenv("AUTH0_NAME_CLAIM")

	CognitoAudience = os.Getenv("COGNITO_CLIENT_ID")
	CognitoDomain = os.Getenv("COGNITO_DOMAIN")

	ClerkDomain = os.Getenv("CLERK_DOMAIN")
	ClerkAudience = os.Getenv("CLERK_AUDIENCE")

	SupabaseDomain = os.Getenv("SUPABASE_DOMAIN")
	SupabaseAudience = os.Getenv("SUPABASE_AUDIENCE")
	SupabaseJwtSecretKey = os.Getenv("SUPABASE_JWT_SECRET_KEY")

	SqsQueueUrl = os.Getenv("SQS_QUEUE_URL")

	MotifsDB = cfg.Databases.Motifs
	WGSDB = cfg.Databases.WGS
	GeneConvDB = cfg.Databases.GeneConv
	PathwayDB = cfg.Databases.Pathway
	CytobandsDir = cfg.Databases.CytobandsDir
	GexDB = cfg.Databases.Gex
	SeqsDB = cfg.Databases.Seqs
	HubsDir = cfg.Databases.HubsDir
	BedsDB = cfg.Databases.Beds
	DnaDir = cfg.Databases.DnaDir
	GenomesDB = cfg.Databases.Genomes
	ScrnaDir = cfg.Databases.ScrnaDir

	//
	// Load EC keys
	//

	bytes, err := os.ReadFile("jwt.es256.private.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtES256PrivateKey, err = jwt.ParseECPrivateKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("jwt.es256.public.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtES256PublicKey, err = jwt.ParseECPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	//
	// Keys for OAuth providers
	//
	bytes, err = os.ReadFile("auth0.key.pub")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtAuth0RsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("clerk.key.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtClerkRsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	Version = sys.Must(sys.LoadVersionInfo("version.json"))
}
