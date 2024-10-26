package repositories

import (
	"fmt"
	"rent_service/constructors"
	"rent_service/server/authenticator/implementations/apikey"
)

type Realisation func() (apikey.ITokenRepository, error)
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

func (self *Constructor) Construct(cleaner *constructors.Cleaner) (apikey.ITokenRepository, error) {
	var repo apikey.ITokenRepository
	config, err := self.parser()

	if nil == err {
		if realisation, found := self.realisations[config.Type]; found {
			repo, err = realisation()
		} else {
			err = fmt.Errorf(
				"Can't find authenticator repository of type '%v'", config.Type,
			)
		}
	}

	if nil == err && nil != cleaner {
		cleaner.AddStage(repo.Clear)
	}

	return repo, err
}

