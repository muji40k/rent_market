package currency_test

import (
	"rent_service/builders/misc/uuidgen"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/misc/types/currency"
	"rent_service/internal/repository/errors/cmnerrors"
	rcurrency "rent_service/internal/repository/implementation/sql/repositories/currency"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CurrencyRepositoryTestSuite struct {
	suite.Suite
	repo *rcurrency.Repository
	psqlcommon.Context
}

func (self *CurrencyRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *CurrencyRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *CurrencyRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Currency repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateCurrencyRepository()
	})
}

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get currency by id",
)

var describeGetId = testcommon.MethodDescriptor(
	"GetId",
	"Get get currency id by name",
)

func (self *CurrencyRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var name string = "rub"
	var id uuid.UUID = psql.GetCurrency(name)

	describeGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	// Empty

	// Act
	var result currency.Currency
	var err error

	t.WithNewStep("Get currency by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(id)
	}, allure.NewParameter("currencyId", id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(name, result.Name, "Same values")
}

func (self *CurrencyRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Currency not found",
		"Checks that method return error NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknonw id", func(sCtx provider.StepCtx) {
			ids := psql.GetAllCurrencies()
			id = uuidgen.Generate()

			for slices.Contains(ids, id) {
				id = uuidgen.Generate()
			}

			sCtx.WithParameters(allure.NewParameter("id", id))
		})
	})

	// Act
	var err error

	t.WithNewStep("Get currency", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("currencyId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *CurrencyRepositoryTestSuite) TestGetIdPositive(t provider.T) {
	var name string = "rub"
	var id uuid.UUID = psql.GetCurrency(name)

	describeGetId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	// Empty

	// Act
	var result uuid.UUID
	var err error

	t.WithNewStep("Get currency by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetId(name)
	}, allure.NewParameter("currencyName", name))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(id, result, "Same values")
}

func (self *CurrencyRepositoryTestSuite) TestGetByNotFound(t provider.T) {
	var name string = "definetly unknown currency"

	describeGetId(t,
		"Currency not found",
		"Checks that method return error NotFound",
	)

	// Arrange
	// Empty

	// Act
	var err error

	t.WithNewStep("Get currency", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetId(name)
	}, allure.NewParameter("currencyName", name))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestCurrencyRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CurrencyRepositoryTestSuite))
}

