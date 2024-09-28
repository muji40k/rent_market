package period

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
)

type IService interface {
	GetPeriods() (Collection[models.Period], error)
}

