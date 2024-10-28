package rent

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/access"
	"rent_service/internal/logic/services/implementations/defservices/misc/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/misc/states"
	"rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	rent_providers "rent_service/internal/repository/context/providers/rent"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type requestRepoProviders struct {
	request rent_providers.IRequestProvider
}

type requestAccessors struct {
	instance    *access.Instance
	user        *access.User
	pickUpPoint *access.PickUpPoint
}

type requestService struct {
	repos         requestRepoProviders
	access        requestAccessors
	smachine      *states.InstanceStateMachine
	authenticator *authenticator.Authenticator
}

func NewRequest(
	smachine *states.InstanceStateMachine,
	authenticator *authenticator.Authenticator,
	requestProvider rent_providers.IRequestProvider,
	instanceAcc *access.Instance,
	userAcc *access.User,
	pickUpPointAcc *access.PickUpPoint,
) rent.IRequestService {
	return &requestService{
		requestRepoProviders{
			requestProvider,
		},
		requestAccessors{
			instanceAcc,
			userAcc,
			pickUpPointAcc,
		},
		smachine,
		authenticator,
	}
}

func mapRequest(value *requests.Rent) rent.RentRequest {
	return rent.RentRequest{
		Id:               value.Id,
		InstanceId:       value.InstanceId,
		UserId:           value.UserId,
		PickUpPointId:    value.PickUpPointId,
		PaymentPeriodId:  value.PaymentPeriodId,
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

func (self *requestService) ListRentRequstsByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[rent.RentRequest], error) {
	var requests Collection[requests.Rent]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.request.GetRentRequestRepository()
		requests, err = repo.GetByUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapRequest, requests), err
}

func (self *requestService) GetRentRequestByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (rent.RentRequest, error) {
	var request requests.Rent
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		repo := self.repos.request.GetRentRequestRepository()
		request, err = repo.GetByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return mapRequest(&request), err
}

func (self *requestService) ListRentRequstsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[rent.RentRequest], error) {
	var requests Collection[requests.Rent]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.request.GetRentRequestRepository()
		requests, err = repo.GetByPickUpPointId(pickUpPointId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapRequest, requests), err
}

func (self *requestService) CreateRentRequest(
	token token.Token,
	form rent.RequestCreateForm,
) (rent.RentRequest, error) {
	var request requests.Rent
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		request, err = self.smachine.CreateRentRequest(
			form.InstanceId,
			user.Id,
			form.PickUpPointId,
			form.PaymentPeriodId,
		)
		err = states.MapError(err)
	}

	return mapRequest(&request), err
}

