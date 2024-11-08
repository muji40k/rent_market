package photoregistry

import (
	"github.com/google/uuid"
)

type IRegistry interface {
	SaveTempPhoto(tempId uuid.UUID, content []byte) error
	// Move temp photo to persistent storage and return new models.Photo
	// identifier
	MoveFromTemp(tempId uuid.UUID) (uuid.UUID, error)
	MoveFromTemps(tempIds ...uuid.UUID) ([]uuid.UUID, error)
	ConvertPath(path string) string
}

