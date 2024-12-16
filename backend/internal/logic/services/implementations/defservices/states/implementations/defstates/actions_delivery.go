package defstates

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	delivery_creator "rent_service/internal/logic/delivery"
	"rent_service/internal/logic/services/errors/cmnerrors"
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
	code := self.code.Generate()

	err := cmnerrors.RepoCallWrap(func() (err error) {
		repo := self.repos.pickUpPoint.GetPickUpPointRepository()
		from, err = repo.GetById(storage.PickUpPointId)
		return
	})

	if nil == err {
		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.pickUpPoint.GetPickUpPointRepository()
			to, err = repo.GetById(toId)
			return
		})
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

		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.delivery.GetDeliveryRepository()
			delivery, err = repo.Create(delivery)
			return
		})
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

		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.delivery.GetDeliveryRepository()
			err = repo.Update(delivery)
			return
		})
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

		err = cmnerrors.RepoCallWrap(func() (err error) {
			repo := self.repos.delivery.GetDeliveryRepository()
			err = repo.Update(delivery)
			return
		})
	}

	return err
}

