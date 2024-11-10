package login_test

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/login"
	"rent_service/internal/logic/services/interfaces/login"
	"rent_service/internal/logic/services/types/token"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"testing"

	ruser "rent_service/internal/repository/implementation/mock/user"

	user_pmock "rent_service/internal/repository/context/mock/user"

	models_om "rent_service/builders/mothers/domain/models"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

func GetService(ctrl *gomock.Controller, f func(repo *ruser.MockIRepository)) login.IService {
	repo := ruser.NewMockIRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.New(user_pmock.New(repo))
}

type LoginServiceTestSuite struct {
	suite.Suite
}

func (self *LoginServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
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

func (self *LoginServiceTestSuite) TestLoginPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service login.IService

	var user models.User

	describeLogin(t,
		"Simple login test",
		"Check that registered user can login without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *ruser.MockIRepository) {
				repo.EXPECT().GetByEmail(user.Email).
					Return(user, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result token.Token
	var err error

	t.WithNewStep("Login",
		func(sCtx provider.StepCtx) {
			result, err = service.Login(user.Email, user.Password)
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

func (self *LoginServiceTestSuite) TestLoginWrongPassword(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service login.IService

	var user models.User

	describeLogin(t,
		"Attemp to loing with wrong password",
		"Check that user can't login with wrong password",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).
					WithPassword("SingleCorrectPassword").
					Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *ruser.MockIRepository) {
				repo.EXPECT().GetByEmail(user.Email).
					Return(user, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error
	var incorrectPassword = "IncorrectPassword"

	t.WithNewStep("Login",
		func(sCtx provider.StepCtx) {
			_, err = service.Login(user.Email, incorrectPassword)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", incorrectPassword),
	)

	// Assert
	var aerr cmnerrors.ErrorAuthentication

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &aerr, "Error is authentication")
}

func (self *LoginServiceTestSuite) TestRegisterPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service login.IService

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

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			toBeInserted := models.User{
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			}
			service = GetService(ctrl, func(repo *ruser.MockIRepository) {
				repo.EXPECT().Create(toBeInserted).
					Return(user, nil).
					Times(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Register",
		func(sCtx provider.StepCtx) {
			err = service.Register(user.Email, user.Password, user.Name)
		},
		allure.NewParameter("email", user.Email),
		allure.NewParameter("password", user.Password),
		allure.NewParameter("username", user.Name),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *LoginServiceTestSuite) TestRegisterUserAlreadyExists(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service login.IService

	var user models.User

	describeRegister(t,
		"User already exists",
		"Check that attempt to register user with same email fails with AlreadyExists error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			toBeInserted := models.User{
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			}
			service = GetService(ctrl, func(repo *ruser.MockIRepository) {
				repo.EXPECT().Create(toBeInserted).
					Return(models.User{}, repo_errors.Duplicate("user_email")).
					Times(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Register",
		func(sCtx provider.StepCtx) {
			err = service.Register(user.Email, user.Password, user.Name)
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

func TestLoginServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(LoginServiceTestSuite))
}

