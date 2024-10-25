package provide

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/access"
	"rent_service/internal/logic/services/implementations/defservices/misc/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/misc/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/misc/states"
	"rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	provision_providers "rent_service/internal/repository/context/providers/provision"
	rent_providers "rent_service/internal/repository/context/providers/rent"
	role_providers "rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type revokeRepoProviders struct {
	renter    role_providers.IRenterProvider
	provision provision_providers.IProvider
	revoke    provision_providers.IRevokeProvider
	rent      rent_providers.IProvider
}

type revokeAccessors struct {
	instance    *access.Instance
	user        *access.User
	pickUpPoint *access.PickUpPoint
	provision   *access.Provision
	revoke      *access.ProvisionRevoke
}

type revokeService struct {
	repos         revokeRepoProviders
	access        revokeAccessors
	smachine      *states.InstanceStateMachine
	authenticator *authenticator.Authenticator
	authorizer    *authorizer.Authorizer
}

func NewRevoke(
	smachine *states.InstanceStateMachine,
	authenticator *authenticator.Authenticator,
	authorizer *authorizer.Authorizer,
	renterProvider role_providers.IRenterProvider,
	provisionProvider provision_providers.IProvider,
	revokeProvider provision_providers.IRevokeProvider,
	rentProvider rent_providers.IProvider,
	instanceAcc *access.Instance,
	userAcc *access.User,
	pickUpPointAcc *access.PickUpPoint,
	provisionAcc *access.Provision,
	revokeAcc *access.ProvisionRevoke,
) provide.IRevokeService {
	return &revokeService{
		revokeRepoProviders{
			renterProvider,
			provisionProvider,
			revokeProvider,
			rentProvider,
		},
		revokeAccessors{
			instanceAcc,
			userAcc,
			pickUpPointAcc,
			provisionAcc,
			revokeAcc,
		},
		smachine,
		authenticator,
		authorizer,
	}
}

func (self *revokeService) mapRevoke(value *requests.Revoke) provide.RevokeRequest {
	var userId uuid.UUID
	repo := self.repos.renter.GetRenterRepository()

	renter, err := repo.GetById(value.RenterId)

	if nil == err {
		userId = renter.UserId
	}

	return provide.RevokeRequest{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           userId,
		PickUpPointId:    value.PickUpPointId,
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

func (self *revokeService) ListProvisionRevokesByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[provide.RevokeRequest], error) {
	var revokes Collection[requests.Revoke]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.revoke.GetRevokeProvisionRepository()
		revokes, err = repo.GetByUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(self.mapRevoke, revokes), err
}

func (self *revokeService) GetProvisionRevokeByInstance(
	token token.Token,
	instance uuid.UUID,
) (provide.RevokeRequest, error) {
	var revoke requests.Revoke
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instance)
	}

	if aerr := (cmnerrors.ErrorNoAccess{}); errors.As(err, &aerr) {
		repo := self.repos.rent.GetRentRepository()
		rent, ierr := repo.GetActiveByInstanceId(instance)

		if nil == ierr && rent.UserId == user.Id {
			err = nil
		} else if cerr := (repo_errors.ErrorNotFound{}); !errors.As(ierr, &cerr) {
			err = cmnerrors.Internal(cmnerrors.DataAccess(ierr))
		}
	}

	if nil == err {
		repo := self.repos.revoke.GetRevokeProvisionRepository()
		revoke, err = repo.GetByInstanceId(instance)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return self.mapRevoke(&revoke), err
}

func (self *revokeService) ListProvisionRetvokesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[provide.RevokeRequest], error) {
	var revokes Collection[requests.Revoke]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.revoke.GetRevokeProvisionRepository()
		revokes, err = repo.GetByPickUpPointId(pickUpPointId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(self.mapRevoke, revokes), err
}

func (self *revokeService) CreateProvisionRevoke(
	token token.Token,
	form provide.RevokeCreateForm,
) (provide.RevokeRequest, error) {
	var provision records.Provision
	var revoke requests.Revoke
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.provision.Access(user.Id, form.ProvisionId)
	}

	if nil == err {
		repo := self.repos.provision.GetProvisionRepository()
		provision, err = repo.GetById(form.ProvisionId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		revoke, err = self.smachine.CreateProvisionReturn(
			provision.InstanceId,
			provision.Id,
			form.PickUpPointId,
		)
		err = states.MapError(err)
	}

	return self.mapRevoke(&revoke), err
}

func (self *revokeService) CancelProvisionRevoke(
	token token.Token,
	requestId uuid.UUID,
) error {
	var revoke requests.Revoke
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.revoke.Access(user.Id, requestId)
	}

	if nil == err {
		repo := self.repos.revoke.GetRevokeProvisionRepository()
		revoke, err = repo.GetById(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = states.MapError(self.smachine.CancelProvisionReturn(
			revoke.InstanceId,
			revoke.Id,
		))
	}

	return err
}

