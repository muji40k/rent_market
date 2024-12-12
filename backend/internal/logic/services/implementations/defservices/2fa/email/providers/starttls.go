package providers

import (
	"crypto/tls"

	"github.com/emersion/go-smtp"
)

type StartTLSProvider struct {
	host    string
	port    string
	setters []func(*tls.Config)
}

func NewStartTLS(host, port string, setters ...func(*tls.Config)) *StartTLSProvider {
	cpy := make([]func(*tls.Config), len(setters))
	copy(cpy, setters)
	return &StartTLSProvider{host, port, cpy}
}

func (self *StartTLSProvider) GetClient() (*smtp.Client, error) {
	tlsconfig := &tls.Config{
		ServerName: self.host,
	}

	for _, s := range self.setters {
		s(tlsconfig)
	}

	return smtp.DialStartTLS(self.host+":"+self.port, tlsconfig)
}

