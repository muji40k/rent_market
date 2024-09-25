package services

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IPickUpPointService interface {
	ListPickUpPoints() (Collection[models.PickUpPoint], error)
	GetPickUpPointById(pickUpPointId uuid.UUID) (models.PickUpPoint, error)
	ListPickUpPointPhotos(
		pickUpPointId uuid.UUID,
	) (Collection[uuid.UUID], error)
	ListPickUpPointWorkingHours(
		pickUpPointId uuid.UUID,
	) (Collection[models.PickUpPointWorkingHours], error)
}

