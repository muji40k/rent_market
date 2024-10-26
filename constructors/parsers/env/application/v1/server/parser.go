package parser

import (
	"rent_service/constructors/application/v1/realisations/server"
	avariables "rent_service/constructors/parsers/env/application/v1/server/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (server.Config, error) {
	var err error
	var out server.Config

	out.Host, _ = variables.GetOr(
		avariables.SERVER_HOST, "0.0.0.0", variables.ParseString,
	)
	out.Port, err = variables.GetOr(
		avariables.SERVER_PORT, 80, variables.ParseUint,
	)

	return out, err
}

