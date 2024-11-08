package rent

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/rent/repository.go

type IRepository interface {
	Create(rent records.Rent) (records.Rent, error)

	Update(rent records.Rent) error

	GetById(rentId uuid.UUID) (records.Rent, error)
	GetByUserId(userId uuid.UUID) (collection.Collection[records.Rent], error)
	GetActiveByInstanceId(instanceId uuid.UUID) (records.Rent, error)
	GetPastByUserId(userId uuid.UUID) (collection.Collection[records.Rent], error)
}

type IRequestRepository interface {
	Create(request requests.Rent) (requests.Rent, error)

	GetById(requestId uuid.UUID) (requests.Rent, error)
	GetByUserId(userId uuid.UUID) (collection.Collection[requests.Rent], error)
	GetByInstanceId(instanceId uuid.UUID) (requests.Rent, error)
	GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Rent], error)

	Remove(requestId uuid.UUID) error
}

type IReturnRepository interface {
	Create(request requests.Return) (requests.Return, error)

	GetById(requestId uuid.UUID) (requests.Return, error)
	GetByUserId(userId uuid.UUID) (collection.Collection[requests.Return], error)
	GetByInstanceId(instanceId uuid.UUID) (requests.Return, error)
	GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Return], error)

	Remove(requestId uuid.UUID) error
}

