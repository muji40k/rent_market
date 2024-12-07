package db

import "os"

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

const (
	TEST_HOST     string = "TEST_DB_HOST"
	TEST_PORT     string = "TEST_DB_PORT"
	TEST_DATABASE string = "TEST_DB_NAME"
	TEST_USER     string = "TEST_DB_USERNAME"
	TEST_PASSWORD string = "TEST_DB_PASSWORD"
)

func getOr(variable string, def string) string {
	if value := os.Getenv(variable); "" != value {
		return value
	} else {
		return def
	}
}

func FromEnv() Config {
	return Config{
		Host:     getOr(TEST_HOST, "localhost"),
		Port:     getOr(TEST_PORT, "5432"),
		Database: getOr(TEST_DATABASE, "rent_market"),
		User:     getOr(TEST_USER, "postgres"),
		Password: getOr(TEST_PASSWORD, "postgres"),
	}
}

