package uuid

import (
	"errors"

	"github.com/google/uuid"
)

func Parse(idRaw string) (uuid.UUID, error) {
	var err error
	var id uuid.UUID

	if "" == idRaw {
		err = errors.New("No identifier provided...")
	} else if parsed, cerr := uuid.Parse(idRaw); nil == cerr {
		id = parsed
	} else {
		err = cerr
	}

	return id, err
}

