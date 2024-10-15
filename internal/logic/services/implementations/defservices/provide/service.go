package provide

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/access"
	"rent_service/internal/logic/services/implementations/defservices/misc/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/misc/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/misc/states"
	"rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	provision_providers "rent_service/internal/repository/context/providers/provision"
	role_providers "rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type repoproviders struct {
	renter    role_providers.IRenterProvider
	provision provision_providers.IProvider
	request   provision_providers.IRequestProvider
	revoke    provision_providers.IRevokeProvider
}

type accessors struct {
	instance *access.Instance
	user     *access.User
	request  *access.ProvisionRequest
	revoke   *access.ProvisionRevoke
}

type service struct {
	repos         repoproviders
	access        accessors
	smachine      *states.InstanceStateMachine
	authenticator *authenticator.Authenticator
	registry      *photoregistry.Registry
}

func New(
	smachine *states.InstanceStateMachine,
	authenticator *authenticator.Authenticator,
	registry *photoregistry.Registry,
	renterProvider role_providers.IRenterProvider,
	provisionProvider provision_providers.IProvider,
	requestProvider provision_providers.IRequestProvider,
	revokeProvider provision_providers.IRevokeProvider,
	instanceAcc *access.Instance,
	userAcc *access.User,
	requestAcc *access.ProvisionRequest,
	revokeAcc *access.ProvisionRevoke,
) provide.IService {
	return &service{
		repoproviders{
			renterProvider,
			provisionProvider,
			requestProvider,
			revokeProvider,
		},
		accessors{instanceAcc, userAcc, requestAcc, revokeAcc},
		smachine,
		authenticator,
		registry,
	}
}

func (self *service) mapf(value *records.Provision) provide.Provision {
	var userId uuid.UUID
	repo := self.repos.renter.GetRenterRepository()

	renter, err := repo.GetById(value.RenterId)

	if nil == err {
		userId = renter.UserId
	}

	out := provide.Provision{
		Id:         value.Id,
		UserId:     userId,
		InstanceId: value.InstanceId,
		StartDate:  date.New(value.StartDate),
	}

	if nil != value.EndDate {
		out.EndDate = new(date.Date)
		*out.EndDate = date.New(*value.EndDate)
	}

	return out
}

func (self *service) ListProvisionsByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[provide.Provision], error) {
	var provisions Collection[records.Provision]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.provision.GetProvisionRepository()
		provisions, err = repo.GetActiveByRenterUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(self.mapf, provisions), err
}

func (self *service) GetProvisionByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (provide.Provision, error) {
	var provision records.Provision
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		repo := self.repos.provision.GetProvisionRepository()
		provision, err = repo.GetActiveByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return self.mapf(&provision), err
}

func clearOverrides(overrides *provide.Overrides) {
	if nil != overrides.Name && "" == *overrides.Name {
		overrides.Name = nil
	}

	if nil != overrides.Description && "" == *overrides.Description {
		overrides.Description = nil
	}

	if nil != overrides.Condition && "" == *overrides.Condition {
		overrides.Condition = nil
	}

	if nil != overrides.PayPlans && 0 == len(overrides.PayPlans) {
		overrides.PayPlans = nil
	}
}

func (self *service) StartProvision(
	token token.Token,
	form provide.StartForm,
) error {
	var request requests.Provide
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.request.Access(user.Id, form.RequestId)
	}

	if nil == err {
		repo := self.repos.request.GetProvisionRequestRepository()
		request, err = repo.GetById(form.RequestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		clearOverrides(&form.Overrides)

		err = states.MapError(self.smachine.AcceptProvisionRequest(
			request.Id,
			form,
		))
	}

	if nil == err {
		_, err = self.registry.MoveFromTemps(form.TempPhotos...)
	}

	return err
}

func (self *service) RejectProvision(
	token token.Token,
	requestId uuid.UUID,
) error {
	var request requests.Provide
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.request.Access(user.Id, requestId)
	}

	if nil == err {
		repo := self.repos.request.GetProvisionRequestRepository()
		request, err = repo.GetById(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = states.MapError(self.smachine.RejectProvisionRequest(request.Id))
	}

	return err
}

func (self *service) StopProvision(
	token token.Token,
	form provide.StopForm,
) error {
	var revoke requests.Revoke
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.revoke.Access(user.Id, form.RevokeId)
	}

	if nil == err {
		repo := self.repos.revoke.GetRevokeProvisionRepository()
		revoke, err = repo.GetById(form.RevokeId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = states.MapError(self.smachine.AcceptProvisionReturn(
			revoke.InstanceId,
			revoke.Id,
			form.VerificationCode,
		))
	}

	if nil == err {
		_, err = self.registry.MoveFromTemps(form.TempPhotos...)
	}

	return err
}

