package repository

import (
	"rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories"
	svariables "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey/repository/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (repositories.Config, error) {
	var config repositories.Config

	config.Type, _ = variables.GetOr(
		svariables.REPO_TYPE, "sql", variables.ParseString,
	)

	return config, nil
}

