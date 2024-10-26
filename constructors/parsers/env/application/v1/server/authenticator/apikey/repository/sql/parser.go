package sql

import (
	"rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey/repositories/realisations/sql"
	svariables "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey/repository/sql/variables"
	"rent_service/constructors/parsers/env/variables"
)

func Parser() (sql.Config, error) {
	var config sql.Config

	config.UserName, _ = variables.GetOr(
		svariables.DB_USERNAME, "postgres", variables.ParseString,
	)
	config.Password, _ = variables.GetOr(
		svariables.DB_PASSWORD, "postgres", variables.ParseString,
	)
	config.Host, _ = variables.GetOr(
		svariables.DB_HOST, "0.0.0.0", variables.ParseString,
	)
	config.Port, _ = variables.GetOr(
		svariables.DB_PORT, "5432", variables.ParseString,
	)
	config.Database, _ = variables.GetOr(
		svariables.DB_NAME, "authentication", variables.ParseString,
	)

	return config, nil
}

