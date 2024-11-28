package period_test

import (
	"rent_service/builders/misc/generator"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/interfaces/period"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/testcommon"
	"rent_service/misc/testcommon/defservices"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func MapPeriod(value *models.Period) period.Period {
	return period.Period{
		Id:       value.Id,
		Duration: value.Duration,
	}
}

type PeriodServiceIntegrationTestSuite struct {
	suite.Suite
	service  period.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *PeriodServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	// t.Parallel()
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreatePeriodService()
	})
}

func (self *PeriodServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *PeriodServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Period service",
	)
}

var describeGetPeriods = testcommon.MethodDescriptor(
	"GetPeriods",
	"Get list of all available periods",
)

func (self *PeriodServiceIntegrationTestSuite) TestGetPeriodsPositive(t provider.T) {
	var reference []period.Period

	describeGetPeriods(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepNewList(t, "periods",
				func(i uint) (models.Period, uuid.UUID) {
					var p models.Period
					switch i % 6 {
					case 0:
						p = models_om.PeriodDay().Build()
					case 1:
						p = models_om.PeriodWeek().Build()
					case 2:
						p = models_om.PeriodMonth().Build()
					case 3:
						p = models_om.PeriodQuarter().Build()
					case 4:
						p = models_om.PeriodHalf().Build()
					case 5:
						p = models_om.PeriodYear().Build()
					}

					reference = append(reference, MapPeriod(&p))

					return p, p.Id
				},
				self.rContext.Inserter.InsertPeriod,
			), 6).
			Generate().
			Finish()
	})

	// Act
	var result collection.Collection[period.Period]
	var err error

	t.WithNewStep("Get all periods", func(sCtx provider.StepCtx) {
		result, err = self.service.GetPeriods()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[period.Period](t).ContainsMultipleFunc(
		testcommon.DeepEqual[period.Period](),
		collection.Collect(result.Iter()), reference,
		"All values must be returned",
	)
}

func TestPeriodServiceIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(PeriodServiceIntegrationTestSuite))
}

