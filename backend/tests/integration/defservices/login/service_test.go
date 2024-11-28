package login_test

import (
	"rent_service/builders/misc/generator"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/login"
	"rent_service/internal/logic/services/types/token"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"rent_service/misc/testcommon/defservices"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type LoginServiceIntegrationTestSuite struct {
	suite.Suite
	service  login.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *LoginServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	// t.Parallel()
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateLoginService()
	})
}

func (self *LoginServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *LoginServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Login service",
	)
}

var describeRegister = testcommon.MethodDescriptor(
	"Register",
	"Register new user",
)

var describeLogin = testcommon.MethodDescriptor(
	"Login",
	"Login user",
)

func (self *LoginServiceIntegrationTestSuite) TestLoginPositive(t provider.T) {
	var user models.User

	describeLogin(t,
		"Simple login test",
		"Check that registered user can login without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepValue(t, "user", &user,
				func() (models.User, uuid.UUID) {
					u := models_om.UserDefault(nullable.None[string]()).Build()
					return u, u.Id
				},
				self.rContext.Inserter.InsertUser,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var result token.Token
	var err error

	t.WithNewStep("Login",
		func(sCtx provider.StepCtx) {
			result, err = self.service.Login(user.Email, user.Password)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", user.Password),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(token.Token(user.Token), result,
		"Right token is returned",
	)
}

func (self *LoginServiceIntegrationTestSuite) TestLoginWrongPassword(t provider.T) {
	var user models.User

	describeLogin(t,
		"Attemp to loing with wrong password",
		"Check that user can't login with wrong password",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepValue(t, "user", &user,
				func() (models.User, uuid.UUID) {
					u := models_om.UserDefault(nullable.None[string]()).
						WithPassword("SingleCorrectPassword").
						Build()
					return u, u.Id
				},
				self.rContext.Inserter.InsertUser,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var err error
	var incorrectPassword = "IncorrectPassword"

	t.WithNewStep("Login",
		func(sCtx provider.StepCtx) {
			_, err = self.service.Login(user.Email, incorrectPassword)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", incorrectPassword),
	)

	// Assert
	var aerr cmnerrors.ErrorAuthentication

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &aerr, "Error is authentication")
}

func (self *LoginServiceIntegrationTestSuite) TestRegisterPositive(t provider.T) {
	var user models.User

	describeRegister(t,
		"Simple register test",
		"Check that new user was added to repository",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Register",
		func(sCtx provider.StepCtx) {
			err = self.service.Register(user.Email, user.Password, user.Name)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", user.Password),
		allure.NewParameter("username", user.Name),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *LoginServiceIntegrationTestSuite) TestRegisterUserAlreadyExists(t provider.T) {
	var user models.User

	describeRegister(t,
		"User already exists",
		"Check that attempt to register user with same email fails with AlreadyExists error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepValue(t, "user", &user,
				func() (models.User, uuid.UUID) {
					u := models_om.UserDefault(nullable.None[string]()).Build()
					return u, u.Id
				},
				self.rContext.Inserter.InsertUser,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var err error

	t.WithNewStep("Register",
		func(sCtx provider.StepCtx) {
			err = self.service.Register(user.Email, user.Password, user.Name)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", user.Password),
		allure.NewParameter("username", user.Name),
	)

	// Assert
	var aeerr cmnerrors.ErrorAlreadyExists

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &aeerr, "Error must be AlreadyExists")
}

func TestLoginServiceIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(LoginServiceIntegrationTestSuite))
}

