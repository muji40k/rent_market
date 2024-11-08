package pickuppoint

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/pickuppoint/repository.go

type IRepository interface {
	GetById(pickUpPointId uuid.UUID) (models.PickUpPoint, error)
	GetAll() (collection.Collection[models.PickUpPoint], error)
}

type IPhotoRepository interface {
	GetById(pickUpPointId uuid.UUID) (collection.Collection[uuid.UUID], error)
}

type IWorkingHoursRepository interface {
	GetById(pickUpPointId uuid.UUID) (models.PickUpPointWorkingHours, error)
}

