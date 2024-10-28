package delivery

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/misc/access"
	"rent_service/internal/logic/services/implementations/defservices/misc/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/misc/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/misc/states"
	"rent_service/internal/logic/services/interfaces/delivery"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
	delivery_providers "rent_service/internal/repository/context/providers/delivery"
	instance_providers "rent_service/internal/repository/context/providers/instance"
	rent_providers "rent_service/internal/repository/context/providers/rent"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	delivery    delivery_providers.IProvider
	photo       instance_providers.IPhotoProvider
	rentRequest rent_providers.IRequestProvider
}

type accessors struct {
	instance    *access.Instance
	pickUpPoint *access.PickUpPoint
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
	delivery delivery_providers.IProvider,
	photo instance_providers.IPhotoProvider,
	rentRequest rent_providers.IRequestProvider,
	instanceAcc *access.Instance,
	pickUpPointAcc *access.PickUpPoint,
) delivery.IService {
	return &service{
		repoproviders{delivery, photo, rentRequest},
		accessors{instanceAcc, pickUpPointAcc},
		smachine,
		authenticator,
		registry,
	}
}

func mapf(value *requests.Delivery) delivery.Delivery {
	out := delivery.Delivery{
		Id:         value.Id,
		CompanyId:  value.CompanyId,
		InstanceId: value.InstanceId,
		FromId:     value.FromId,
		ToId:       value.ToId,
		BeginDate: delivery.Dates{
			Scheduled: date.New(value.ScheduledBeginDate),
		},
		EndDate: delivery.Dates{
			Scheduled: date.New(value.ScheduledEndDate),
		},
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}

	if nil != value.ActualBeginDate {
		out.BeginDate.Actual = new(date.Date)
		*out.BeginDate.Actual = date.New(*value.ActualBeginDate)
	}

	if nil != value.ActualEndDate {
		out.EndDate.Actual = new(date.Date)
		*out.EndDate.Actual = date.New(*value.ActualEndDate)
	}

	return out
}

func (self *service) ListDeliveriesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (Collection[delivery.Delivery], error) {
	var deliveries Collection[requests.Delivery]
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, pickUpPointId)
	}

	if nil == err {
		repo := self.repos.delivery.GetDeliveryRepository()
		deliveries, err = repo.GetActiveByPickUpPointId(pickUpPointId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapf, deliveries), err
}

func (self *service) GetDeliveryByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (delivery.Delivery, error) {
	var delivery requests.Delivery
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.instance.Access(user.Id, instanceId)
	}

	if aerr := (cmnerrors.ErrorNoAccess{}); errors.As(err, &aerr) {
		repo := self.repos.rentRequest.GetRentRequestRepository()
		request, ierr := repo.GetByInstanceId(instanceId)

		if nil == ierr && request.UserId == user.Id {
			err = nil
		} else if cerr := (repo_errors.ErrorNotFound{}); !errors.As(ierr, &cerr) {
			err = cmnerrors.Internal(cmnerrors.DataAccess(ierr))
		}
	}

	if nil == err {
		repo := self.repos.delivery.GetDeliveryRepository()
		delivery, err = repo.GetActiveByInstanceId(instanceId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return mapf(&delivery), err
}

func (self *service) CreateDelivery(
	token token.Token,
	form delivery.CreateForm,
) (delivery.Delivery, error) {
	var delivery requests.Delivery
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, form.From)
	}

	if nil == err {
		err = self.access.instance.Access(user.Id, form.InstanceId)
	}

	if nil == err {
		delivery, err = self.smachine.CreateDelivery(
			form.InstanceId, form.From, form.To,
		)
		err = states.MapError(err)
	}

	return mapf(&delivery), err
}

func (self *service) SendDelivery(
	token token.Token,
	form delivery.SendForm,
) error {
	var delivery requests.Delivery
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.delivery.GetDeliveryRepository()
		delivery, err = repo.GetById(form.DeliveryId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, delivery.FromId)
	}

	if nil == err {
		err = states.MapError(self.smachine.SendDelivery(
			delivery.InstanceId,
			delivery.Id,
			form.VerificationCode,
		))
	}

	var ids []uuid.UUID
	if nil == err {
		ids, err = self.registry.MoveFromTemps(form.TempPhotos...)
	}

	if nil == err {
		repo := self.repos.photo.GetInstancePhotoRepository()

		for i := 0; len(ids) > i; i++ {
			cerr := repo.Create(delivery.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
	}

	return err
}

func (self *service) AcceptDelivery(
	token token.Token,
	form delivery.AcceptForm,
) error {
	var delivery requests.Delivery
	user, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.delivery.GetDeliveryRepository()
		delivery, err = repo.GetById(form.DeliveryId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		err = self.access.pickUpPoint.Access(user.Id, delivery.ToId)
	}

	if nil == err {
		err = states.MapError(self.smachine.AcceptDelivery(
			delivery.InstanceId,
			delivery.Id,
			form.Comment,
			form.VerificationCode,
		))
	}

	var ids []uuid.UUID
	if nil == err {
		ids, err = self.registry.MoveFromTemps(form.TempPhotos...)
	}

	if nil == err {
		repo := self.repos.photo.GetInstancePhotoRepository()

		for i := 0; len(ids) > i; i++ {
			cerr := repo.Create(delivery.InstanceId, ids[i])

			if nil == err {
				err = cerr
			}
		}
	}

	return err
}

type companyRepoProviders struct {
	company delivery_providers.ICompanyProvider
}

type companyService struct {
	repos         companyRepoProviders
	authenticator *authenticator.Authenticator
}

func NewCompany(
	authenticator *authenticator.Authenticator,
	company delivery_providers.ICompanyProvider,
) delivery.ICompanyService {
	return &companyService{companyRepoProviders{company}, authenticator}
}

func mapCompany(value *models.DeliveryCompany) delivery.DeliveryCompany {
	return delivery.DeliveryCompany{
		Id:          value.Id,
		Name:        value.Name,
		Site:        value.Site,
		PhoneNumber: value.PhoneNumber,
		Description: value.Description,
	}
}

func (self *companyService) ListDeliveryCompanies(
	token token.Token,
) (Collection[delivery.DeliveryCompany], error) {
	var companies Collection[models.DeliveryCompany]
	_, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.company.GetDeliveryCompanyRepository()
		companies, err = repo.GetAll()

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapCompany, companies), err
}

func (self *companyService) GetDeliveryCompanyById(
	token token.Token,
	companyId uuid.UUID,
) (delivery.DeliveryCompany, error) {
	var company models.DeliveryCompany
	_, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.company.GetDeliveryCompanyRepository()
		company, err = repo.GetById(companyId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return mapCompany(&company), err
}

