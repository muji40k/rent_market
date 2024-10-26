package defservices

import (
	svariables "rent_service/constructors/parsers/env/service/v1/defservices/variables"
	"rent_service/constructors/parsers/env/variables"
	"rent_service/constructors/service/factory/v1/realisations/defservices"
)

func Parser() (defservices.Config, error) {
	var out defservices.Config
	var err error

	out.CodegenLength, err = variables.GetOr(
		svariables.CODE_LENGTH, 6, variables.ParseUint,
	)

	out.Photo.Main, _ = variables.GetOr(
		svariables.PHOTO_MAIN_PATH, "/server/media", variables.ParseString,
	)
	out.Photo.Temp, _ = variables.GetOr(
		svariables.PHOTO_TEMP_PATH, "/server/temp", variables.ParseString,
	)
	out.Photo.BaseUrl, _ = variables.GetOr(
		svariables.PHOTO_BASE_URL, "http://localhost/static", variables.ParseString,
	)

	return out, err
}

