package photoregistry

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=registry.go -destination=implementations/mock/registry.go

type IRegistry interface {
	SaveTempPhoto(tempId uuid.UUID, content []byte) error
	// Move temp photo to persistent storage and return new models.Photo
	// identifier
	MoveFromTemp(tempId uuid.UUID) (uuid.UUID, error)
	ConvertPath(path string) string
}

func MoveFromTemps(self IRegistry, tempIds ...uuid.UUID) ([]uuid.UUID, error) {
	var err error
	out := make([]uuid.UUID, 0, len(tempIds))

	for i := 0; nil == err && len(tempIds) > i; i++ {
		id, cerr := self.MoveFromTemp(tempIds[i])

		if nil == cerr {
			out = append(out, id)
		} else {
			err = cerr
		}
	}

	return out, err
}

