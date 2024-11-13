package logger

import (
	"fmt"
	"rent_service/constructors"
	"rent_service/logger"
)

type Realisation func() (logger.ILogger, error)
type Provider func() (string, Realisation)

type Parser func() (Config, error)

type Config struct {
	Type string
}

type entry struct {
	realisation Realisation
	logger      logger.ILogger
}

type Constructor struct {
	parser       Parser
	realisations map[string]entry
}

func New(parser Parser, realisations ...Provider) *Constructor {
	out := &Constructor{parser, make(map[string]entry)}

	for _, provider := range realisations {
		name, realisation := provider()
		out.realisations[name] = entry{realisation, nil}
	}

	return out
}

func (self *Constructor) Construct(
	cleaner *constructors.Cleaner,
) (logger.ILogger, error) {
	var log logger.ILogger
	config, err := self.parser()

	if "" == config.Type {
		return nil, nil
	}

	if nil == err {
		if entry, found := self.realisations[config.Type]; found {
			if nil != entry.logger {
				log = entry.logger
			} else {
				log, err = entry.realisation()
				entry.logger = log
				self.realisations[config.Type] = entry

				if nil == err && nil != cleaner {
					cleaner.AddStage(log.Close)
				}
			}
		} else {
			err = fmt.Errorf("Can't find logger of type '%v'", config.Type)
		}
	}

	return log, err
}

