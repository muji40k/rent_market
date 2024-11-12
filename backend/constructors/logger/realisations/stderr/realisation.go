package stderr

import (
	"os"
	builder "rent_service/builders/logger/writer"
	constructor "rent_service/constructors/logger"
	"rent_service/logger"
)

type Parser func() (Config, error)

type Config struct {
	Hostname string
}

func New(parser Parser) constructor.Provider {
	return func() (string, constructor.Realisation) {
		return "stderr", newConstructor(parser)
	}
}

func newConstructor(parser Parser) constructor.Realisation {
	return func() (logger.ILogger, error) {
		var log logger.ILogger

		var hostname string
		conf, err := parser()

		if nil == err {
			if name, cerr := os.Hostname(); nil == cerr {
				hostname = name
			} else {
				hostname = conf.Hostname
			}
		}

		if nil == err {
			log, err = builder.New().
				WithWriter(os.Stderr).
				WithHost(&hostname).
				Build()
		}

		return log, err
	}
}

