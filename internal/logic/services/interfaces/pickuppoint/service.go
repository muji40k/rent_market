package pickuppoint

import (
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListPickUpPoints() (Collection[PickUpPoint], error)
	GetPickUpPointById(pickUpPointId uuid.UUID) (PickUpPoint, error)
}

type IPhotoService interface {
	ListPickUpPointPhotos(
		pickUpPointId uuid.UUID,
	) (Collection[uuid.UUID], error)
}

type IWorkingHoursService interface {
	ListPickUpPointWorkingHours(
		pickUpPointId uuid.UUID,
	) (Collection[WorkingHours], error)
}

