package provision

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../implementation/mock/provision/repository.go

type IRepository interface {
	Create(provision records.Provision) (records.Provision, error)

	Update(provision records.Provision) error

	GetById(provisionId uuid.UUID) (records.Provision, error)
	GetByRenterUserId(userId uuid.UUID) (collection.Collection[records.Provision], error)
	GetByInstanceId(instanceId uuid.UUID) (collection.Collection[records.Provision], error)
	GetActiveByInstanceId(instanceId uuid.UUID) (records.Provision, error)
}

type IRequestRepository interface {
	Create(request requests.Provide) (requests.Provide, error)

	GetById(requestId uuid.UUID) (requests.Provide, error)
	GetByUserId(userId uuid.UUID) (collection.Collection[requests.Provide], error)
	GetByInstanceId(instanceId uuid.UUID) (requests.Provide, error)
	GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Provide], error)

	Remove(requestId uuid.UUID) error
}

type IRevokeRepository interface {
	Create(request requests.Revoke) (requests.Revoke, error)

	GetById(requestId uuid.UUID) (requests.Revoke, error)
	GetByUserId(userId uuid.UUID) (collection.Collection[requests.Revoke], error)
	GetByInstanceId(instanceId uuid.UUID) (requests.Revoke, error)
	GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Revoke], error)

	Remove(requestId uuid.UUID) error
}

