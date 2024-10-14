package photo

import (
	"fmt"

	"github.com/google/uuid"

	"rent_service/internal/logic/services/types/token"
)

type IService interface {
	CreateTempPhoto(token token.Token, photo Description) (uuid.UUID, error)
	UploadTempPhoto(token token.Token, photoId uuid.UUID, content []byte) error

	GetTempPhoto(token token.Token, photoId uuid.UUID) (TempPhoto, error)
	GetPhoto(photoId uuid.UUID) (Photo, error)
}

type ErrorUnsupportedMime struct{ Mime string }

func (e ErrorUnsupportedMime) Error() string {
	return fmt.Sprintf("Unsupported mime for the photo '%v'", e.Mime)
}

