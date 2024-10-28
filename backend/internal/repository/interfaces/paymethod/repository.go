package paymethod

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
)

type IRepository interface {
	GetAll() (collection.Collection[models.PayMethod], error)
}

