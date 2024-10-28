package period

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IRepository interface {
	GetById(periodId uuid.UUID) (models.Period, error)
	GetAll() (collection.Collection[models.Period], error)
}

