package category_test

import (
	// "errors"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	defcategory "rent_service/internal/logic/services/implementations/defservices/category"
	"rent_service/internal/logic/services/interfaces/category"
	"rent_service/internal/misc/types/collection"
	category_pmock "rent_service/internal/repository/context/mock/category"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	category_mock "rent_service/internal/repository/implementation/mock/repositories/category"
	"slices"

	// "slices"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CategoryServiceSuite struct {
	suite.Suite
}

func getService(f func(repository *category_mock.MockRepository)) category.IService {
	repo := category_mock.New()

	if nil != f {
		f(repo)
	}

	return defcategory.New(category_pmock.New(repo))
}

func (self *CategoryServiceSuite) BeforeEach(t provider.T) {
	t.AddParentSuite("DefServices")
	t.Epic("Default services implementation")
	t.Feature("Category service")
}

func (self *CategoryServiceSuite) describeList(t provider.T, title string, description string) {
	t.AddSubSuite("ListCategories")
	t.Story("List all")
	t.Title(title)
	t.Description(description)
}

func (self *CategoryServiceSuite) TestListCategoriesPositive(t provider.T) {
	var categories []models.Category
	var reference []category.Category
	var service category.IService

	self.describeList(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create reference categories", func(sCtx provider.StepCtx) {
			categories = models_om.CategoryToPath(models_om.CategoryDefaultPath()...)
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
			service = getService(func(builder *category_mock.MockRepository) {
				builder.WithGetAll(func() (collection.Collection[models.Category], error) {
					return collection.SliceCollection(categories), nil
				})
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

func (self *CategoryServiceSuite) TestListCategoriesInternalError(t provider.T) {
	var service category.IService

	self.describeList(t,
		"Error mappging",
		"Checks that any error is mapped to Interanal:DataAccess error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create empty service", func(sCtx provider.StepCtx) {
			service = getService(nil)
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

func (self *CategoryServiceSuite) describeFull(t provider.T, title string, description string) {
	t.AddSubSuite("GetFullCategory")
	t.Story("List path to leaf")
	t.Title(title)
	t.Description(description)
}

func generateDefaultPath(
	t provider.T,
	categories *[]models.Category,
	reference *[]category.Category,
	service *category.IService,
) {
	t.WithNewStep("Create reference categories", func(sCtx provider.StepCtx) {
		*categories = models_om.CategoryToPath(models_om.CategoryDefaultPath()...)
		*reference = make([]category.Category, len(*categories))
		for i, v := range *categories {
			(*reference)[i] = category.Category{
				Id:       v.Id,
				ParentId: v.ParentId,
				Name:     v.Name,
			}
		}
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		*service = getService(func(builder *category_mock.MockRepository) {
			builder.WithGetPath(func(leaf uuid.UUID) (collection.Collection[models.Category], error) {
				if i := slices.IndexFunc(*categories, func(c models.Category) bool {
					return leaf == c.Id
				}); 0 <= i {
					return collection.SliceCollection((*categories)[0 : i+1]), nil
				} else {
					return nil, repo_errors.NotFound("category_id")
				}
			})
		})
	})
}

func (self *CategoryServiceSuite) TestGetFullPathPositive(t provider.T) {
	var categories []models.Category
	var reference []category.Category
	var service category.IService

	self.describeFull(t,
		"Simple return predefined test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generateDefaultPath(t, &categories, &reference, &service)
	})

	// Act
	var result collection.Collection[category.Category]
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		result, err = service.GetFullCategory(categories[len(categories)-1].Id)
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *CategoryServiceSuite) TestGetFullPathNotFound(t provider.T) {
	var categories []models.Category
	var reference []category.Category
	var service category.IService
	var id uuid.UUID

	self.describeFull(t,
		"Category not found",
		"Not found must be returned if unknown id is specified",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generateDefaultPath(t, &categories, &reference, &service)

		t.WithNewStep("Generate unknown id", func(sCtx provider.StepCtx) {
			value, err := uuid.NewRandom()

			if nil != err {
				t.Fatalf("Unable to generate uuid: %v", err)
			}

			for slices.ContainsFunc(categories, func(c models.Category) bool {
				return c.Id == value
			}) {
				value, err = uuid.NewRandom()

				if nil != err {
					t.Fatalf("Unable to generate uuid: %v", err)
				}
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		_, err = service.GetFullCategory(id)
	})

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is not found")
}

func (self *CategoryServiceSuite) TestGetFullPathInternalError(t provider.T) {
	var service category.IService

	self.describeFull(t,
		"Error mappging",
		"Checks that error is mapped to Interanal:DataAccess error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create empty service", func(sCtx provider.StepCtx) {
			service = getService(nil)
		})
	})

	// Act
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		_, err = service.ListCategories()
	})

	// Assert
	var ierr cmnerrors.ErrorInternal
	var derr cmnerrors.ErrorDataAccess

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &ierr, "Error is internal")
	t.Assert().ErrorAs(ierr, &derr, "Error is data access")
}

func TestCategorySuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CategoryServiceSuite))
}

