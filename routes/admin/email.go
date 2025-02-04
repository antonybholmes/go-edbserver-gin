package adminroutes

import (
	"bytes"
	"html/template"

	"net/mail"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-mailer/mailserver"
)

const FILE = "templates/email/account/created.html"

type EmailBody struct {
	Name       string
	From       string
	Time       string
	Link       string
	DoNotReply string
}

func SendAccountCreatedEmail(
	authUser *auth.AuthUser,
	address *mail.Address) error {

	var body bytes.Buffer

	t, err := template.ParseFiles(FILE)

	if err != nil {
		c.Error(err)
		return
	}

	var firstName string = ""

	if len(authUser.FirstName) > 0 {
		firstName = authUser.FirstName
	} else {
		firstName = strings.Split(address.Address, "@")[0]
	}

	firstName = strings.Split(firstName, " ")[0]

	err = t.Execute(&body, EmailBody{
		Name:       firstName,
		Link:       consts.APP_URL,
		From:       consts.NAME,
		DoNotReply: consts.DO_NOT_REPLY,
	})

	if err != nil {
		c.Error(err)
		return
	}

	err = mailserver.SendHtmlEmail(address, "New account created", body.String())

	if err != nil {
		c.Error(err)
		return
	}

	return nil
}
