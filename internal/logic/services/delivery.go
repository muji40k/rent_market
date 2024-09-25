package services

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type DeliveryCreateForm struct {
	InstanceId uuid.UUID
	From       uuid.UUID
	To         uuid.UUID
}

type DeliverySendForm struct {
	DeliveryId       uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type DeliveryAcceptForm struct {
	DeliveryId       uuid.UUID
	Comment          *string
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type IDeliveryService interface {
	ListDeliveriesByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Delivery], error)
	GetDeliveryByInstance(
		token Token,
		instanceId uuid.UUID,
	) (requests.Delivery, error)

	CreateDelivery(
		token Token,
		form DeliveryCreateForm,
	) (requests.Delivery, error)
	SendDelivery(token Token, form DeliverySendForm) error
	AcceptDelivery(token Token, form DeliveryAcceptForm) error
}

type IDeliveryCompanyService interface {
	ListDeliveryCompanies(
		token Token,
	) (Collection[models.DeliveryCompany], error)
	GetDeliveryCompanyById(
		token Token,
		companyId uuid.UUID,
	) (models.DeliveryCompany, error)
}

