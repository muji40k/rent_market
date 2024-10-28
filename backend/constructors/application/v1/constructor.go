package v1

import (
	"fmt"
	"rent_service/application"
	"rent_service/constructors"
	scontext "rent_service/internal/logic/context/v1"
)

type Realisation func(context *scontext.Context) (application.IApplication, error)
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
	context *scontext.Context,
	cleaner *constructors.Cleaner,
) (application.IApplication, error) {
	var application application.IApplication
	config, err := self.parser()

	if nil == err {
		if realisation, found := self.realisations[config.Type]; found {
			application, err = realisation(context)
		} else {
			err = fmt.Errorf("Can't find app of type '%v'", config.Type)
		}
	}

	if nil == err && nil != cleaner {
		cleaner.AddStage(application.Clear)
	}

	return application, err
}

