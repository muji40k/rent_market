package storage

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListStoragesByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[records.Storage], error)
	GetStorageByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (records.Storage, error)
}

