package delivery

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/delivery/repository.go

type IRepository interface {
	Create(delivery requests.Delivery) (requests.Delivery, error)

	Update(delivery requests.Delivery) error

	GetById(deliveryId uuid.UUID) (requests.Delivery, error)
	GetActiveByPickUpPointId(
		pickUpPointId uuid.UUID,
	) (collection.Collection[requests.Delivery], error)
	GetActiveByInstanceId(instanceId uuid.UUID) (requests.Delivery, error)
}

type ICompanyRepository interface {
	GetById(companyId uuid.UUID) (models.DeliveryCompany, error)
	GetAll() (collection.Collection[models.DeliveryCompany], error)
}

