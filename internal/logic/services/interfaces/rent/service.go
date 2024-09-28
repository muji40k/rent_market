package rent

import (
	"time"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
)

type StartForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type StopForm struct {
	ReturnId         uuid.UUID
	VerificationCode string
	Comment          *string
	TempPhotos       Collection[uuid.UUID]
}

type IService interface {
	ListRentsByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[records.Rent], error)
	GetRentByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (records.Rent, error)

	StartRent(userId uuid.UUID, form StartForm) error
	RejectRent(userId uuid.UUID, requestId uuid.UUID) error
	StopRent(userId uuid.UUID, form StopForm) error
}

type RequestCreateForm struct {
	InstanceId      uuid.UUID
	PickUpPointId   uuid.UUID
	PaymentPeriodId uuid.UUID
}

type IRequestService interface {
	ListRentRequstsByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[requests.Rent], error)
	GetRentRequestByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (requests.Rent, error)
	ListRentRequstsByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Rent], error)

	CreateRentRequest(
		token models.Token,
		form RequestCreateForm,
	) (requests.Rent, error)
}

type ReturnCreateForm struct {
	RentId        uuid.UUID
	PickUpPointId uuid.UUID
	EndDate       time.Time
}

type IReturnService interface {
	ListRentReturnsByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[requests.Return], error)
	GetRentReturnByInstance(
		token models.Token,
		instance uuid.UUID,
	) (requests.Return, error)
	ListRentReturnsByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Return], error)

	CreateRentReturn(
		token models.Token,
		form ReturnCreateForm,
	) (requests.Return, error)
	CancelRentReturn(token models.Token, requestId uuid.UUID) error
}

