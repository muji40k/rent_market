package delivery

import (
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListDeliveriesByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[Delivery], error)
	GetDeliveryByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (Delivery, error)

	CreateDelivery(token token.Token, form CreateForm) (Delivery, error)
	SendDelivery(token token.Token, form SendForm) error
	AcceptDelivery(token token.Token, form AcceptForm) error
}

type ICompanyService interface {
	ListDeliveryCompanies(
		token token.Token,
	) (Collection[DeliveryCompany], error)
	GetDeliveryCompanyById(
		token token.Token,
		companyId uuid.UUID,
	) (DeliveryCompany, error)
}

