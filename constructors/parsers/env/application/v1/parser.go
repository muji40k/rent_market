package v1

import (
	"rent_service/constructors/application/v1"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (v1.Config, error) {
	var err error
	var out v1.Config

	out.Type, _ = variables.GetOr(
		variables.APP_TYPE,
		"web",
		variables.ParseString,
	)

	return out, err
}

