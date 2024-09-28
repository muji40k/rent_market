package pickuppoint

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListPickUpPoints() (Collection[models.PickUpPoint], error)
	GetPickUpPointById(pickUpPointId uuid.UUID) (models.PickUpPoint, error)
}

type IPhotoService interface {
	ListPickUpPointPhotos(
		pickUpPointId uuid.UUID,
	) (Collection[uuid.UUID], error)
}

type IWorkingHoursService interface {
	ListPickUpPointWorkingHours(
		pickUpPointId uuid.UUID,
	) (Collection[models.PickUpPointWorkingHours], error)
}

