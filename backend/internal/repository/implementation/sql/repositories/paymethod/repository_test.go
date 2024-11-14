package paymethod_test

import (
	"rent_service/builders/misc/collect"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/paymethod"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type PayMethodRepositoryTestSuite struct {
	suite.Suite
	repo paymethod.IRepository
	psqlcommon.Context
}

func (self *PayMethodRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PayMethodRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PayMethodRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Pay Method repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePayMethodRepository()
	})
}

var describeGetAll = testcommon.MethodDescriptor(
	"GetAll",
	"Get app pay methods",
)

func (self *PayMethodRepositoryTestSuite) TestGetAllPositive(t provider.T) {
	var reference []models.PayMethod

	describeGetAll(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert pay methods", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "methods",
				collect.Do(
					models_om.PayMethodExample("1"),
					models_om.PayMethodExample("2"),
					models_om.PayMethodExample("3"),
					models_om.PayMethodExample("4"),
					models_om.PayMethodExample("5"),
				),
			)
			psql.BulkInsert(self.Inserter.InsertPayMethod, reference...)
		})
	})

	// Act
	var result collection.Collection[models.PayMethod]
	var err error

	t.WithNewStep("Get all pay methods", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(reference, collection.Collect(result.Iter()),
		"Same company values")
}

func (self *PayMethodRepositoryTestSuite) TestGetAllEmpty(t provider.T) {
	describeGetAll(t,
		"No pay methods",
		"Checks that method return empty collection withour error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[models.PayMethod]
	var err error

	t.WithNewStep("Get all pay methods", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()),
		"Collection is empty")
}

func TestPayMethodRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PayMethodRepositoryTestSuite))
}

