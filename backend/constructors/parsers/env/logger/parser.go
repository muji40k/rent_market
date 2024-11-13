package logger

import (
	"rent_service/constructors/logger"
	lvariables "rent_service/constructors/parsers/env/logger/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (logger.Config, error) {
	var err error
	var out logger.Config

	out.Type, _ = variables.GetOr(
		lvariables.LOGGER_TYPE,
		"",
		variables.ParseString,
	)

	return out, err
}

