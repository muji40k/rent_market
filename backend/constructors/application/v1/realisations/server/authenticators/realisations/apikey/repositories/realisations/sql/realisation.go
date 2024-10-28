package sql

import (
	builder "rent_service/builders/application/web/authenticator/apikey/repository/psql"
	constructor "rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories"
	"rent_service/server/authenticator/implementations/apikey"
)

type Parser func() (Config, error)

type Config struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
}

func New(parser Parser) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "sql", newConstructor(parser)
	}
}

func newConstructor(parser Parser) constructor.Realisation {
	return func() (apikey.ITokenRepository, error) {
		var repo apikey.ITokenRepository
		conf, err := parser()

		if nil == err {
			repo, err = builder.New().
				WithUser(conf.UserName).
				WithPassword(conf.Password).
				WithHost(conf.Host).
				WithPort(conf.Port).
				WithDatabase(conf.Database).
				Build()
		}

		return repo, err
	}
}

