package payment

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IRepository interface {
	GetByInstanceId(instanceId uuid.UUID) (Collection[models.Payment], error)
	GetByRentId(rentId uuid.UUID) (Collection[models.Payment], error)
}

