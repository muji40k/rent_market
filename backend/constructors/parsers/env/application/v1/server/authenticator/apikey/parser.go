package apikey

import (
	"rent_service/constructors/application/v1/realisations/server/authenticators/realisations/apikey"
	svariables "rent_service/constructors/parsers/env/application/v1/server/authenticator/apikey/variables"
	"rent_service/constructors/parsers/env/variables"
	"strconv"
	"time"
)

func parseTime(value string) (time.Duration, error) {
	v, err := strconv.ParseInt(value, 10, 64)

	if nil == err {
		return time.Duration(v) * time.Hour, nil
	} else {
		return 0, err
	}
}

func Parser() (apikey.Config, error) {
	var config apikey.Config
	var err error

	config.AccessTime, err = variables.GetOr(
		svariables.ACCESS_TIME, 24*time.Hour, parseTime,
	)

	if nil == err {
		config.RenewTime, err = variables.GetOr(
			svariables.RENEW_TIME, 7*24*time.Hour, parseTime,
		)
	}

	return config, err
}

