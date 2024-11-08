package defstates

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	istates "rent_service/internal/logic/services/implementations/defservices/states"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type recordsList struct {
	storage         *records.Storage
	delivery        *requests.Delivery
	rent            *records.Rent
	rentRequest     *requests.Rent
	rentReturn      *requests.Return
	provision       *records.Provision
	provisionRevoke *requests.Revoke
}

var STORAGE_MASK = getMask(
	STATE_STORAGE,
	STATE_RENT_AWAIT_DELIVERY_START, STATE_RENT_AWAIT,
	STATE_REVOKE_AWAIT_DELIVERY_START, STATE_REVOKE_AWAIT,
	STATE_DELIVERY_AWAIT_START,
)
var NO_STORAGE_MASK = getInvertMask(STORAGE_MASK)

func (self *InstanceStateMachine) checkStorage(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.storage.GetStorageRepository()
	storage, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_STORAGE_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= STORAGE_MASK
		recList.storage = new(records.Storage)
		*recList.storage = storage
	}

	return nil
}

var DELIVERY_MASK = getMask(
	STATE_RENT_AWAIT_DELIVERY_START, STATE_RENT_AWAIT_DELIVERY,
	STATE_REVOKE_AWAIT_DELIVERY_START, STATE_REVOKE_AWAIT_DELIVERY,
	STATE_DELIVERY_AWAIT_START, STATE_DELIVERY_AWAIT,
)
var NO_DELIVERY_MASK = getInvertMask(DELIVERY_MASK)

func (self *InstanceStateMachine) checkDelivery(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.delivery.GetDeliveryRepository()
	delivery, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_DELIVERY_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= DELIVERY_MASK
		recList.delivery = new(requests.Delivery)
		*recList.delivery = delivery
	}

	return nil
}

var RENT_MASK = getMask(
	STATE_RENT, STATE_RETURN_FORCED_ISSUE,
	STATE_RETURN_AWAIT, STATE_RETURN_FORCED_AWAIT,
)
var NO_RENT_MASK = getInvertMask(RENT_MASK)

func (self *InstanceStateMachine) checkRent(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.rent.self.GetRentRepository()
	rent, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_RENT_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= RENT_MASK
		recList.rent = new(records.Rent)
		*recList.rent = rent
	}

	return nil
}

var RENT_REQUEST_MASK = getMask(
	STATE_RENT_AWAIT_DELIVERY_START, STATE_RENT_AWAIT_DELIVERY,
	STATE_RENT_AWAIT,
)
var NO_RENT_REQUEST_MASK = getInvertMask(RENT_REQUEST_MASK)

func (self *InstanceStateMachine) checkRentRequest(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.rent.request.GetRentRequestRepository()
	rentRequest, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_RENT_REQUEST_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= RENT_REQUEST_MASK
		recList.rentRequest = new(requests.Rent)
		*recList.rentRequest = rentRequest
	}

	return nil
}

var RENT_RETURN_MASK = getMask(
	STATE_RETURN_AWAIT, STATE_RETURN_FORCED_AWAIT,
)
var NO_RENT_RETURN_MASK = getInvertMask(RENT_RETURN_MASK)

func (self *InstanceStateMachine) checkRentReturn(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.rent.retrn.GetRentReturnRepository()
	rentReturn, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_RENT_RETURN_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= RENT_RETURN_MASK
		recList.rentReturn = new(requests.Return)
		*recList.rentReturn = rentReturn
	}

	return nil
}

var NO_PROVISION_MASK = getMask(STATE_REVOKED)
var PROVISION_MASK = getInvertMask(NO_PROVISION_MASK)

func (self *InstanceStateMachine) checkProvision(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.provision.self.GetProvisionRepository()
	provision, err := repo.GetActiveByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_PROVISION_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= PROVISION_MASK
		recList.provision = new(records.Provision)
		*recList.provision = provision
	}

	return nil
}

var PROVISION_REVOKE_MASK = getMask(
	STATE_RETURN_FORCED_ISSUE, STATE_RETURN_FORCED_AWAIT,
	STATE_REVOKE_AWAIT_DELIVERY_START, STATE_REVOKE_AWAIT_DELIVERY,
	STATE_REVOKE_AWAIT,
)
var NO_PROVISION_REVOKE_MASK = getInvertMask(PROVISION_REVOKE_MASK)

func (self *InstanceStateMachine) checkProvisionRevoke(
	instanceId uuid.UUID,
	possible *uint,
	recList *recordsList,
) error {
	repo := self.repos.provision.revoke.GetRevokeProvisionRepository()
	provisionRevoke, err := repo.GetByInstanceId(instanceId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		*possible &= NO_PROVISION_REVOKE_MASK
	} else if nil != err {
		return cmnerrors.Internal(cmnerrors.DataAccess(err))
	} else {
		*possible &= PROVISION_REVOKE_MASK
		recList.provisionRevoke = new(requests.Revoke)
		*recList.provisionRevoke = provisionRevoke
	}

	return nil
}

// Get all possible information about instance and return it's state.
// If state is correct, actions and transmissions can be applied without
// verification
func (self *InstanceStateMachine) parseInstance(instanceId uuid.UUID) (instance, recordsList, error) {
	var empty instance
	var recList recordsList
	var err error
	var minstance models.Instance
	possible := getAll()

	repo := self.repos.instance.self.GetInstanceRepository()
	minstance, err = repo.GetById(instanceId)
	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		return empty, recList, istates.ErrorUnknownInstance
	} else if nil != err {
		return empty, recList, cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	var selectors = [...]func(uuid.UUID, *uint, *recordsList) error{
		self.checkStorage,
		self.checkDelivery,
		self.checkRent,
		self.checkRentRequest,
		self.checkRentReturn,
		self.checkProvision,
		self.checkProvisionRevoke,
	}

	for _, selector := range selectors {
		if err := selector(instanceId, &possible, &recList); nil != err {
			return empty, recList, err
		}
	}

	states := getStates(possible)

	if 1 != len(states) {
		return empty, recList, istates.ErrorBrokenState
	}

	return instance{instanceId, minstance, states[0]}, recList, nil
}

