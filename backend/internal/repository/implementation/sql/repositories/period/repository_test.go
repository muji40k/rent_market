package period_test

import (
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/period"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type PeriodRepositoryTestSuite struct {
	suite.Suite
	repo period.IRepository
	psqlcommon.Context
}

func (self *PeriodRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PeriodRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PeriodRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Period repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePeriodRepository()
	})
}

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get period by id",
)

var describeGetAll = testcommon.MethodDescriptor(
	"GetAll",
	"Get all delivery comapnies",
)

func (self *PeriodRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var reference models.Period

	describeGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert periods", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "period",
				models_om.PeriodWeek().Build(),
			)
			self.Inserter.InsertPeriod(&reference)
		})
	})

	// Act
	var result models.Period
	var err error

	t.WithNewStep("Get period by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("periodId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same period value")
}

func (self *PeriodRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Delivery company not found",
		"Checks that method return error NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknonw id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})
	})

	// Act
	var err error

	t.WithNewStep("Get period by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("periodId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *PeriodRepositoryTestSuite) TestGetAllPositive(t provider.T) {
	var reference []models.Period

	describeGetAll(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert periods", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "periods",
				collect.Do(
					models_om.PeriodDay(),
					models_om.PeriodWeek(),
					models_om.PeriodMonth(),
					models_om.PeriodQuarter(),
					models_om.PeriodHalf(),
					models_om.PeriodYear(),
				),
			)
			psql.BulkInsert(self.Inserter.InsertPeriod, reference...)
		})
	})

	// Act
	var result collection.Collection[models.Period]
	var err error

	t.WithNewStep("Get all periods", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(reference, collection.Collect(result.Iter()),
		"Same company values")
}

func (self *PeriodRepositoryTestSuite) TestGetAllEmpty(t provider.T) {
	describeGetAll(t,
		"No periods",
		"Checks that method return empty collection withour error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[models.Period]
	var err error

	t.WithNewStep("Get all periods", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()),
		"Collection is empty")
}

func TestPeriodRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PeriodRepositoryTestSuite))
}

