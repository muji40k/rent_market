package psql

import (
	"fmt"
	"os"
	builder "rent_service/builders/repository/factory/v1/psql"
	constructor "rent_service/constructors/repository/factory/v1"
	v1 "rent_service/internal/factory/repositories/v1"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical/implementations/simple"
)

type Parser func() (Config, error)

type Config struct {
	Hostname string
	Hasher   string
	DB       struct {
		UserName string
		Password string
		Host     string
		Port     string
		Database string
	}
}

func New(parser Parser, hashers map[string]user.Hasher) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "psql", newConstructor(parser, hashers)
	}
}

func newConstructor(parser Parser, hashers map[string]user.Hasher) constructor.Realisation {
	return func() (v1.IFactory, error) {
		var rfactory v1.IFactory
		var hostname string
		var hasher user.Hasher
		conf, err := parser()

		if nil == err {
			if name, cerr := os.Hostname(); nil == cerr {
				hostname = name
			} else {
				hostname = conf.Hostname
			}

			if h, found := hashers[conf.Hasher]; !found {
				err = fmt.Errorf("Hasher '%v' not found", conf.Hasher)
			} else {
				hasher = h
			}
		}

		if nil == err {
			rfactory, err = builder.New().
				WithUser(conf.DB.UserName).
				WithPassword(conf.DB.Password).
				WithHost(conf.DB.Host).
				WithPort(conf.DB.Port).
				WithDatabase(conf.DB.Database).
				WithHasher(hasher).
				WithSetter(simple.New(fmt.Sprintf("backend@%v", hostname))).
				Build()
		}

		return rfactory, err
	}
}

