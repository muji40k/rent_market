package storage

import (
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListStoragesByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[Storage], error)
	GetStorageByInstance(instanceId uuid.UUID) (Storage, error)
}

