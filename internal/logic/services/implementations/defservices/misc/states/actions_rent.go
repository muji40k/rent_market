package states

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"time"

	"github.com/google/uuid"
)

func (self *InstanceStateMachine) actionCreateRent(
	request requests.Rent,
) (records.Rent, error) {
	repo := self.repos.rent.self.GetRentRepository()
	rent, err := repo.Create(records.Rent{
		UserId:          request.UserId,
		InstanceId:      request.InstanceId,
		StartDate:       time.Now(),
		PaymentPeriodId: request.PaymentPeriodId,
	})

	if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return rent, err
}

func (self *InstanceStateMachine) actionStopRent(
	rent records.Rent,
) error {
	rent.EndDate = new(time.Time)
	*rent.EndDate = time.Now()

	repo := self.repos.rent.self.GetRentRepository()
	err := repo.Update(rent)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *InstanceStateMachine) actionCreateRentRequest(
	instance models.Instance,
	userId uuid.UUID,
	pickUpPointId uuid.UUID,
	paymentPeriodId uuid.UUID,
) (requests.Rent, error) {
	var err error
	repo := self.repos.rent.request.GetRentRequestRepository()
	request := requests.Rent{
		InstanceId:       instance.Id,
		UserId:           userId,
		PickUpPointId:    pickUpPointId,
		PaymentPeriodId:  paymentPeriodId,
		VerificationCode: self.code.Generate(),
		CreateDate:       time.Now(),
	}

	request, err = repo.Create(request)

	if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return request, err
}

func (self *InstanceStateMachine) actionAcceptRentRequest(
	request requests.Rent,
	verificationCode string,
) error {
	var err error

	if verificationCode != request.VerificationCode {
		err = cmnerrors.Incorrect("verification_code")
	}

	if nil == err {
		repo := self.repos.rent.request.GetRentRequestRepository()
		err := repo.Remove(request.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *InstanceStateMachine) actionRejectRentRequest(
	request requests.Rent,
) error {
	repo := self.repos.rent.request.GetRentRequestRepository()
	err := repo.Remove(request.Id)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *InstanceStateMachine) actionCreateRentReturn(
	instance models.Instance,
	rent records.Rent,
	pickUpPointId uuid.UUID,
	endDate time.Time,
) (requests.Return, error) {
	end := endDate
	repo := self.repos.period.GetPeriodRepository()
	period, err := repo.GetById(rent.PaymentPeriodId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err {
		end = endDate.Truncate(24 * time.Hour)           // Valid, since all
		start := rent.StartDate.Truncate(24 * time.Hour) // periods are multiple
		duration := end.Sub(start)                       // of 1 day

		if start.After(end) {
			err = cmnerrors.Incorrect("End date before start")
		} else if offset := duration % period.Duration; 0 != offset {
			end.Add(period.Duration - offset)
		}
	}

	var request requests.Return
	if nil == err {
		repo := self.repos.rent.retrn.GetRentReturnRepository()
		request = requests.Return{
			InstanceId:       instance.Id,
			UserId:           rent.UserId,
			PickUpPointId:    pickUpPointId,
			RentEndDate:      end,
			VerificationCode: self.code.Generate(),
			CreateDate:       time.Now(),
		}

		request, err = repo.Create(request)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return request, err
}

func (self *InstanceStateMachine) actionAcceptRentReturn(
	request requests.Return,
	verificationCode string,
) error {
	var err error

	if verificationCode != request.VerificationCode {
		err = cmnerrors.Incorrect("verification_code")
	}

	if request.RentEndDate.After(time.Now()) {
		err = cmnerrors.Incorrect("Rent end date didn't pass")
	}

	if nil == err {
		repo := self.repos.rent.retrn.GetRentReturnRepository()
		err := repo.Remove(request.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *InstanceStateMachine) actionCancelRentReturn(
	request requests.Return,
) error {
	repo := self.repos.rent.retrn.GetRentReturnRepository()
	err := repo.Remove(request.Id)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

