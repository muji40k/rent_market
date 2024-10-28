package apikey

import (
	"errors"
	"rent_service/internal/logic/context/providers/login"
	"rent_service/server/authenticator"
	"rent_service/server/authenticator/implementations/apikey"
	"time"
)

type Builder struct {
	login      login.IProvider
	repository apikey.ITokenRepository
	accesTime  time.Duration
	renewTime  time.Duration
}

func New() *Builder {
	return &Builder{}
}

func (self *Builder) WithLogin(login login.IProvider) *Builder {
	self.login = login
	return self
}

func (self *Builder) WithTokenRepository(repository apikey.ITokenRepository) *Builder {
	self.repository = repository
	return self
}

func (self *Builder) WithAccesTime(accesTime time.Duration) *Builder {
	self.accesTime = accesTime
	return self
}

func (self *Builder) WithRenewTime(renewTime time.Duration) *Builder {
	self.renewTime = renewTime
	return self
}

func (self *Builder) Build() (authenticator.IAuthenticator, error) {
	var err error
	var out authenticator.IAuthenticator

	if nil == self.login {
		err = errors.New("ApiKeyBuilder: login provider not set")
	}

	if nil == err && nil == self.repository {
		err = errors.New("ApiKeyBuilder: token repository not set")
	}

	if nil == err {
		out = apikey.New(
			self.login,
			self.repository,
			self.accesTime,
			self.renewTime,
		)
	}

	return out, err
}

