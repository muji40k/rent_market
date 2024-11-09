package category_test

import (
	"errors"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	defcategory "rent_service/internal/logic/services/implementations/defservices/services/category"
	"rent_service/internal/logic/services/interfaces/category"
	"rent_service/internal/misc/types/collection"
	category_pmock "rent_service/internal/repository/context/mock/category"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	category_mock "rent_service/internal/repository/implementation/mock/category"
	"rent_service/misc/testcommon"

	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

func GetService(
	ctrl *gomock.Controller,
	f func(repository *category_mock.MockIRepository),
) category.IService {
	repo := category_mock.NewMockIRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return defcategory.New(category_pmock.New(repo))
}

func generateDefaultPath(
	t provider.T,
	ctrl *gomock.Controller,
) (category.IService, []models.Category, []category.Category) {
	var service category.IService
	var categories []models.Category
	var reference []category.Category

	t.WithNewStep("Create reference categories", func(sCtx provider.StepCtx) {
		categories = testcommon.AssignParameter(sCtx, "categories",
			models_om.CategoryToPath(models_om.CategoryDefaultPath()...),
		)
		reference = make([]category.Category, len(categories))
		for i, v := range categories {
			reference[i] = category.Category{
				Id:       v.Id,
				ParentId: v.ParentId,
				Name:     v.Name,
			}
		}
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		service = GetService(ctrl, func(repo *category_mock.MockIRepository) {
			repo.EXPECT().GetPath(gomock.Any()).
				DoAndReturn(func(leaf uuid.UUID) (collection.Collection[models.Category], error) {
					if i := slices.IndexFunc(categories, func(c models.Category) bool {
						return leaf == c.Id
					}); 0 <= i {
						return collection.SliceCollection(categories[0 : i+1]), nil
					} else {
						return nil, repo_errors.NotFound("category_id")
					}
				}).
				MinTimes(1)
		})
	})

	return service, categories, reference
}

type CategoryServiceTestSuite struct {
	suite.Suite
}

func (self *CategoryServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Category service",
	)
}

var describeListCategories = testcommon.MethodDescriptor(
	"ListCategories",
	"List All",
)
var describeGetFullCategory = testcommon.MethodDescriptor(
	"GetFullCategory",
	"List path to leaf",
)

func (self *CategoryServiceTestSuite) TestListCategoriesPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var categories []models.Category
	var reference []category.Category
	var service category.IService

	describeListCategories(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create reference categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				models_om.CategoryToPath(models_om.CategoryDefaultPath()...),
			)
			reference = make([]category.Category, len(categories))
			for i, v := range categories {
				reference[i] = category.Category{
					Id:       v.Id,
					ParentId: v.ParentId,
					Name:     v.Name,
				}
			}
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *category_mock.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(collection.SliceCollection(categories), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[category.Category]
	var err error

	t.WithNewStep("Get all categories", func(sCtx provider.StepCtx) {
		result, err = service.ListCategories()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *CategoryServiceTestSuite) TestListCategoriesInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var service category.IService

	describeListCategories(t,
		"Error mappging",
		"Checks that any error is mapped to Interanal:DataAccess error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create empty service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *category_mock.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(nil, errors.New("Some Internal Error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all categories", func(sCtx provider.StepCtx) {
		_, err = service.ListCategories()
	})

	// Assert
	var ierr cmnerrors.ErrorInternal
	var derr cmnerrors.ErrorDataAccess

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &ierr, "Error is internal")
	t.Assert().ErrorAs(ierr, &derr, "Error is data access")
}

func (self *CategoryServiceTestSuite) TestGetFullCategoryPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var categories []models.Category
	var reference []category.Category
	var service category.IService

	describeGetFullCategory(t,
		"Simple return predefined test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		service, categories, reference = generateDefaultPath(t, ctrl)
	})

	// Act
	var result collection.Collection[category.Category]
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		result, err = service.GetFullCategory(categories[len(categories)-1].Id)
	}, allure.NewParameter("leaf", categories[len(categories)-1].Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *CategoryServiceTestSuite) TestGetFullCategoryNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var categories []models.Category
	var service category.IService
	var id uuid.UUID

	describeGetFullCategory(t,
		"Category not found",
		"Not found must be returned if unknown id is specified",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		service, categories, _ = generateDefaultPath(t, ctrl)

		t.WithNewStep("Generate unknown id", func(sCtx provider.StepCtx) {
			var err error
			id, err = uuid.NewRandom()
			sCtx.Require().Nil(err, "Unable to generate uuid")

			for slices.ContainsFunc(categories, func(c models.Category) bool {
				return c.Id == id
			}) {
				id, err = uuid.NewRandom()
				sCtx.Require().Nil(err, "Unable to generate uuid")
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		_, err = service.GetFullCategory(id)
	}, allure.NewParameter("leaf", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is not found")
}

func (self *CategoryServiceTestSuite) TestGetFullCategoryInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service category.IService
	var id uuid.UUID

	describeGetFullCategory(t,
		"Error mappging",
		"Checks that error is mapped to Interanal:DataAccess error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate some uuid", func(sCtx provider.StepCtx) {
			var err error
			id, err = uuid.NewRandom()
			sCtx.Require().Nil(err)
		})

		t.WithNewStep("Create empty service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *category_mock.MockIRepository) {
				repo.EXPECT().GetPath(gomock.Any()).
					Return(nil, errors.New("Some Internal Error")).
					MinTimes(1)
			})

		})
	})

	// Act
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		_, err = service.GetFullCategory(id)
	}, allure.NewParameter("leaf", id))

	// Assert
	var ierr cmnerrors.ErrorInternal
	var derr cmnerrors.ErrorDataAccess

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &ierr, "Error is internal")
	t.Assert().ErrorAs(ierr, &derr, "Error is data access")
}

func TestCategoryServiceTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CategoryServiceTestSuite))
}

