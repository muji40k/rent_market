package address_test

import (
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/implementation/sql/repositories/address"
	"rent_service/misc/testcommon"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type AddressRepositoryTestSuite struct {
	suite.Suite
	inserter *psql.Inserter
}

func (self *AddressRepositoryTestSuite) BeforeAll(t provider.T) {
	self.inserter = psql.NewInserter()
}

func (self *AddressRepositoryTestSuite) AfterAll(t provider.T) {
	self.inserter.Close()
}

func (self *AddressRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Address repository",
	)
	self.inserter.ClearDB()
}

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get address by id (used only by pick up point repository)",
)

func (self *AddressRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var repo *address.Repository

	var reference models.Address

	describeGetById(t,
		"Simple return all test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert reference address", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "address",
				models_om.AddressExmapleWithFlat("1").Build(),
			)
			self.inserter.InsertAddress(&reference)
		})

		t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
			factory, err := psql.PSQLRepositoryFactory().Build()

			if nil != err {
				t.Breakf("Unable to create repository: %s", err)
			}

			repo = factory.CreateAddressRepository()
		})
	})

	// Act
	var result models.Address
	var err error

	t.WithNewStep("Get all addresses", func(sCtx provider.StepCtx) {
		result, err = repo.GetById(reference.Id)
	}, allure.NewParameter("addressId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Same value")
}

func (self *AddressRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var repo *address.Repository

	var id uuid.UUID

	describeGetById(t,
		"Address not found",
		"Checks that method return error NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknonw id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
			factory, err := psql.PSQLRepositoryFactory().Build()

			if nil != err {
				t.Breakf("Unable to create repository: %s", err)
			}

			repo = factory.CreateAddressRepository()
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all addresses", func(sCtx provider.StepCtx) {
		_, err = repo.GetById(id)
	}, allure.NewParameter("addressId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestAddressRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(AddressRepositoryTestSuite))
}

