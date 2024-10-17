package exist

import (
	"fmt"
	"rent_service/internal/repository/errors/cmnerrors"

	"github.com/jmoiron/sqlx"
)

func Check(what string, db *sqlx.DB, counter string, args ...interface{}) error {
	var amount uint
	err := db.Get(&amount, counter, args...)

	if nil == err && 0 == amount {
		err = cmnerrors.NotFound(what)
	} else if nil == err && 1 != amount {
		err = fmt.Errorf("Found multiple instances of %v ...", what)
	}

	return err
}

func CheckMultiple(what string, db *sqlx.DB, counter string, args ...interface{}) error {
	var amount uint
	err := db.Get(&amount, counter, args...)

	if nil == err && 0 == amount {
		err = cmnerrors.NotFound(what)
	}

	return err
}

func GenericCounter(tableName string) string {
	return fmt.Sprintf("select count(*) from %v where id = $1", tableName)
}

