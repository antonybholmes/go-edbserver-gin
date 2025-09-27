package authentication

import (
	"fmt"
	"math"

	edbmail "github.com/antonybholmes/go-edbmailserver/mail"
	mailserver "github.com/antonybholmes/go-mailserver"
	"github.com/antonybholmes/go-mailserver/mailqueue"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)



type OTPRoutes struct {
	OTP *auth.OTP
}

func NewOTPRoutes(otp *auth.OTP) *OTPRoutes {
	return &OTPRoutes{
		OTP: otp,
	}
}

func (otpRoutes *OTPRoutes) Email6DigitOTPRoute(c *gin.Context) {

	validator, err := NewValidator(c).CheckEmailIsWellFormed().Ok()

	if err != nil {
		web.BaseBadReqResp(c, err)
		return
	}

	//log.Debug().Msgf("EmailOTPRoute")

	//user := validator.AuthUser
	address := validator.Address

	code, err := otpRoutes.OTP.Cache6DigitOTP(address.Address)

	if err != nil {
		web.BaseInternalErrorResp(c, err)
		return
	}

	mins :=  int(math.Round(otpRoutes.OTP.TTL().Minutes()))  

	email := mailserver.MailItem{
		Name:      address.Address,
		To:        address.Address,
		Payload:   &mailserver.Payload{DataType: "code", Data: code},
		TTL:       fmt.Sprintf("%d minutes", mins),
		EmailType: edbmail.QUEUE_EMAIL_TYPE_OTP,
	}
	err = mailqueue.SendMail(&email)

	if err != nil {
		log.Debug().Msgf("error sending email %v", err)
		web.BaseInternalErrorResp(c, err)
		return
	}

	web.MakeOkResp(c, fmt.Sprintf("A 6 digit one-time code has been sent to the email address. The code is valid for %d minutes.", mins))
}
