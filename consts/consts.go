package consts

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-sys/env"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web/auth"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Name      = "Experiments Server"
	AppName   = "edbserver"
	Copyright = "Copyright (C) 2024-2025 Antony Holmes"
)

var (
	AppUrl    string
	AppDomain string
	Version   sys.VersionInfo

	JwtRsaPrivateKey     *rsa.PrivateKey //[]byte
	JwtRsaPublicKey      *rsa.PublicKey  //[]byte
	JwtAuth0RsaPublicKey *rsa.PublicKey
	JwtClerkRsaPublicKey *rsa.PublicKey
	JwtSupabaseSecretKey string
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

	CognitoClientId string
	CognitoDomain   string

	PasswordlessTokenTtlMins time.Duration
	AccessTokenTtlMins       time.Duration
	OtpTokenTtlMins          time.Duration
	ShortTtlMins             time.Duration

	UrlResetEmail    string
	UrlResetPassword string
	UrlVerifyEmail   string

	SqsQueueUrl string
)

func init() {
	env.Load("consts.env")
	env.Load("version.env")

	AppUrl = os.Getenv("APP_URL")
	AppDomain = os.Getenv("APP_DOMAIN")
	//Version = os.Getenv("VERSION")
	//Updated = os.Getenv("UPDATED")

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	UrlResetEmail = os.Getenv("URL_RESET_EMAIL")
	UrlResetPassword = os.Getenv("URL_RESET_PASSWORD")
	UrlVerifyEmail = os.Getenv("URL_VERIFY_EMAIL")

	//JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	//JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SessionName = os.Getenv("SESSION_NAME")
	SessionKey = os.Getenv("SESSION_KEY")
	SessionEncryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")

	PasswordlessTokenTtlMins = env.GetMin("PASSWORDLESS_TOKEN_TTL_MINS", auth.Ttl10Mins)
	AccessTokenTtlMins = env.GetMin("ACCESS_TOKEN_TTL_MINS", auth.Ttl15Mins)
	OtpTokenTtlMins = env.GetMin("OTP_TOKEN_TTL_MINS", auth.Ttl20Mins)
	ShortTtlMins = env.GetMin("SHORT_TTL_MINS", auth.Ttl10Mins)

	Auth0Audience = os.Getenv("AUTH0_AUDIENCE")
	Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	Auth0EmailClaim = os.Getenv("AUTH0_EMAIL_CLAIM")
	Auth0NameClaim = os.Getenv("AUTH0_NAME_CLAIM")

	CognitoClientId = os.Getenv("COGNITO_CLIENT_ID")
	CognitoDomain = os.Getenv("COGNITO_DOMAIN")

	JwtSupabaseSecretKey = os.Getenv("JWT_SUPABASE_SECRET_KEY")

	SqsQueueUrl = os.Getenv("SQS_QUEUE_URL")

	bytes, err := os.ReadFile("jwtRS256.key")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtRsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	bytes, err = os.ReadFile("jwtRS256.key.pub")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	JwtRsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

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
