package email

import (
	"fmt"
	"io"
	"rent_service/internal/domain/models"
	"time"

	"github.com/emersion/go-smtp"
)

type IClientProvider interface {
	GetClient() (*smtp.Client, error)
}

type IAuthenticator interface {
	Auth(client *smtp.Client) (string, error)
}

type Email2FA struct {
	client        IClientProvider
	authenticator IAuthenticator
}

func New(client IClientProvider, auth IAuthenticator) *Email2FA {
	return &Email2FA{client, auth}
}

const FORMAT string = "Mon, 02 Jan 2006 15:04:05 -0700"

func SendMail(
	c *smtp.Client,
	from string,
	to string,
	subject string,
	body func(io.Writer) error,
) error {
	var wc io.WriteCloser
	err := c.Mail(from, nil)

	if nil == err {
		err = c.Rcpt(to, nil)
	}

	if nil == err {
		wc, err = c.Data()
	}

	if nil == err {
		_, err = fmt.Fprintf(wc, "Date: %v\n"+
			"From: %v\n"+
			"To: %v\n"+
			"Subject: %v\n\n",
			time.Now().Format(FORMAT), from, to, subject,
		)
	}

	if nil == err {
		err = body(wc)
	}

	if nil != wc {
		cerr := wc.Close()

		if nil != cerr {
			if nil == err {
				err = cerr
			} else {
				err = fmt.Errorf(
					"General error: %w\nClose error: %w", err, cerr,
				)
			}
		}
	}

	return err
}

func (self *Email2FA) SendCode(user models.User, code string) error {
	var from string
	c, err := self.client.GetClient()

	if nil == err {
		from, err = self.authenticator.Auth(c)
	}

	if nil == err {
		err = SendMail(c, from, user.Email,
			"Подтверждение смены пароля",
			func(w io.Writer) error {
				_, err := fmt.Fprintf(w, "Здравствуйте, %v!\n"+
					"Ваш код подтверждения: %v.\n\n"+
					"Никому не сообщайте код подтверждения.",
					user.Name, code,
				)

				return err
			},
		)
	}

	if nil != c {
		cerr := c.Quit()

		if nil != cerr {
			if nil == err {
				err = cerr
			} else {
				err = fmt.Errorf(
					"General error: %w\nQuit error: %w", err, cerr,
				)
			}
		}
	}

	return err
}

