package photo

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/photo/repository.go

type IRepository interface {
	Create(photo models.Photo) (models.Photo, error)

	GetById(photoId uuid.UUID) (models.Photo, error)
}

type ITempRepository interface {
	Create(photo models.TempPhoto) (models.TempPhoto, error)

	Update(photo models.TempPhoto) error

	GetById(photoId uuid.UUID) (models.TempPhoto, error)

	Remove(photoId uuid.UUID) error
}

