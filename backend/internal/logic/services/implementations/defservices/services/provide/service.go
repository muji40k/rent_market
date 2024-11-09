package provide

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/states"
	"rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	instance_providers "rent_service/internal/repository/context/providers/instance"
	provision_providers "rent_service/internal/repository/context/providers/provision"
	role_providers "rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type repoproviders struct {
	renter    role_providers.IRenterProvider
	provision provision_providers.IProvider
	request   provision_providers.IRequestProvider
	revoke    provision_providers.IRevokeProvider
	photo     instance_providers.IPhotoProvider
}

type accessors struct {
	instance access.IInstance
	user     access.IUser
	request  access.IProvisionRequest
	revoke   access.IProvisionRevoke
}

type service struct {
	repos         repoproviders
	access        accessors
	smachine      states.IInstanceStateMachine
	authenticator authenticator.IAuthenticator
	registry      photoregistry.IRegistry
}

func New(
	smachine states.IInstanceStateMachine,
	authenticator authenticator.IAuthenticator,
	registry photoregistry.IRegistry,
	renterProvider role_providers.IRenterProvider,
	provisionProvider provision_providers.IProvider,
	requestProvider provision_providers.IRequestProvider,
	revokeProvider provision_providers.IRevokeProvider,
	photoProvider instance_providers.IPhotoProvider,
	instanceAcc access.IInstance,
	userAcc access.IUser,
	requestAcc access.IProvisionRequest,
	revokeAcc access.IProvisionRevoke,
) provide.IService {
	return &service{
		repoproviders{
			renterProvider,
			provisionProvider,
			requestProvider,
			revokeProvider,
			photoProvider,
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
		provisions, err = repo.GetByRenterUserId(userId)

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
	var provision records.Provision
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

		provision, err = self.smachine.AcceptProvisionRequest(
			request.Id,
			form,
		)
		err = states.MapError(err)
	}

	var ids []uuid.UUID
	if nil == err {
		ids, err = photoregistry.MoveFromTemps(self.registry, form.TempPhotos...)
	}

	if nil == err {
		repo := self.repos.photo.GetInstancePhotoRepository()

		for i := 0; len(ids) > i; i++ {
			cerr := repo.Create(provision.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
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

	var ids []uuid.UUID
	if nil == err {
		ids, err = photoregistry.MoveFromTemps(self.registry, form.TempPhotos...)
	}

	if nil == err {
		repo := self.repos.photo.GetInstancePhotoRepository()

		for i := 0; len(ids) > i; i++ {
			cerr := repo.Create(revoke.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
	}

	return err
}

