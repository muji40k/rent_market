package provide

import (
	"github.com/google/uuid"

	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/misc/types/currency"
)

type StartForm struct {
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

type StopForm struct {
	RevokeId         uuid.UUID
	VerificationCode string
	TempPhotos       Collection[uuid.UUID]
}

type IService interface {
	ListProvisionsByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[records.Provision], error)
	GetProvisionByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (records.Provision, error)

	StartProvision(token models.Token, form StartForm) error
	RejectProvision(token models.Token, requestId uuid.UUID) error
	StopProvision(token models.Token, form StopForm) error
}

type RequestCreateForm struct {
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

type IRequestService interface {
	ListProvisionRequstsByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[requests.Provide], error)
	GetProvisionRequestByInstance(
		token models.Token,
		instanceId uuid.UUID,
	) (requests.Provide, error)
	ListProvisionRequstsByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Provide], error)

	CreateProvisionRequest(
		token models.Token,
		form RequestCreateForm,
	) (requests.Provide, error)
}

type RevokeCreateForm struct {
	ProvisionId   uuid.UUID
	PickUpPointId uuid.UUID
}

type IRevokeService interface {
	ListProvisionRevokesByUser(
		token models.Token,
		userId uuid.UUID,
	) (Collection[requests.Revoke], error)
	GetProvisionRevokeByInstance(
		token models.Token,
		instance uuid.UUID,
	) (requests.Revoke, error)
	ListProvisionRetvokesByPickUpPoint(
		token models.Token,
		pickUpPointId uuid.UUID,
	) (Collection[requests.Revoke], error)

	CreateProvisionRevoke(
		token models.Token,
		form RevokeCreateForm,
	) (requests.Revoke, error)
	CancelProvisionRevoke(token models.Token, requestId uuid.UUID) error
}

