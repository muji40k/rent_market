package apikey

import (
	builder "rent_service/builders/application/web/authenticator/apikey"
	"rent_service/constructors"
	constructor "rent_service/constructors/application/v1/realisations/server/authenticators"
	rconstructor "rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories"
	"rent_service/internal/logic/context/providers/login"
	"rent_service/internal/logic/services/types/token"
	"rent_service/server/authenticator"
	"rent_service/server/authenticator/implementations/apikey"
	"time"
)

type Config struct {
	AccessTime time.Duration
	RenewTime  time.Duration
}

type Parser func() (Config, error)

func New(parser Parser, repo *rconstructor.Constructor) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "apikey", newConstructor(parser, repo)
	}
}

type cleanerWrap struct {
	authenticator authenticator.IAuthenticator
	cleaner       *constructors.Cleaner
}

func (self *cleanerWrap) Login(email string, password string) (authenticator.ApiToken, error) {
	return self.authenticator.Login(email, password)
}

func (self *cleanerWrap) GetToken(access string) (token.Token, error) {
	return self.authenticator.GetToken(access)
}

func (self *cleanerWrap) RenewKey(token authenticator.ApiToken) (authenticator.ApiToken, error) {
	return self.authenticator.RenewKey(token)
}

func (self *cleanerWrap) Logout(access string) error {
	return self.authenticator.Logout(access)
}

func (self *cleanerWrap) Clear() {
	self.authenticator.Clear()
	self.cleaner.Clean()
}

func newConstructor(
	parser Parser,
	constructor *rconstructor.Constructor,
) constructor.Realisation {
	return func(login login.IProvider) (authenticator.IAuthenticator, error) {
		var out = &cleanerWrap{nil, constructors.NewCleaner()}
		var repo apikey.ITokenRepository
		config, err := parser()

		if nil == err {
			repo, err = constructor.Construct(out.cleaner)
		}

		if nil == err {
			out.authenticator, err = builder.New().
				WithAccesTime(config.AccessTime).
				WithRenewTime(config.RenewTime).
				WithLogin(login).
				WithTokenRepository(repo).
				Build()
		}

		return out, nil
	}
}

