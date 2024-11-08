package paymethod

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/paymethod/repository.go

type IRepository interface {
	GetAll() (collection.Collection[models.PayMethod], error)
}

