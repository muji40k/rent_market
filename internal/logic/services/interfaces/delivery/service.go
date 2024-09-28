package delivery

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type CreateForm struct {
	InstanceId uuid.UUID
	From       uuid.UUID
	To         uuid.UUID
}

type SendForm struct {
	DeliveryId       uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type AcceptForm struct {
	DeliveryId       uuid.UUID
	Comment          *string
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type IService interface {
	ListDeliveriesByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Delivery], error)
	GetDeliveryByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (requests.Delivery, error)

	CreateDelivery(
		token models.Token,
		form CreateForm,
	) (requests.Delivery, error)
	SendDelivery(token models.Token, form SendForm) error
	AcceptDelivery(token models.Token, form AcceptForm) error
}

type ICompanyService interface {
	ListDeliveryCompanies(
		token models.Token,
	) (Collection[models.DeliveryCompany], error)
	GetDeliveryCompanyById(
		token models.Token,
		companyId uuid.UUID,
	) (models.DeliveryCompany, error)
}

