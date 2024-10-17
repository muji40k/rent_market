package uuid

import (
	"errors"
	"rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GenerateAvailable(
	db *sqlx.DB,
	checker func(*sqlx.DB, uuid.UUID) error,
) (uuid.UUID, error) {
	var id uuid.NullUUID
	var err error

	for nil == err && id.Valid {
		id.UUID, err = uuid.NewRandom()

		if nil == err {
			err = checker(db, id.UUID)

			if cerr := (cmnerrors.ErrorNotFound{}); errors.As(err, &cerr) {
				err = nil
				id.Valid = true
			}
		}
	}

	return id.UUID, err
}

