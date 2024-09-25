package services

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
)

type IPeriodService interface {
	GetPeriods() (Collection[models.Period], error)
}

