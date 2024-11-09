package rent

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
	"rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	instance_providers "rent_service/internal/repository/context/providers/instance"
	rent_providers "rent_service/internal/repository/context/providers/rent"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type repoproviders struct {
	rent    rent_providers.IProvider
	request rent_providers.IRequestProvider
	ret     rent_providers.IReturnProvider
	photo   instance_providers.IPhotoProvider
}

type accessors struct {
	instance access.IInstance
	user     access.IUser
	request  access.IRentRequest
	ret      access.IRentReturn
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
	rentProvider rent_providers.IProvider,
	requestProvider rent_providers.IRequestProvider,
	retProvider rent_providers.IReturnProvider,
	photoProvider instance_providers.IPhotoProvider,
	instanceAcc access.IInstance,
	userAcc access.IUser,
	requestAcc access.IRentRequest,
	retAcc access.IRentReturn,
) rent.IService {
	return &service{
		repoproviders{rentProvider, requestProvider, retProvider, photoProvider},
		accessors{instanceAcc, userAcc, requestAcc, retAcc},
		smachine,
		authenticator,
		registry,
	}
}

func mapf(value *records.Rent) rent.Rent {
	out := rent.Rent{
		Id:              value.Id,
		UserId:          value.UserId,
		InstanceId:      value.InstanceId,
		StartDate:       date.New(value.StartDate),
		PaymentPeriodId: value.PaymentPeriodId,
	}

	if nil != value.EndDate {
		out.EndDate = new(date.Date)
		*out.EndDate = date.New(*value.EndDate)
	}

	return out
}

func (self *service) ListRentsByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[rent.Rent], error) {
	var rents Collection[records.Rent]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.rent.GetRentRepository()
		rents, err = repo.GetByUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapf, rents), err
}

func (self *service) GetRentByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (rent.Rent, error) {
	var rent records.Rent
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		repo := self.repos.rent.GetRentRepository()
		rent, err = repo.GetActiveByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return mapf(&rent), err
}

func (self *service) StartRent(token token.Token, form rent.StartForm) error {
	var request requests.Rent
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.request.Access(user.Id, form.RequestId)
	}

	if nil == err {
		repo := self.repos.request.GetRentRequestRepository()
		request, err = repo.GetById(form.RequestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		_, err = self.smachine.AcceptRentRequest(
			request.InstanceId,
			request.Id,
			form.VerificationCode,
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
			cerr := repo.Create(request.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
	}

	return err
}

func (self *service) RejectRent(token token.Token, requestId uuid.UUID) error {
	var request requests.Rent
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.request.Access(user.Id, requestId)
	}

	if nil == err {
		repo := self.repos.request.GetRentRequestRepository()
		request, err = repo.GetById(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = states.MapError(self.smachine.RejectRentRequest(
			request.InstanceId,
			request.Id,
		))
	}

	return err
}

func (self *service) StopRent(token token.Token, form rent.StopForm) error {
	var request requests.Return
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.ret.Access(user.Id, form.ReturnId)
	}

	if nil == err {
		repo := self.repos.ret.GetRentReturnRepository()
		request, err = repo.GetById(form.ReturnId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		_, err = self.smachine.AcceptRentReturn(
			request.InstanceId,
			request.Id,
			form.Comment,
			form.VerificationCode,
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
			cerr := repo.Create(request.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
	}

	return err
}

