package services

import (
	"time"

	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
)

type IRentSevice interface {
	ListRentsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[records.Rent], error)
	GetRentByInstance(token Token, instanceId uuid.UUID) (records.Rent, error)
}

type RentCreateForm struct {
	InstanceId      uuid.UUID
	PickUpPointId   uuid.UUID
	PaymentPeriodId uuid.UUID
}

type RentStartForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type IRentRequestService interface {
	ListRentRequstsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[requests.Rent], error)
	GetRentRequestByInstance(
		token Token,
		instanceId uuid.UUID,
	) (requests.Rent, error)
	ListRentRequstsByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Rent], error)

	CreateRentRequest(
		token Token,
		form RentCreateForm,
	) (requests.Rent, error)
	StartRent(userId uuid.UUID, form RentStartForm) error
	RejectRent(userId uuid.UUID, requestId uuid.UUID) error
}

type ReturnCreateForm struct {
	RentId        uuid.UUID
	PickUpPointId uuid.UUID
	EndDate       time.Time
}

type RentStopForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	Comment          *string
	TempPhotos       Collection[uuid.UUID]
}

type IReturnRequestService interface {
	ListReturnRequestsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[requests.Return], error)
	GetReturnRequestByInstance(
		token Token,
		instance uuid.UUID,
	) (requests.Return, error)
	ListReturnRequestsByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Return], error)

	CreateReturnRequest(
		token Token,
		form ReturnCreateForm,
	) (requests.Return, error)
	CancelReturnRequest(token Token, requestId uuid.UUID) error
	StopRent(userId uuid.UUID, form RentStopForm) error
}

