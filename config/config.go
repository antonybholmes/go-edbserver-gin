package config

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"os"
	"sync"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/antonybholmes/go-sys/log"

	"github.com/golang-jwt/jwt/v5"
	"go.yaml.in/yaml/v4"
)

type AppConfig struct {
	App struct {
		URL            string          `yaml:"url"`
		Domain         string          `yaml:"domain"`
		LogFile        string          `yaml:"logFile"`
		Port           string          `yaml:"port"`
		TrustedProxies []string        `yaml:"trustedProxies"`
		Version        sys.VersionInfo `yaml:"-"`
	} `yaml:"app"`
	Cors struct {
		AllowOrigins     []string `yaml:"allowOrigins"`
		AllowMethods     []string `yaml:"allowMethods"`
		AllowHeaders     []string `yaml:"allowHeaders"`
		AllowCredentials bool     `yaml:"allowCredentials"`
		MaxAge           int      `yaml:"maxAge"`
	} `yaml:"cors"`
	Session struct {
		Name          string `yaml:"name"`
		TTLHours      int    `yaml:"ttlHours"`
		Key           string `yaml:"-"`
		EncryptionKey string `yaml:"-"`
	} `yaml:"session"`
	Auth struct {
		PasswordlessTokenTtlMins int `yaml:"passwordlessTokenTtlMins"`
		AccessTokenTtlMins       int `yaml:"accessTokenTtlMins"`
		OtpTokenTtlMins          int `yaml:"otpTokenTtlMins"`
		ShortTtlMins             int `yaml:"shortTtlMins"`
		Oauth2                   struct {
			Auth0 struct {
				Audience        string         `yaml:"-"`
				Domain          string         `yaml:"-"`
				EmailClaim      string         `yaml:"-"`
				NameClaim       string         `yaml:"-"`
				JwtRsaPublicKey *rsa.PublicKey `yaml:"-"`
			} `yaml:"-"`
			Clerk struct {
				Audience        string         `yaml:"-"`
				Domain          string         `yaml:"-"`
				JwtRsaPublicKey *rsa.PublicKey `yaml:"-"`
			} `yaml:"-"`
			Supabase struct {
				Audience     string `yaml:"-"`
				Domain       string `yaml:"-"`
				JwtSecretKey string `yaml:"-"`
			} `yaml:"-"`
			Cognito struct {
				Audience string `yaml:"-"`
				Domain   string `yaml:"-"`
			} `yaml:"-"`
		} `yaml:"-"`
		Keys struct {
			JwtES256PrivateKey *ecdsa.PrivateKey `yaml:"-"`
			JwtES256PublicKey  *ecdsa.PublicKey  `yaml:"-"`
		} `yaml:"-"`
	} `yaml:"auth"`
	Redis struct {
		Addr     string `yaml:"-"`
		Password string `yaml:"-"`
	} `yaml:"-"`
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
	Aws struct {
		Sqs struct {
			QueueUrl string `yaml:"-"`
		} `yaml:"-"`
	} `yaml:"-"`
	Account struct {
		ResetEmailUrl    string `yaml:"-"`
		ResetPasswordUrl string `yaml:"-"`
		VerifyEmailUrl   string `yaml:"-"`
	} `yaml:"-"`
}

const (
	Name       = "Experiments Server"
	AppName    = "edbserver"
	Copyright  = "Copyright (C) 2024-2025 Antony Holmes"
	ConfigFile = "config/app.yaml"
)

var (
	once sync.Once
	cfg  *AppConfig
	// AppUrl         string
	// AppDomain      string
	// AppPort        string
	// Version        sys.VersionInfo
	// TrustedProxies []string

	// AllowOrigins     []string
	// AllowMethods     []string
	// AllowHeaders     []string
	// AllowCredentials bool
	// CorsMaxAge       time.Duration

	// PasswordlessTokenTtlMins time.Duration
	// AccessTokenTtlMins       time.Duration
	// OtpTokenTtlMins          time.Duration
	// ShortTtlMins             time.Duration

	//JwtRsaPrivateKey *rsa.PrivateKey //[]byte
	//JwtRsaPublicKey  *rsa.PublicKey  //[]byte

	//JwtES256PrivateKey *ecdsa.PrivateKey
	//JwtES256PublicKey  *ecdsa.PublicKey

	//JwtClerkRsaPublicKey *rsa.PublicKey
	//SupabaseJwtSecretKey string
	//SessionName          string
	//SessionKey           string
	//SessionEncryptionKey string
	//Updated              string

	//UrlResetEmail    string
	//UrlResetPassword string
	//UrlVerifyEmail   string

	//SqsQueueUrl string

	// MotifsDB     string
	// WGSDB        string
	// GeneConvDB   string
	// PathwayDB    string
	// CytobandsDir string
	// GexDB        string
	// SeqsDB       string
	// HubsDir      string
	// BedsDB       string
	// DnaDir       string
	// GenomesDB    string
	// ScrnaDir     string
)

func loadYaml() AppConfig {
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

func load() *AppConfig {
	// val, exists := os.LookupEnv("APP_ENV")

	// if !exists || val != "production" {
	// 	env.Load("consts.env")
	// }

	env.Load("version.env")

	cfg := loadYaml()

	// AppUrl = cfg.App.URL
	// AppDomain = cfg.App.Domain
	// AppPort = cfg.App.Port
	// TrustedProxies = cfg.App.TrustedProxies
	// AllowOrigins = cfg.Cors.AllowOrigins
	// AllowMethods = cfg.Cors.AllowMethods
	// AllowHeaders = cfg.Cors.AllowHeaders
	// AllowCredentials = cfg.Cors.AllowCredentials
	// CorsMaxAge = time.Duration(cfg.Cors.MaxAge) * time.Second

	//SessionName = cfg.Session.Name

	// PasswordlessTokenTtlMins = time.Duration(cfg.Auth.PasswordlessTokenTtlMins) * time.Minute
	// AccessTokenTtlMins = time.Duration(cfg.Auth.AccessTokenTtlMins) * time.Minute
	// OtpTokenTtlMins = time.Duration(cfg.Auth.OtpTokenTtlMins) * time.Minute
	// ShortTtlMins = time.Duration(cfg.Auth.ShortTtlMins) * time.Minute

	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	cfg.Account.ResetEmailUrl = os.Getenv("URL_RESET_EMAIL")
	cfg.Account.ResetPasswordUrl = os.Getenv("URL_RESET_PASSWORD")
	cfg.Account.VerifyEmailUrl = os.Getenv("URL_VERIFY_EMAIL")

	cfg.Session.Key = os.Getenv("SESSION_KEY")
	cfg.Session.EncryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")

	cfg.Auth.Oauth2.Auth0.Audience = os.Getenv("AUTH0_AUDIENCE")
	cfg.Auth.Oauth2.Auth0.Domain = os.Getenv("AUTH0_DOMAIN")
	cfg.Auth.Oauth2.Auth0.EmailClaim = os.Getenv("AUTH0_EMAIL_CLAIM")
	cfg.Auth.Oauth2.Auth0.NameClaim = os.Getenv("AUTH0_NAME_CLAIM")

	cfg.Auth.Oauth2.Cognito.Audience = os.Getenv("COGNITO_CLIENT_ID")
	cfg.Auth.Oauth2.Cognito.Domain = os.Getenv("COGNITO_DOMAIN")

	cfg.Auth.Oauth2.Clerk.Domain = os.Getenv("CLERK_DOMAIN")
	cfg.Auth.Oauth2.Clerk.Audience = os.Getenv("CLERK_AUDIENCE")

	cfg.Auth.Oauth2.Supabase.Domain = os.Getenv("SUPABASE_DOMAIN")
	cfg.Auth.Oauth2.Supabase.Audience = os.Getenv("SUPABASE_AUDIENCE")

	cfg.Auth.Oauth2.Supabase.JwtSecretKey = os.Getenv("SUPABASE_JWT_SECRET_KEY")

	cfg.Aws.Sqs.QueueUrl = os.Getenv("SQS_QUEUE_URL")

	// MotifsDB = cfg.Databases.Motifs
	// WGSDB = cfg.Databases.WGS
	// GeneConvDB = cfg.Databases.GeneConv
	// PathwayDB = cfg.Databases.Pathway
	// CytobandsDir = cfg.Databases.CytobandsDir
	// GexDB = cfg.Databases.Gex
	// SeqsDB = cfg.Databases.Seqs
	// HubsDir = cfg.Databases.HubsDir
	// BedsDB = cfg.Databases.Beds
	// DnaDir = cfg.Databases.DnaDir
	// GenomesDB = cfg.Databases.Genomes
	// ScrnaDir = cfg.Databases.ScrnaDir

	//
	// Load EC keys
	//

	bytes, err := os.ReadFile("jwt.es256.private.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	cfg.Auth.Keys.JwtES256PrivateKey, err = jwt.ParseECPrivateKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("jwt.es256.public.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	cfg.Auth.Keys.JwtES256PublicKey, err = jwt.ParseECPublicKeyFromPEM(bytes)
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

	cfg.Auth.Oauth2.Auth0.JwtRsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("clerk.key.pem")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	cfg.Auth.Oauth2.Clerk.JwtRsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	cfg.App.Version = sys.Must(sys.LoadVersionInfo("version.json"))

	return &cfg
}

func Instance() *AppConfig {
	once.Do(func() {
		cfg = load()
	})

	return cfg
}
