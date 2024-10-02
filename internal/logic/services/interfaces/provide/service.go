package provide

import (
	"github.com/google/uuid"

	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
)

type IService interface {
	ListProvisionsByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[Provision], error)
	GetProvisionByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (Provision, error)

	StartProvision(token token.Token, form StartForm) error
	RejectProvision(token token.Token, requestId uuid.UUID) error
	StopProvision(token token.Token, form StopForm) error
}

type IRequestService interface {
	ListProvisionRequstsByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[ProvideRequest], error)
	GetProvisionRequestByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (ProvideRequest, error)
	ListProvisionRequstsByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[ProvideRequest], error)

	CreateProvisionRequest(
		token token.Token,
		form RequestCreateForm,
	) (ProvideRequest, error)
}

type IRevokeService interface {
	ListProvisionRevokesByUser(
		token token.Token,
		userId uuid.UUID,
	) (Collection[RevokeRequest], error)
	GetProvisionRevokeByInstance(
		token token.Token,
		instance uuid.UUID,
	) (RevokeRequest, error)
	ListProvisionRetvokesByPickUpPoint(
		token token.Token,
		pickUpPointId uuid.UUID,
	) (Collection[RevokeRequest], error)

	CreateProvisionRevoke(
		token token.Token,
		form RevokeCreateForm,
	) (RevokeRequest, error)
	CancelProvisionRevoke(token token.Token, requestId uuid.UUID) error
}

