package states

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/provide"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"strings"

	"github.com/google/uuid"
)

func override[T any](dst *T, override *T) {
	if nil != override {
		*dst = *override
	}
}

func overrideMap(
	dst *map[uuid.UUID]models.PayPlan,
	override []provide.PayPlan,
) error {
	if nil == override {
		return nil
	}

	var err error
	previous := *dst
	*dst = make(map[uuid.UUID]models.PayPlan)

	for i := 0; nil == err && len(override) > i; i++ {
		v := override[i]

		if plan, found := previous[v.PeriodId]; !found || plan.Id != v.Id {
			err = cmnerrors.Incorrect("Pay plan doesn't exist")
		} else {
			(*dst)[v.PeriodId] = models.PayPlan{
				Id:       v.Id,
				PeriodId: v.PeriodId,
				Price:    v.Price,
			}
		}
	}

	return err
}

func (self *InstanceStateMachine) actionCreateInstance(
	request requests.Provide,
	form provide.StartForm,
) (models.Instance, error) {
	var instance models.Instance
	override(&request.ProductId, form.Overrides.ProductId)
	override(&request.Name, form.Overrides.Name)
	override(&request.Description, form.Overrides.Description)
	override(&request.Condition, form.Overrides.Condition)
	err := overrideMap(&request.PayPlans, form.Overrides.PayPlans)

	if nil == err {
		repo := self.repos.instance.self.GetInstanceRepository()
		instance, err = repo.Create(models.Instance{
			ProductId:   request.ProductId,
			Name:        request.Name,
			Description: request.Description,
			Condition:   request.Condition,
		})

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	if nil == err {
		repo := self.repos.instance.payPlans.GetInstancePayPlansRepository()
		_, err = repo.Create(models.InstancePayPlans{
			InstanceId: instance.Id,
			Map:        request.PayPlans,
		})

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.AlreadyExists(cerr.What...))
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return instance, err
}

type instancef func(*models.Instance)

func instanceUpdateDescription(description *string) instancef {
	if nil == description || "" == *description {
		return nil
	}

	tdescription := strings.TrimSpace(*description)

	return func(instance *models.Instance) {
		instance.Description += "\n" + tdescription
	}
}

func (self *InstanceStateMachine) actionUpdateInstance(
	instance models.Instance,
	f instancef,
) error {
	if nil == f {
		return nil
	}

	f(&instance)

	repo := self.repos.instance.self.GetInstanceRepository()
	err := repo.Update(instance)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.Internal(cmnerrors.NotFound(cerr.What...))
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return err
}

