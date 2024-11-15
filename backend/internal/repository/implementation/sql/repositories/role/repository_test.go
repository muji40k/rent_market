package role_test

import (
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/role"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type AdministratorRepositoryTestSuite struct {
	suite.Suite
	repo role.IAdministratorRepository
	psqlcommon.Context
}

func (self *AdministratorRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *AdministratorRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *AdministratorRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Role Administrator repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateRoleAdministratorRepository()
	})
}

var describeAdministratorGetByUserId = testcommon.MethodDescriptor(
	"GetByUserId",
	"Get administrator record with user id",
)

func (self *AdministratorRepositoryTestSuite) TestGetByUserIdPositive(t provider.T) {
	var (
		user      models.User
		reference models.Administrator
	)

	describeAdministratorGetByUserId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserAdministrator(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert administrator", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "administrator",
				models_om.AdministratorWithUserId(user.Id).Build(),
			)
			self.Inserter.InsertAdministrator(&reference)
		})
	})

	// Act
	var result models.Administrator
	var err error

	t.WithNewStep("Get administrator by user id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetByUserId(user.Id)
	}, allure.NewParameter("userId", user.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same administrator value")
}

func (self *AdministratorRepositoryTestSuite) TestGetByUserIdNotFound(t provider.T) {
	var id uuid.UUID

	describeAdministratorGetByUserId(t,
		"User not found",
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

	t.WithNewStep("Get administrator by user id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByUserId(id)
	}, allure.NewParameter("userId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type RenterRepositoryTestSuite struct {
	suite.Suite
	repo role.IRenterRepository
	psqlcommon.Context
}

func (self *RenterRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *RenterRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *RenterRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Role Renter repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateRoleRenterRepository()
	})
}

var describeRenterCreate = testcommon.MethodDescriptor(
	"Create",
	"Register user as renter",
)

var describeRenterGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get renter record with user id",
)

var describeRenterGetByUserId = testcommon.MethodDescriptor(
	"GetByUserId",
	"Get renter record with user id",
)

func CompareCreated(exp models.Renter, act models.Renter) bool {
	return uuid.UUID{} != act.Id && exp.UserId == act.UserId
}

func (self *RenterRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var (
		user      models.User
		reference models.Renter
	)

	describeRenterCreate(t,
		"Simple create test",
		"Checks that method returns value with set uuid without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserRenter(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create renter", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "renter",
				models_om.RenterWithUserId(user.Id).
					WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result models.Renter
	var err error

	t.WithNewStep("Get renter by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(user.Id)
	}, allure.NewParameter("userId", user.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.Renter](t).EqualFunc(
		CompareCreated, reference, result, "Same renter value",
	)
}

func (self *RenterRepositoryTestSuite) TestCreateDuplicate(t provider.T) {
	var user models.User

	describeRenterCreate(t,
		"Second attempt to create renter record for user",
		"Checks that method returns error and its Duplicate",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserRenter(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert renter", func(sCtx provider.StepCtx) {
			builder := models_om.RenterWithUserId(user.Id)
			created := testcommon.AssignParameter(sCtx, "created",
				builder.Build(),
			)
			self.Inserter.InsertRenter(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Get renter by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(user.Id)
	}, allure.NewParameter("userId", user.Id))

	// Assert
	var derr cmnerrors.ErrorDuplicate
	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &derr, "Error is Duplicate")
}

func (self *RenterRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var (
		user      models.User
		reference models.Renter
	)

	describeRenterGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserRenter(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert renter", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "renter",
				models_om.RenterWithUserId(user.Id).Build(),
			)
			self.Inserter.InsertRenter(&reference)
		})
	})

	// Act
	var result models.Renter
	var err error

	t.WithNewStep("Get renter by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("retnerId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same renter value")
}

func (self *RenterRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeRenterGetById(t,
		"Renter not found",
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

	t.WithNewStep("Get renter by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("renterId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *RenterRepositoryTestSuite) TestGetByUserIdPositive(t provider.T) {
	var (
		user      models.User
		reference models.Renter
	)

	describeRenterGetByUserId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserRenter(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert renter", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "renter",
				models_om.RenterWithUserId(user.Id).Build(),
			)
			self.Inserter.InsertRenter(&reference)
		})
	})

	// Act
	var result models.Renter
	var err error

	t.WithNewStep("Get renter by user id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetByUserId(user.Id)
	}, allure.NewParameter("userId", user.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same renter value")
}

func (self *RenterRepositoryTestSuite) TestGetByUserIdNotFound(t provider.T) {
	var id uuid.UUID

	describeRenterGetByUserId(t,
		"User not found",
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

	t.WithNewStep("Get renter by user id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByUserId(id)
	}, allure.NewParameter("userId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type StorekeeperRepositoryTestSuite struct {
	suite.Suite
	repo role.IStorekeeperRepository
	psqlcommon.Context
}

func (self *StorekeeperRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *StorekeeperRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *StorekeeperRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Role Storekeeper repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateRoleStorekeeperRepository()
	})
}

var describeStorekeeperGetByUserId = testcommon.MethodDescriptor(
	"GetByUserId",
	"Get storekeeper record with user id",
)

func (self *StorekeeperRepositoryTestSuite) TestGetByUserIdPositive(t provider.T) {
	var (
		pup       models.PickUpPoint
		user      models.User
		reference models.Storekeeper
	)

	describeStorekeeperGetByUserId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create and insert user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
			self.Inserter.InsertUser(&user)
		})

		t.WithNewStep("Create and insert storekeeper", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storekeeper",
				models_om.StorekeeperWithUserId(user.Id, pup.Id).Build(),
			)
			self.Inserter.InsertStorekeeper(&reference)
		})
	})

	// Act
	var result models.Storekeeper
	var err error

	t.WithNewStep("Get storekeeper by user id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetByUserId(user.Id)
	}, allure.NewParameter("userId", user.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same storekeeper value")
}

func (self *StorekeeperRepositoryTestSuite) TestGetByUserIdNotFound(t provider.T) {
	var id uuid.UUID

	describeStorekeeperGetByUserId(t,
		"User not found",
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

	t.WithNewStep("Get storekeeper by user id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByUserId(id)
	}, allure.NewParameter("userId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestAdministratorRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(AdministratorRepositoryTestSuite))
}

func TestRenterRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(RenterRepositoryTestSuite))
}

func TestStorekeeperRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(StorekeeperRepositoryTestSuite))
}

