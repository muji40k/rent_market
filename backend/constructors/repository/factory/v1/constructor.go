package v1

import (
	"fmt"
	"rent_service/constructors"
	v1 "rent_service/internal/factory/repositories/v1"
)

type Realisation func() (v1.IFactory, error)
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

func (self *Constructor) Construct(
	cleaner *constructors.Cleaner,
) (v1.IFactory, error) {
	var factory v1.IFactory
	config, err := self.parser()

	if nil == err {
		if realisation, found := self.realisations[config.Type]; found {
			factory, err = realisation()
		} else {
			err = fmt.Errorf("Can't find repository of type '%v'", config.Type)
		}
	}

	if nil == err && nil != cleaner {
		cleaner.AddStage(factory.Clear)
	}

	return factory, err
}

