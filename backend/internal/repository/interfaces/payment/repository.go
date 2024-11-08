package payment

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/payment/repository.go

type IRepository interface {
	GetByInstanceId(instanceId uuid.UUID) (collection.Collection[models.Payment], error)
	GetByRentId(rentId uuid.UUID) (collection.Collection[models.Payment], error)
}

