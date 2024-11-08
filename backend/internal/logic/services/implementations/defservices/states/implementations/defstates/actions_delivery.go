package defstates

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	delivery_creator "rent_service/internal/logic/delivery"
	"rent_service/internal/logic/services/errors/cmnerrors"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"time"

	"github.com/google/uuid"
)

func (self *InstanceStateMachine) actionCreateDelivery(
	instance models.Instance,
	storage records.Storage,
	toId uuid.UUID,
) (requests.Delivery, error) {
	var delivery requests.Delivery
	var companyResponse delivery_creator.Delivery
	var from models.PickUpPoint
	var to models.PickUpPoint
	var err error
	code := self.code.Generate()

	repo := self.repos.pickUpPoint.GetPickUpPointRepository()
	from, err = repo.GetById(storage.PickUpPointId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err {
		repo := self.repos.pickUpPoint.GetPickUpPointRepository()
		to, err = repo.GetById(toId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		companyResponse, err = self.delivery.CreateDelivery(
			from.Address,
			to.Address,
			code,
		)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	if nil == err {
		repo := self.repos.delivery.GetDeliveryRepository()
		delivery = requests.Delivery{
			CompanyId:          companyResponse.CompanyId,
			InstanceId:         instance.Id,
			FromId:             storage.PickUpPointId,
			ToId:               toId,
			DeliveryId:         companyResponse.DeliveryId,
			ScheduledBeginDate: companyResponse.ScheduledBeginDate,
			ScheduledEndDate:   companyResponse.ScheduledEndDate,
			VerificationCode:   code,
			CreateDate:         time.Now(),
		}

		delivery, err = repo.Create(delivery)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return delivery, err
}

func (self *InstanceStateMachine) actionSendDelivery(
	delivery requests.Delivery,
	verificationCode string,
) error {
	var err error

	if verificationCode != delivery.VerificationCode {
		err = cmnerrors.Incorrect("verification_code")
	}

	if nil == err {
		delivery.ActualBeginDate = new(time.Time)
		*delivery.ActualBeginDate = time.Now()

		repo := self.repos.delivery.GetDeliveryRepository()
		err = repo.Update(delivery)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *InstanceStateMachine) actionAcceptDelivery(
	delivery requests.Delivery,
	verificationCode string,
) error {
	var err error

	if verificationCode != delivery.VerificationCode {
		err = cmnerrors.Incorrect("verification_code")
	}

	if nil == err {
		delivery.ActualEndDate = new(time.Time)
		*delivery.ActualEndDate = time.Now()

		repo := self.repos.delivery.GetDeliveryRepository()
		err = repo.Update(delivery)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

