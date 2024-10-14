package states

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/provide"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"time"

	"github.com/google/uuid"
)

func (self *InstanceStateMachine) actionCreateProvision(
	instance models.Instance,
	request requests.Provide,
) (records.Provision, error) {
	repo := self.repos.provision.self.GetProvisionRepository()
	provision, err := repo.Create(records.Provision{
		RenterId:   request.RenterId,
		InstanceId: instance.Id,
		StartDate:  time.Now(),
	})

	if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return provision, err
}

func (self *InstanceStateMachine) actionStopProvision(
	provision records.Provision,
) error {
	repo := self.repos.provision.self.GetProvisionRepository()
	provision.EndDate = new(time.Time)
	*provision.EndDate = time.Now()

	err := repo.Update(provision)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *InstanceStateMachine) actionCreateProvisionRequest(
	renterId uuid.UUID,
	form provide.RequestCreateForm,
) (requests.Provide, error) {
	var err error
	repo := self.repos.provision.request.GetProvisionRequestRepository()
	request := requests.Provide{
		ProductId:        form.ProductId,
		RenterId:         renterId,
		PickUpPointId:    form.PickUpPointId,
		Name:             form.Name,
		Description:      form.Description,
		Condition:        form.Condition,
		VerificationCode: self.code.Generate(),
		CreateDate:       time.Now(),
	}

	request.PayPlans = make(map[uuid.UUID]models.PayPlan)

	for _, v := range form.PayPlans {
		request.PayPlans[v.PeriodId] = models.PayPlan{
			PeriodId: v.PeriodId,
			Price:    v.Price,
		}
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

func (self *InstanceStateMachine) actionRejectProvisionRequest(
	request requests.Provide,
) error {
	repo := self.repos.provision.request.GetProvisionRequestRepository()
	err := repo.Remove(request.Id)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *InstanceStateMachine) actionAcceptProvisionRequest(
	request requests.Provide,
) error {
	repo := self.repos.provision.request.GetProvisionRequestRepository()
	err := repo.Remove(request.Id)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

func (self *InstanceStateMachine) actionCreateProvisionReturn(
	instance models.Instance,
	provision records.Provision,
	pickUpPointId uuid.UUID,
) (requests.Revoke, error) {
	var err error
	repo := self.repos.provision.revoke.GetRevokeProvisionRepository()
	request := requests.Revoke{
		InstanceId:       instance.Id,
		RenterId:         provision.RenterId,
		PickUpPointId:    pickUpPointId,
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

func (self *InstanceStateMachine) actionAcceptProvisionReturn(
	request requests.Revoke,
	verificationCode string,
) error {
	var err error

	if verificationCode != request.VerificationCode {
		err = cmnerrors.Incorrect("verification_code")
	}

	if nil == err {
		repo := self.repos.provision.revoke.GetRevokeProvisionRepository()
		err = repo.Remove(request.Id)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

func (self *InstanceStateMachine) actionCancelProvisionReturn(
	request requests.Revoke,
) error {
	repo := self.repos.provision.revoke.GetRevokeProvisionRepository()
	err := repo.Remove(request.Id)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

