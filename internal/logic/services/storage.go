package services

import (
	"rent_service/internal/domain/records"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IStorageService interface {
	ListStoragesByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[records.Storage], error)
	GetStorageByInstance(
		token Token,
		instanceId uuid.UUID,
	) (records.Storage, error)
}

