package period_test

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/period"
	"rent_service/internal/logic/services/interfaces/period"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/testcommon"
	"testing"

	rperiod "rent_service/internal/repository/implementation/mock/period"

	period_pmock "rent_service/internal/repository/context/mock/period"

	"rent_service/builders/misc/collect"
	models_om "rent_service/builders/mothers/domain/models"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

func GetService(ctrl *gomock.Controller, f func(repo *rperiod.MockIRepository)) period.IService {
	repo := rperiod.NewMockIRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.New(period_pmock.New(repo))
}

func MapPeriod(value *models.Period) period.Period {
	return period.Period{
		Id:       value.Id,
		Duration: value.Duration,
	}
}

type PeriodServiceTestSuite struct {
	suite.Suite
}

func (self *PeriodServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default period implementation",
		"Period service",
	)
}

var describeGetPeriods = testcommon.MethodDescriptor(
	"GetPeriods",
	"Get list of all available periods",
)

func (self *PeriodServiceTestSuite) TestGetPeriodsPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service period.IService

	var periods []models.Period
	var reference []period.Period

	describeGetPeriods(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create reference periods", func(sCtx provider.StepCtx) {
			periods = testcommon.AssignParameter(sCtx, "periods",
				collect.Do(
					models_om.PeriodDay(),
					models_om.PeriodWeek(),
					models_om.PeriodMonth(),
					models_om.PeriodQuarter(),
					models_om.PeriodHalf(),
					models_om.PeriodYear(),
				),
			)

			reference = collection.Collect(
				collection.MapIterator(
					MapPeriod, collection.SliceIterator(periods),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rperiod.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(collection.SliceCollection(periods), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[period.Period]
	var err error

	t.WithNewStep("Get all periods", func(sCtx provider.StepCtx) {
		result, err = service.GetPeriods()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PeriodServiceTestSuite) TestGetPeriodsInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service period.IService

	describeGetPeriods(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rperiod.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(nil, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all periods", func(sCtx provider.StepCtx) {
		_, err = service.GetPeriods()
	})

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "No error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func TestPeriodServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PeriodServiceTestSuite))
}

