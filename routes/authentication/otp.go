package authentication

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/queue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const OTP_TTL = 5 * time.Minute

type OTP struct {
	ctx context.Context
	rdb *redis.Client
}

func NewOTP(rdb *redis.Client) *OTP {
	return &OTP{
		ctx: context.Background(),
		rdb: rdb,
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

func (otp *OTP) Email6DigitCodeRoute(c *gin.Context) {

	req, err := ParseLoginRequestBody(c)

	if err != nil {
		web.BaseBadReqResp(c, err)
		return
	}

	emailAddress := req.Email

	code, err := otp.CacheDigitCode(emailAddress)

	if err != nil {
		web.BaseInternalErrorResp(c, err)
		return
	}

	address, err := mail.ParseAddress(emailAddress)

	if err != nil {
		web.BaseBadReqResp(c, err)
		return
	}

	email := mailer.QueueEmail{
		Name:      emailAddress,
		To:        address.Address,
		Token:     code,
		Ttl:       fmt.Sprintf("%d minutes", int(OTP_TTL.Minutes())),
		EmailType: mailer.QUEUE_EMAIL_TYPE_TOTP}
	queue.PublishEmail(&email)

	web.MakeOkResp(c, "6 digit code sent to email")

}

func (otp *OTP) delete2FACode(username string) error {
	key := "2fa:" + username
	return otp.rdb.Del(otp.ctx, key).Err()
}

func (otp *OTP) get2FACode(username string) (string, error) {
	key := "2fa:" + username
	return otp.rdb.Get(otp.ctx, key).Result()
}

func (otp *OTP) store2FACode(username, code string) error {
	key := "2fa:" + username
	return otp.rdb.Set(otp.ctx, key, code, OTP_TTL).Err() // expires in 5 mins
}
