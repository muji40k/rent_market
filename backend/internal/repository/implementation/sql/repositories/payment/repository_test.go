package payment_test

import (
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/payment"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	repo payment.IRepository
	psqlcommon.Context
}

func (self *PaymentRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PaymentRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PaymentRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Payment repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePaymentRepository()
	})
}

var describeGetByInstanceId = testcommon.MethodDescriptor(
	"GetByInstanceId",
	"Get payments by instance id",
)

var describeGetByRentId = testcommon.MethodDescriptor(
	"GetByRentId",
	"Get payments by rent id",
)

func (self *PaymentRepositoryTestSuite) TestGetByInstanceIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetByInstanceId(t,
		"Instance not found",
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

	t.WithNewStep("Get payments for instance", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByInstanceId(id)
	}, allure.NewParameter("instanceId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *PaymentRepositoryTestSuite) TestGetByRentIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetByRentId(t,
		"Rent not found",
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

	t.WithNewStep("Get payments for rent", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByRentId(id)
	}, allure.NewParameter("rentId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPaymentRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PaymentRepositoryTestSuite))
}

