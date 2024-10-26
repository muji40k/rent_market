package v1

import (
	"rent_service/constructors/parsers/env/variables"
	"rent_service/constructors/repository/factory/v1"
)

func Parser() (v1.Config, error) {
	var err error
	var out v1.Config

	out.Type, _ = variables.GetOr(
		variables.REPOSITORY_TYPE,
		"psql",
		variables.ParseString,
	)

	return out, err
}

