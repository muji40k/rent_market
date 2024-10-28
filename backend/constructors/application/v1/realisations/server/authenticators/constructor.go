package authenticators

import (
	"fmt"
	"rent_service/constructors"
	"rent_service/internal/logic/context/providers/login"
	"rent_service/server/authenticator"
)

type Realisation func(login login.IProvider) (authenticator.IAuthenticator, error)
type Provider func() (string, Realisation)

type Parser func() (Config, error)

type Config struct {
	Type string
}

type Constructor struct {
	parser       Parser
	realisations map[string]Realisation
}

func New(parser Parser, realisations ...Provider) *Constructor {
	out := &Constructor{parser, make(map[string]Realisation)}

	for _, provider := range realisations {
		name, realisation := provider()
		out.realisations[name] = realisation
	}

	return out
}

func (self *Constructor) Construct(login login.IProvider, cleaner *constructors.Cleaner) (authenticator.IAuthenticator, error) {
	var auth authenticator.IAuthenticator
	config, err := self.parser()

	if nil == err {
		if realisation, found := self.realisations[config.Type]; found {
			auth, err = realisation(login)
		} else {
			err = fmt.Errorf("Can't find authenticator of type '%v'", config.Type)
		}
	}

	if nil == err && nil != cleaner {
		cleaner.AddStage(auth.Clear)
	}

	return auth, err
}

