package authenticationroutes

import (
	"bytes"
	"fmt"
	"html/template"

	"net/mail"
	"net/url"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-auth/tokengen"
	"github.com/antonybholmes/go-auth/userdbcache"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/rdb"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-mailer"
	"github.com/antonybholmes/go-mailer/mailserver"
	"github.com/gin-gonic/gin"
)

const JWT_PARAM = "jwt"
const URL_PARAM = "url"

type EmailBody struct {
	Name       string
	From       string
	Time       string
	Link       string
	DoNotReply string
}

func SendEmailWithToken(subject string,
	authUser *auth.AuthUser,
	file string,
	token string,
	callbackUrl string,
	vistUrl string) error {

	address, err := mail.ParseAddress(authUser.Email)

	if err != nil {
		return err
	}

	return BaseSendEmailWithToken(subject, authUser, address, file, token, callbackUrl, vistUrl)
}

// Generic method for sending an email with a token in it. For APIS this is a token to use in the request, for websites
// it can craft a callback url with the token added as a parameter so that the web app can deal with the response.
func BaseSendEmailWithToken(subject string,
	authUser *auth.AuthUser,
	address *mail.Address,
	file string,
	token string,
	callbackUrl string,
	vistUrl string) error {

	var body bytes.Buffer

	t, err := template.ParseFiles(file)

	if err != nil {
		return err
	}

	var firstName string = ""

	if len(authUser.FirstName) > 0 {
		firstName = authUser.FirstName
	} else {
		firstName = strings.Split(address.Address, "@")[0]
	}

	firstName = strings.Split(firstName, " ")[0]

	time := fmt.Sprintf("%d minutes", int(auth.TTL_10_MINS.Minutes()))

	if callbackUrl != "" {
		callbackUrl, err := url.Parse(callbackUrl)

		if err != nil {
			return err
		}

		params, err := url.ParseQuery(callbackUrl.RawQuery)

		if err != nil {
			return err
		}

		if vistUrl != "" {
			params.Set(URL_PARAM, vistUrl)
		}

		params.Set(JWT_PARAM, token)

		callbackUrl.RawQuery = params.Encode()

		link := callbackUrl.String()

		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       link,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return err
		}
	} else {
		err = t.Execute(&body, EmailBody{
			Name:       firstName,
			Link:       token,
			From:       consts.NAME,
			Time:       time,
			DoNotReply: consts.DO_NOT_REPLY,
		})

		if err != nil {
			return err
		}
	}

	err = mailserver.SendHtmlEmail(address, subject, body.String())

	if err != nil {
		return err
	}

	return nil
}

func EmailUpdatedResp(c *gin.Context) {
	routes.MakeOkPrettyResp(c, "email updated")
}

// Start passwordless login by sending an email
func SendResetEmailEmailRoute(c *gin.Context) {
	return NewValidator(c).ParseLoginRequestBody().LoadAuthUserFromToken().Success(func(validator *Validator) {
		authUser := validator.AuthUser
		req := validator.LoginBodyReq

		newEmail, err := mail.ParseAddress(req.Email)

		if err != nil {
			return
		}

		otpToken, err := tokengen.ResetEmailToken(c, authUser, newEmail)

		if err != nil {
			return
		}

		// var file string

		// if req.CallbackUrl != "" {
		// 	file = "templates/email/email/reset/web.html"
		// } else {
		// 	file = "templates/email/email/reset/api.html"
		// }

		// go BaseSendEmailWithToken("Update Email",
		// 	authUser,
		// 	newEmail,
		// 	file,
		// 	otpToken,
		// 	req.CallbackUrl,
		// 	req.VisitUrl)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:          authUser.Email,
			Token:       otpToken,
			EmailType:   mailer.REDIS_EMAIL_TYPE_EMAIL_RESET,
			Ttl:         fmt.Sprintf("%d minutes", int(consts.SHORT_TTL_MINS.Minutes())),
			CallBackUrl: req.CallbackUrl,
			VisitUrl:    req.VisitUrl}
		rdb.PublishEmail(&email)

		//if err != nil {
		//	return routes.ErrorReq(err)
		//}

		routes.MakeOkPrettyResp(c, "check your email for a change email link")
	})
}

func UpdateEmailRoute(c *gin.Context) {
	return NewValidator(c).CheckEmailIsWellFormed().LoadAuthUserFromToken().Success(func(validator *Validator) error {

		if validator.Claims.Type != auth.CHANGE_EMAIL_TOKEN {
			return routes.WrongTokentTypeReq()
		}

		err := auth.CheckOTPValid(validator.AuthUser, validator.Claims.OneTimePasscode)

		if err != nil {
			return err
		}

		authUser := validator.AuthUser
		uuid := authUser.Uuid

		err = userdbcache.SetEmailAddress(authUser, validator.Address, false)

		if err != nil {
			return err
		}

		authUser, err = userdbcache.FindUserByUuid(uuid)

		if err != nil {
			return err
		}

		//return SendEmailChangedEmail(c, authUser)

		email := mailer.RedisQueueEmail{Name: authUser.FirstName,
			To:        authUser.Email,
			EmailType: mailer.REDIS_EMAIL_TYPE_EMAIL_UPDATED}
		rdb.PublishEmail(&email)

		return routes.MakeOkPrettyResp(c, "email updated confirmation email sent")
	})
}

func SendEmailChangedEmail(c *gin.Context, authUser *auth.AuthUser) error {

	file := "templates/email/email/updated.html"

	go SendEmailWithToken("Email Address Changed",
		authUser,
		file,
		"",
		"",
		"")

	//if err != nil {
	//	return routes.ErrorReq(err)
	//}

	return EmailUpdatedResp(c)

}
