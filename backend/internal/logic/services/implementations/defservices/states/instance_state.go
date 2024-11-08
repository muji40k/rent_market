package states

import (
	"errors"
	"fmt"
	"time"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/interfaces/provide"

	"github.com/google/uuid"
)

type IInstanceStateMachine interface {
	CreateProvisionRequest(renterId uuid.UUID, form provide.RequestCreateForm) (requests.Provide, error)
	RejectProvisionRequest(requestId uuid.UUID) error
	AcceptProvisionRequest(requestId uuid.UUID, form provide.StartForm) (records.Provision, error)
	CreateDelivery(instanceId uuid.UUID, fromId uuid.UUID, toId uuid.UUID) (requests.Delivery, error)
	SendDelivery(instanceId uuid.UUID, deliveryId uuid.UUID, verificationCode string) error
	AcceptDelivery(instanceId uuid.UUID, deliveryId uuid.UUID, comment *string, verificationCode string) error
	CreateRentRequest(instanceId uuid.UUID, userId uuid.UUID, pickUpPointId uuid.UUID, paymentPeriodId uuid.UUID) (requests.Rent, error)
	AcceptRentRequest(instanceId uuid.UUID, requestId uuid.UUID, verificationCode string) (records.Rent, error)
	RejectRentRequest(instanceId uuid.UUID, requestId uuid.UUID) error
	CreateRentReturn(instanceId uuid.UUID, rentId uuid.UUID, pickUpPointId uuid.UUID, endDate time.Time) (requests.Return, error)
	AcceptRentReturn(instanceId uuid.UUID, returnId uuid.UUID, comment *string, verificationCode string) (records.Storage, error)
	CancelRentReturn(instanceId uuid.UUID, returnId uuid.UUID) error
	CreateProvisionReturn(instanceId uuid.UUID, provisionId uuid.UUID, pickUpPointId uuid.UUID) (requests.Revoke, error)
	AcceptProvisionReturn(instanceId uuid.UUID, revokeId uuid.UUID, verificationCode string) error
	CancelProvisionReturn(instanceId uuid.UUID, revokeId uuid.UUID) error
}

var ErrorUnknownInstance = errors.New("Unknown instance")
var ErrorBrokenState = errors.New("Broken state")

type ErrorForbiddenMethod struct {
	Err error
}

func (self ErrorForbiddenMethod) Error() string {
	return fmt.Sprintf("Method forbidden: '%v'", self.Err)
}

func (self ErrorForbiddenMethod) Unwrap() error {
	return self.Err
}

