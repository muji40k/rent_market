package authenticator

import (
	"rent_service/constructors/application/v1/realisations/server/authenticators"
	svariables "rent_service/constructors/parsers/env/application/v1/server/authenticator/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (authenticators.Config, error) {
	var config authenticators.Config

	config.Type, _ = variables.GetOr(
		svariables.AUTH_TYPE, "apikey", variables.ParseString,
	)

	return config, nil
}

