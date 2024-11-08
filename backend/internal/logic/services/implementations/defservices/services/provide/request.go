package provide

import (
	"errors"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/states"
	"rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	provision_providers "rent_service/internal/repository/context/providers/provision"
	role_providers "rent_service/internal/repository/context/providers/role"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

type requestRepoProviders struct {
	renter  role_providers.IRenterProvider
	request provision_providers.IRequestProvider
}

type requestAccessors struct {
	instance    access.IInstance
	user        access.IUser
	pickUpPoint access.IPickUpPoint
}

type requestService struct {
	repos         requestRepoProviders
	access        requestAccessors
	smachine      states.IInstanceStateMachine
	authenticator authenticator.IAuthenticator
	authorizer    authorizer.IAuthorizer
}

func NewRequest(
	smachine states.IInstanceStateMachine,
	authenticator authenticator.IAuthenticator,
	authorizer authorizer.IAuthorizer,
	renterProvider role_providers.IRenterProvider,
	requestProvider provision_providers.IRequestProvider,
	instanceAcc access.IInstance,
	userAcc access.IUser,
	pickUpPointAcc access.IPickUpPoint,
) provide.IRequestService {
	return &requestService{
		requestRepoProviders{
			renterProvider,
			requestProvider,
		},
		requestAccessors{
			instanceAcc,
			userAcc,
			pickUpPointAcc,
		},
		smachine,
		authenticator,
		authorizer,
	}
}

func (self *requestService) mapRequest(value *requests.Provide) provide.ProvideRequest {
	var userId uuid.UUID
	repo := self.repos.renter.GetRenterRepository()

	renter, err := repo.GetById(value.RenterId)

	if nil == err {
		userId = renter.UserId
	}

	payPlans := make([]provide.PayPlan, 0, len(value.PayPlans))

	for _, v := range value.PayPlans {
		payPlans = append(payPlans, provide.PayPlan{
			Id:       v.Id,
			PeriodId: v.PeriodId,
			Price:    v.Price,
		})
	}

	return provide.ProvideRequest{
		Id:               value.Id,
		ProductId:        value.ProductId,
		UserId:           userId,
		PickUpPointId:    value.PickUpPointId,
		Name:             value.Name,
		Description:      value.Description,
		Condition:        value.Condition,
		PayPlans:         payPlans,
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

func (self *requestService) ListProvisionRequstsByUser(
	token token.Token,
	userId uuid.UUID,
) (Collection[provide.ProvideRequest], error) {
	var requests Collection[requests.Provide]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.user.Access(user.Id, userId)
	}

	if nil == err {
		repo := self.repos.request.GetProvisionRequestRepository()
		requests, err = repo.GetByUserId(userId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(self.mapRequest, requests), err
}

func (self *requestService) GetProvisionRequestByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (provide.ProvideRequest, error) {
	var request requests.Provide
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if nil == err {
		repo := self.repos.request.GetProvisionRequestRepository()
		request, err = repo.GetByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return self.mapRequest(&request), err
}

func (self *requestService) ListProvisionRequstsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[provide.ProvideRequest], error) {
	var requests Collection[requests.Provide]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.request.GetProvisionRequestRepository()
		requests, err = repo.GetByPickUpPointId(pickUpPointId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(self.mapRequest, requests), err
}

func (self *requestService) CreateProvisionRequest(
	token token.Token,
	form provide.RequestCreateForm,
) (provide.ProvideRequest, error) {
	var request requests.Provide
	var renter models.Renter
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		renter, err = self.authorizer.IsRenter(user.Id)
	}

	if nil == err {
		request, err = self.smachine.CreateProvisionRequest(
			renter.Id,
			form,
		)
		err = states.MapError(err)
	}

	return self.mapRequest(&request), err
}

