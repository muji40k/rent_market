package stderr

import (
	"rent_service/constructors/logger/realisations/stderr"
	rvariables "rent_service/constructors/parsers/env/logger/stderr/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (stderr.Config, error) {
	var out stderr.Config
	var err error

	out.Hostname, _ = variables.GetOr(
		rvariables.APP_HOSTNAME, "unknown", variables.ParseString,
	)

	return out, err
}

