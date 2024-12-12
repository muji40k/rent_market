package authenticators

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type PlainOptions struct {
	Email    string
	Identity string
	Username string
	Password string
}

type Plain struct {
	opts PlainOptions
}

func NewPlain(opts PlainOptions) *Plain {
	opts.trim()
	return &Plain{opts}
}

func (self *PlainOptions) trim() {
	if "" == self.Username {
		self.Username = self.Email
	}
}

func (self *Plain) Auth(client *smtp.Client) (string, error) {
	return self.opts.Email, client.Auth(
		sasl.NewPlainClient(
			self.opts.Identity,
			self.opts.Username,
			self.opts.Password,
		),
	)
}

