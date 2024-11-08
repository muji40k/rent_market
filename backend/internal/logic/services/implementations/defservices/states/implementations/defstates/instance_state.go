package defstates

import (
	"errors"
	"fmt"
	"reflect"
	"rent_service/internal/domain/models"
	"runtime"
	"time"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	delivery_creator "rent_service/internal/logic/delivery"

	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/codegen"
	"rent_service/internal/logic/services/implementations/defservices/states"
	"rent_service/internal/logic/services/interfaces/provide"
	delivery_provider "rent_service/internal/repository/context/providers/delivery"
	instance_provider "rent_service/internal/repository/context/providers/instance"
	period_provider "rent_service/internal/repository/context/providers/period"
	pickuppoint_provider "rent_service/internal/repository/context/providers/pickuppoint"
	provision_provider "rent_service/internal/repository/context/providers/provision"
	rent_provider "rent_service/internal/repository/context/providers/rent"
	storage_provider "rent_service/internal/repository/context/providers/storage"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"slices"

	"github.com/google/uuid"
)

type instance struct {
	instanceId uuid.UUID
	instance   models.Instance
	state      State
}

type provisionRepos struct {
	self    provision_provider.IProvider
	request provision_provider.IRequestProvider
	revoke  provision_provider.IRevokeProvider
}

type rentRepos struct {
	self    rent_provider.IProvider
	request rent_provider.IRequestProvider
	retrn   rent_provider.IReturnProvider
}

type instanceRepos struct {
	self     instance_provider.IProvider
	payPlans instance_provider.IPayPlansProvider
}

type repos struct {
	rent        rentRepos
	provision   provisionRepos
	instance    instanceRepos
	delivery    delivery_provider.IProvider
	storage     storage_provider.IProvider
	pickUpPoint pickuppoint_provider.IProvider
	period      period_provider.IProvider
}

type InstanceStateMachine struct {
	repos    repos
	delivery delivery_creator.ICreator
	code     codegen.IGenerator
}

func New(
	delivery delivery_creator.ICreator,
	code codegen.IGenerator,
	deliveryProvider delivery_provider.IProvider,
	storageProvider storage_provider.IProvider,
	pickUpPointProvider pickuppoint_provider.IProvider,
	periodProvider period_provider.IProvider,
	rentProvider rent_provider.IProvider,
	rentRequestProvider rent_provider.IRequestProvider,
	rentReturnProvider rent_provider.IReturnProvider,
	provisionProvider provision_provider.IProvider,
	provisionRequestProvider provision_provider.IRequestProvider,
	provisionRevokeProvider provision_provider.IRevokeProvider,
	instanceProvider instance_provider.IProvider,
	instancePayPlansProvider instance_provider.IPayPlansProvider,
) states.IInstanceStateMachine {
	return &InstanceStateMachine{
		repos{
			rentRepos{
				rentProvider,
				rentRequestProvider,
				rentReturnProvider,
			},
			provisionRepos{
				provisionProvider,
				provisionRequestProvider,
				provisionRevokeProvider,
			},
			instanceRepos{
				instanceProvider,
				instancePayPlansProvider,
			},
			deliveryProvider,
			storageProvider,
			pickUpPointProvider,
			periodProvider,
		},
		delivery,
		code,
	}
}

func (self *InstanceStateMachine) CreateProvisionRequest(
	renterId uuid.UUID,
	form provide.RequestCreateForm,
) (requests.Provide, error) {
	return self.actionCreateProvisionRequest(renterId, form)
}

func (self *InstanceStateMachine) RejectProvisionRequest(
	requestId uuid.UUID,
) error {
	// Validate
	repo := self.repos.provision.request.GetProvisionRequestRepository()
	request, err := repo.GetById(requestId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	// Act
	if nil == err {
		err = self.actionRejectProvisionRequest(request)
	}

	// Transmit
	// Empty

	return err
}

func (self *InstanceStateMachine) AcceptProvisionRequest(
	requestId uuid.UUID,
	form provide.StartForm,
) (records.Provision, error) {
	var provision records.Provision
	// Validate
	repo := self.repos.provision.request.GetProvisionRequestRepository()
	request, err := repo.GetById(requestId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	// Act
	if nil == err {
		err = self.actionAcceptProvisionRequest(request)
	}

	// Transmit
	var instance models.Instance

	if nil == err {
		instance, err = self.actionCreateInstance(request, form)
	}

	if nil == err {
		provision, err = self.actionCreateProvision(instance, request)
	}

	if nil == err {
		_, err = self.actionCreateStorage(instance, request.PickUpPointId)
	}

	return provision, err
}

func (self *InstanceStateMachine) CreateDelivery(
	instanceId uuid.UUID,
	fromId uuid.UUID,
	toId uuid.UUID,
) (requests.Delivery, error) {
	var request requests.Delivery
	var allowed_states = [...]State{STATE_STORAGE}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CreateDelivery)
	}

	if nil == err && fromId != rl.storage.PickUpPointId {
		err = cmnerrors.Incorrect("from_id")
	}

	// Act
	if nil == err {
		request, err = self.actionCreateDelivery(
			instance.instance,
			*rl.storage,
			toId,
		)
	}

	// Transmit
	// Empty

	return request, err
}

func (self *InstanceStateMachine) SendDelivery(
	instanceId uuid.UUID,
	deliveryId uuid.UUID,
	verificationCode string,
) error {
	var allowed_states = [...]State{
		STATE_DELIVERY_AWAIT_START, STATE_RENT_AWAIT_DELIVERY_START,
		STATE_REVOKE_AWAIT_DELIVERY_START,
	}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.SendDelivery)
	}

	if nil == err && deliveryId != rl.delivery.Id {
		err = cmnerrors.Incorrect("delivery_id")
	}

	// Act
	if nil == err {
		err = self.actionSendDelivery(*rl.delivery, verificationCode)
	}

	// Transmit
	if nil == err {
		err = self.actionStopStorage(*rl.storage)
	}

	return err
}

func (self *InstanceStateMachine) AcceptDelivery(
	instanceId uuid.UUID,
	deliveryId uuid.UUID,
	comment *string,
	verificationCode string,
) error {
	var allowed_states = [...]State{
		STATE_DELIVERY_AWAIT, STATE_RENT_AWAIT_DELIVERY,
		STATE_REVOKE_AWAIT_DELIVERY,
	}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.AcceptDelivery)
	}

	if nil == err && deliveryId != rl.delivery.Id {
		err = cmnerrors.Incorrect("delivery_id")
	}

	// Act
	if nil == err {
		err = self.actionAcceptDelivery(*rl.delivery, verificationCode)
	}

	// Transmit
	if nil == err {
		_, err = self.actionCreateStorage(instance.instance, rl.delivery.ToId)
	}

	if nil == err {
		err = self.actionUpdateInstance(
			instance.instance,
			instanceUpdateDescription(comment),
		)
	}

	return err
}

func (self *InstanceStateMachine) CreateRentRequest(
	instanceId uuid.UUID,
	userId uuid.UUID,
	pickUpPointId uuid.UUID,
	paymentPeriodId uuid.UUID,
) (requests.Rent, error) {
	var request requests.Rent
	var allowed_states = [...]State{STATE_STORAGE}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CreateRentRequest)
	}

	// Act
	if nil == err {
		request, err = self.actionCreateRentRequest(
			instance.instance,
			userId,
			pickUpPointId,
			paymentPeriodId,
		)
	}

	// Transmit
	if nil == err && pickUpPointId != rl.storage.PickUpPointId {
		_, err = self.actionCreateDelivery(
			instance.instance,
			*rl.storage,
			pickUpPointId,
		)

	}

	return request, err
}

func (self *InstanceStateMachine) AcceptRentRequest(
	instanceId uuid.UUID,
	requestId uuid.UUID,
	verificationCode string,
) (records.Rent, error) {
	var rent records.Rent
	var allowed_states = [...]State{STATE_RENT_AWAIT}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.AcceptRentRequest)
	}

	if nil == err && requestId != rl.rentRequest.Id {
		err = cmnerrors.Incorrect("rent_request_id")
	}

	// Act
	if nil == err {
		err = self.actionAcceptRentRequest(*rl.rentRequest, verificationCode)
	}

	// Transmit
	if nil == err {
		rent, err = self.actionCreateRent(*rl.rentRequest)
	}

	if nil == err {
		err = self.actionStopStorage(*rl.storage)
	}

	return rent, err
}

func (self *InstanceStateMachine) RejectRentRequest(
	instanceId uuid.UUID,
	requestId uuid.UUID,
) error {
	var allowed_states = [...]State{STATE_RENT_AWAIT}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.RejectRentRequest)
	}

	if nil == err && requestId != rl.rentRequest.Id {
		err = cmnerrors.Incorrect("rent_request_id")
	}

	// Act
	if nil == err {
		err = self.actionRejectRentRequest(*rl.rentRequest)
	}

	// Transmit
	// Empty

	return err
}

func (self *InstanceStateMachine) CreateRentReturn(
	instanceId uuid.UUID,
	rentId uuid.UUID,
	pickUpPointId uuid.UUID,
	endDate time.Time,
) (requests.Return, error) {
	var request requests.Return
	var allowed_states = [...]State{STATE_RENT, STATE_RETURN_FORCED_ISSUE}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CreateRentReturn)
	}

	if nil == err && rentId != rl.rent.Id {
		err = cmnerrors.Incorrect("rent_id")
	}

	// Act
	if nil == err {
		request, err = self.actionCreateRentReturn(
			instance.instance,
			*rl.rent,
			pickUpPointId,
			endDate,
		)
	}

	// Transmit
	// Empty

	return request, err
}

func (self *InstanceStateMachine) AcceptRentReturn(
	instanceId uuid.UUID,
	returnId uuid.UUID,
	comment *string,
	verificationCode string,
) (records.Storage, error) {
	var storage records.Storage
	var allowed_states = [...]State{
		STATE_RETURN_AWAIT, STATE_RETURN_FORCED_AWAIT,
	}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.AcceptRentReturn)
	}

	if nil == err && returnId != rl.rentReturn.Id {
		err = cmnerrors.Incorrect("rent_return_id")
	}

	// Act
	if nil == err {
		err = self.actionAcceptRentReturn(*rl.rentReturn, verificationCode)
	}

	// Transmit
	if nil == err {
		err = self.actionStopRent(*rl.rent)
	}

	if nil == err {
		storage, err = self.actionCreateStorage(
			instance.instance,
			rl.rentReturn.PickUpPointId,
		)
	}

	if nil == err && STATE_RETURN_FORCED_AWAIT == instance.state &&
		storage.PickUpPointId != rl.provisionRevoke.PickUpPointId {
		_, err = self.actionCreateDelivery(
			instance.instance,
			storage,
			rl.provisionRevoke.PickUpPointId,
		)
	}

	if nil == err {
		err = self.actionUpdateInstance(
			instance.instance,
			instanceUpdateDescription(comment),
		)
	}

	return storage, err
}

func (self *InstanceStateMachine) CancelRentReturn(
	instanceId uuid.UUID,
	returnId uuid.UUID,
) error {
	var allowed_states = [...]State{STATE_RETURN_AWAIT}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CancelRentReturn)
	}

	if nil == err && returnId != rl.rentReturn.Id {
		err = cmnerrors.Incorrect("rent_return_id")
	}

	// Act
	if nil == err {
		err = self.actionCancelRentReturn(*rl.rentReturn)
	}

	// Transmit
	// Empty

	return err
}

func (self *InstanceStateMachine) CreateProvisionReturn(
	instanceId uuid.UUID,
	provisionId uuid.UUID,
	pickUpPointId uuid.UUID,
) (requests.Revoke, error) {
	var revoke requests.Revoke
	var allowed_states = [...]State{
		STATE_STORAGE, STATE_RENT, STATE_RETURN_AWAIT,
	}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CreateProvisionReturn)
	}

	if nil == err && provisionId != rl.provision.Id {
		err = cmnerrors.Incorrect("provision_id")
	}

	// Act
	if nil == err {
		revoke, err = self.actionCreateProvisionReturn(
			instance.instance,
			*rl.provision,
			pickUpPointId,
		)
	}

	// Transmit
	if nil == err && STATE_STORAGE == instance.state &&
		pickUpPointId != rl.storage.PickUpPointId {
		_, err = self.actionCreateDelivery(
			instance.instance,
			*rl.storage,
			pickUpPointId,
		)
	}

	return revoke, err
}

func (self *InstanceStateMachine) AcceptProvisionReturn(
	instanceId uuid.UUID,
	revokeId uuid.UUID,
	verificationCode string,
) error {
	var allowed_states = [...]State{STATE_REVOKE_AWAIT}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.AcceptProvisionReturn)
	}

	if nil == err && revokeId != rl.provisionRevoke.Id {
		err = cmnerrors.Incorrect("provision_revoke_id")
	}

	// Act
	if nil == err {
		err = self.actionAcceptProvisionReturn(
			*rl.provisionRevoke,
			verificationCode,
		)
	}

	// Transmit
	if nil == err {
		err = self.actionStopStorage(*rl.storage)
	}

	if nil == err {
		err = self.actionStopProvision(*rl.provision)
	}

	return err
}

func (self *InstanceStateMachine) CancelProvisionReturn(
	instanceId uuid.UUID,
	revokeId uuid.UUID,
) error {
	var allowed_states = [...]State{
		STATE_RETURN_FORCED_ISSUE, STATE_RETURN_FORCED_AWAIT,
	}

	// Validate
	instance, rl, err := self.parseInstance(instanceId)
	if nil == err && !slices.Contains(allowed_states[:], instance.state) {
		err = ForbiddenMethod(instance.state, self.CancelProvisionReturn)
	}

	if nil == err && revokeId != rl.provisionRevoke.Id {
		err = cmnerrors.Incorrect("provision_revoke_id")
	}

	// Act
	if nil == err {
		err = self.actionCancelProvisionReturn(*rl.provisionRevoke)
	}

	// Transmit
	// Empty

	return err
}

func ForbiddenMethod(state State, method interface{}) error {
	// .Pointer() panics but it's intended, so dont't pass anything besides
	// function pointer...
	return states.ErrorForbiddenMethod{
		Err: fmt.Errorf(
			"Instance in state '%v' doesn't allow method '%v'",
			state, runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(),
		),
	}
}

