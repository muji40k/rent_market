package photo

import (
	"fmt"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
)

type Description struct {
	Mime        string
	Placeholder string
	Description string
}

type IService interface {
	CreateTempPhoto(token models.Token, photo Description) (uuid.UUID, error)
	UploadTempPhoto(
		token models.Token,
		photoId uuid.UUID,
		content []byte,
	) error
	SavePhotoFromTemp(
		token models.Token,
		photoId uuid.UUID,
	) (models.Photo, error)

	GetPhoto(token models.Token, photoId uuid.UUID) (models.Photo, error)
}

type ErrorUnsupportedMime struct{ Mime string }

func (e ErrorUnsupportedMime) Error() string {
	return fmt.Sprintf("Unsupported mime for the photo '%v'", e.Mime)
}

