package rent

import (
	"github.com/google/uuid"

	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
)

type IService interface {
	ListRentsByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[Rent], error)
	GetRentByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (Rent, error)

	StartRent(token token.Token, form StartForm) error
	RejectRent(token token.Token, requestId uuid.UUID) error
	StopRent(token token.Token, form StopForm) error
}

type IRequestService interface {
	ListRentRequstsByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[RentRequest], error)
	GetRentRequestByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (RentRequest, error)
	ListRentRequstsByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[RentRequest], error)

	CreateRentRequest(
		token token.Token,
		form RequestCreateForm,
	) (RentRequest, error)
}

type IReturnService interface {
	ListRentReturnsByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[ReturnRequest], error)
	GetRentReturnByInstance(
		token token.Token,
		instance uuid.UUID,
	) (ReturnRequest, error)
	ListRentReturnsByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[ReturnRequest], error)

	CreateRentReturn(
		token token.Token,
		form ReturnCreateForm,
	) (ReturnRequest, error)
	CancelRentReturn(token token.Token, requestId uuid.UUID) error
}

