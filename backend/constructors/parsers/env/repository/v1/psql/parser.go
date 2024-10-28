package psql

import (
	rvariables "rent_service/constructors/parsers/env/repository/v1/psql/variables"
	"rent_service/constructors/parsers/env/variables"
	"rent_service/constructors/repository/factory/v1/realisations/psql"
)

func Parser() (psql.Config, error) {
	var out psql.Config
	var err error

	out.DB.UserName, _ = variables.GetOr(
		rvariables.DB_USERNAME, "postgres", variables.ParseString,
	)
	out.DB.Password, _ = variables.GetOr(
		rvariables.DB_PASSWORD, "postgres", variables.ParseString,
	)
	out.DB.Host, _ = variables.GetOr(
		rvariables.DB_HOST, "0.0.0.0", variables.ParseString,
	)
	out.DB.Port, _ = variables.GetOr(
		rvariables.DB_PORT, "5432", variables.ParseString,
	)
	out.DB.Database, _ = variables.GetOr(
		rvariables.DB_NAME, "rent_market", variables.ParseString,
	)
	out.Hostname, _ = variables.GetOr(
		rvariables.APP_HOSTNAME, "unknown", variables.ParseString,
	)
	out.Hasher, _ = variables.GetOr(
		rvariables.REPOSITORY_HASHER, "md5", variables.ParseString,
	)

	return out, err
}

