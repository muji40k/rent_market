package rent

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/states"
	"rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	rent_providers "rent_service/internal/repository/context/providers/rent"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type returnRepoProviders struct {
	rent  rent_providers.IProvider
	retrn rent_providers.IReturnProvider
}

type returnAccessors struct {
	instance    access.IInstance
	user        access.IUser
	pickUpPoint access.IPickUpPoint
	rent        access.IRent
	retrn       access.IRentReturn
}

type returnService struct {
	repos         returnRepoProviders
	access        returnAccessors
	smachine      states.IInstanceStateMachine
	authenticator authenticator.IAuthenticator
}

func NewReturn(
	smachine states.IInstanceStateMachine,
	authenticator authenticator.IAuthenticator,
	rentProvider rent_providers.IProvider,
	retrnProvider rent_providers.IReturnProvider,
	instanceAcc access.IInstance,
	userAcc access.IUser,
	pickUpPointAcc access.IPickUpPoint,
	rentAcc access.IRent,
	retrnAcc access.IRentReturn,
) rent.IReturnService {
	return &returnService{
		returnRepoProviders{rentProvider, retrnProvider},
		returnAccessors{
			instanceAcc,
			userAcc,
			pickUpPointAcc,
			rentAcc,
			retrnAcc,
		},
		smachine,
		authenticator,
	}
}

func mapReturn(value *requests.Return) rent.ReturnRequest {
	return rent.ReturnRequest{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		RentEndDate:      date.New(value.RentEndDate),
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

func (self *returnService) ListRentReturnsByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[rent.ReturnRequest], error) {
	var returns Collection[requests.Return]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.retrn.GetRentReturnRepository()
		returns, err = repo.GetByUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapReturn, returns), err
}

func (self *returnService) GetRentReturnByInstance(
	token token.Token,
	instance uuid.UUID,
) (rent.ReturnRequest, error) {
	var retrn requests.Return
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instance)
	}

	if nil == err {
		repo := self.repos.retrn.GetRentReturnRepository()
		retrn, err = repo.GetByInstanceId(instance)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return mapReturn(&retrn), err
}

func (self *returnService) ListRentReturnsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[rent.ReturnRequest], error) {
	var returns Collection[requests.Return]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.retrn.GetRentReturnRepository()
		returns, err = repo.GetByPickUpPointId(pickUpPointId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapReturn, returns), err
}

func (self *returnService) CreateRentReturn(
	token token.Token,
	form rent.ReturnCreateForm,
) (rent.ReturnRequest, error) {
	var rent records.Rent
	var retrn requests.Return
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.rent.Access(user.Id, form.RentId)
	}

	if nil == err {
		repo := self.repos.rent.GetRentRepository()
		rent, err = repo.GetById(form.RentId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		retrn, err = self.smachine.CreateRentReturn(
			rent.InstanceId,
			rent.Id,
			form.PickUpPointId,
			form.EndDate.Time,
		)
		err = states.MapError(err)
	}

	return mapReturn(&retrn), err
}

func (self *returnService) CancelRentReturn(
	token token.Token,
	requestId uuid.UUID,
) error {
	var retrn requests.Return
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.retrn.Access(user.Id, requestId)
	}

	if nil == err {
		repo := self.repos.retrn.GetRentReturnRepository()
		retrn, err = repo.GetById(requestId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = states.MapError(self.smachine.CancelRentReturn(
			retrn.InstanceId,
			retrn.Id,
		))
	}

	return err
}

