package variables

import (
	"os"
	"strconv"
)

const (
	APP_TYPE        string = "APP_TYPE"
	SERVICE_TYPE    string = "APP_SERVICE_TYPE"
	REPOSITORY_TYPE string = "APP_REPOSITORY_TYPE"
)

func GetOr[T any](name string, def T, parser func(string) (T, error)) (T, error) {
	if value := os.Getenv(name); "" == value {
		return def, nil
	} else {
		return parser(value)
	}
}

func ParseString(value string) (string, error) {
	return value, nil
}

func ParseUint(value string) (uint, error) {
	v, e := strconv.ParseUint(value, 10, 64)
	return uint(v), e
}

