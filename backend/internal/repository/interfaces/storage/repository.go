package storage

import (
	"rent_service/internal/domain/records"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IRepository interface {
	Create(storage records.Storage) (records.Storage, error)

	Update(storage records.Storage) error

	GetActiveByPickUpPointId(pickUpPointId uuid.UUID) (Collection[records.Storage], error)
	GetActiveByInstanceId(instanceId uuid.UUID) (records.Storage, error)
}

