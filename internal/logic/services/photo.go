package services

import (
	"fmt"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
)

type PhotoDescription struct {
	Mime        string
	Placeholder string
	Description string
}

type IPhotoService interface {
	CreateTempPhoto(token Token, photo PhotoDescription) (uuid.UUID, error)
	UploadTempPhoto(token Token, photoId uuid.UUID, content []byte) error
	SavePhotoFromTemp(token Token, photoId uuid.UUID) (models.Photo, error)

	GetPhoto(token Token, photoId uuid.UUID) (models.Photo, error)
}

type ErrorUnsupportedMime struct{ Mime string }

func (e ErrorUnsupportedMime) Error() string {
	return fmt.Sprintf("Unsupported mime for the photo '%v'", e.Mime)
}

