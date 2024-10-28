package v1

import (
	"rent_service/constructors/parsers/env/variables"
	"rent_service/constructors/service/factory/v1"
)

func Parser() (v1.Config, error) {
	var err error
	var out v1.Config

	out.Type, err = variables.GetOr(
		variables.SERVICE_TYPE,
		"default",
		variables.ParseString,
	)

	return out, err
}

