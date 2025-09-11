package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const OTP_TTL = 10 * time.Minute
const KEY = "otp:"

type OTP struct {
	Context     context.Context
	RedisClient *redis.Client
}

func NewOTP(rdb *redis.Client) *OTP {
	return &OTP{
		Context:     context.Background(),
		RedisClient: rdb,
	}
}

func (otp *OTP) CacheDigitCode(username string) (string, error) {
	code, err := auth.Generate6DigitCode()

	if err != nil {
		return "", err
	}

	err = otp.store2FACode(username, code)

	if err != nil {
		return "", err
	}

	return code, nil
}

func (otp *OTP) delete2FACode(username string) error {
	key := KEY + username
	return otp.RedisClient.Del(otp.Context, key).Err()
}

func (otp *OTP) get2FACode(username string) (string, error) {
	key := KEY + username
	return otp.RedisClient.Get(otp.Context, key).Result()
}

func (otp *OTP) store2FACode(username string, code string) error {
	key := KEY + username
	return otp.RedisClient.Set(otp.Context, key, code, OTP_TTL).Err() // expires in 5 mins
}

func (otp *OTP) validate2FACode(username string, input string) (bool, error) {

	stored, err := otp.get2FACode(username)

	log.Debug().Msgf("validating %s %s %s", username, input, stored)

	if err == redis.Nil {
		return false, nil // not found or expired
	} else if err != nil {
		return false, err
	}

	if stored != input {
		return false, nil
	}

	// Remove after use
	err = otp.delete2FACode(username)

	if err != nil {
		return false, err
	}

	return true, nil
}

type OTPRoutes struct {
	OTP *OTP
}

func NewOTPRoutes(otp *OTP) *OTPRoutes {
	return &OTPRoutes{
		OTP: otp,
	}
}

func (otpRoutes *OTPRoutes) EmailOTPRoute(c *gin.Context) {
	otpRoutes.Email6DigitCodeRoute(c)
}

func (otpRoutes *OTPRoutes) Email6DigitCodeRoute(c *gin.Context) {

	validator, err := NewValidator(c).CheckEmailIsWellFormed().Ok()

	if err != nil {
		web.BaseBadReqResp(c, err)
		return
	}

	//user := validator.AuthUser
	address := validator.Address

	code, err := otpRoutes.OTP.CacheDigitCode(address.Address)

	if err != nil {
		web.BaseInternalErrorResp(c, err)
		return
	}

	email := mailer.QueueEmail{
		Name:      address.Address,
		To:        address.Address,
		Token:     code,
		TTL:       fmt.Sprintf("%d minutes", int(OTP_TTL.Minutes())),
		EmailType: mailer.QUEUE_EMAIL_TYPE_OTP}
	err = queue.PublishEmail(&email)

	if err != nil {
		log.Debug().Msgf("error sending email %v", err)
		web.BaseInternalErrorResp(c, err)
		return
	}

	web.MakeOkResp(c, "6 digit code sent to email")
}
