package services

import (
	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/misc/types/currency"
)

type IProvisionSevice interface {
	ListProvisionsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[records.Provision], error)
	GetProvisionByInstance(
		token Token,
		instanceId uuid.UUID,
	) (records.Provision, error)
}

type ProvisionCreateForm struct {
	ProductId     uuid.UUID
	PickUpPointId uuid.UUID
	Name          string
	Description   string
	Condition     string
	PayPlans      map[uuid.UUID]struct {
		PeriodId uuid.UUID // Index
		Price    currency.Currency
	}
}

type ProvisionStartForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	Overrides        struct {
		ProductId   *uuid.UUID
		Name        *string
		Description *string
		Condition   *string
		PayPlans    map[uuid.UUID]struct {
			PeriodId uuid.UUID // Index
			Price    currency.Currency
		}
	}
	TempPhotos Collection[uuid.UUID]
}

type IProvisionRequestService interface {
	ListProvisionRequstsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[requests.Provide], error)
	GetProvisionRequestByInstance(
		token Token,
		instanceId uuid.UUID,
	) (requests.Provide, error)
	ListProvisionRequstsByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Provide], error)

	CreateProvisionRequest(
		token Token,
		form ProvisionCreateForm,
	) (requests.Provide, error)
	StartProvision(userId uuid.UUID, form ProvisionStartForm) error
	RejectProvision(userId uuid.UUID, requestId uuid.UUID) error
}

type RevokeCreateForm struct {
	ProvisionId   uuid.UUID
	PickUpPointId uuid.UUID
}

type ProvisionStopForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type IRevokeRequestService interface {
	ListRevokeRequestsByUser(
		token Token,
		userId uuid.UUID,
	) (Collection[requests.Revoke], error)
	GetRevokesRequestByInstance(
		token Token,
		instance uuid.UUID,
	) (requests.Revoke, error)
	ListRetvokeRequestsByPickUpPoint(
		token Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Revoke], error)

	CreateRevokeRequest(
		token Token,
		form RevokeCreateForm,
	) (requests.Revoke, error)
	CancelRevokeRequest(token Token, requestId uuid.UUID) error
	StopProvision(userId uuid.UUID, form ProvisionStopForm) error
}

