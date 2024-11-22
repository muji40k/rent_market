package category_test

import (
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/category"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/testcommon"
	"rent_service/misc/testcommon/defservices"
	psqlcommon "rent_service/misc/testcommon/psql"

	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func MapCategory(value *models.Category) category.Category {
	return category.Category{
		Id:       value.Id,
		ParentId: value.ParentId,
		Name:     value.Name,
	}
}

type CategoryServiceIntegrationTestSuite struct {
	suite.Suite
	service  category.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *CategoryServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *CategoryServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *CategoryServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Category service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateCategoryService()
	})
}

func (self *CategoryServiceIntegrationTestSuite) AfterEach(t provider.T) {
	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})
}

var describeListCategories = testcommon.MethodDescriptor(
	"ListCategories",
	"List All",
)
var describeGetFullCategory = testcommon.MethodDescriptor(
	"GetFullCategory",
	"List path to leaf",
)

func (self *CategoryServiceIntegrationTestSuite) TestListCategoriesPositive(t provider.T) {
	var (
		categories []models.Category
		reference  []category.Category
	)

	describeListCategories(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(models_om.CategoryDefaultPath()...),
			)
			reference = collection.Collect(collection.MapIterator(
				MapCategory, collection.SliceIterator(categories),
			))
			psql.BulkInsert(self.rContext.Inserter.InsertCategory, categories...)
		})
	})

	// Act
	var result collection.Collection[category.Category]
	var err error

	t.WithNewStep("Get all categories", func(sCtx provider.StepCtx) {
		result, err = self.service.ListCategories()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *CategoryServiceIntegrationTestSuite) TestGetFullCategoryPositive(t provider.T) {
	var (
		categories []models.Category
		reference  []category.Category
	)

	describeGetFullCategory(t,
		"Simple return predefined test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(models_om.CategoryDefaultPath()...),
			)
			reference = collection.Collect(collection.MapIterator(
				MapCategory, collection.SliceIterator(categories),
			))
			psql.BulkInsert(self.rContext.Inserter.InsertCategory, categories...)
		})
	})

	// Act
	var result collection.Collection[category.Category]
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		result, err = self.service.GetFullCategory(categories[len(categories)-1].Id)
	}, allure.NewParameter("leaf", categories[len(categories)-1].Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *CategoryServiceIntegrationTestSuite) TestGetFullCategoryNotFound(t provider.T) {
	var (
		id         uuid.UUID
		categories []models.Category
	)

	describeGetFullCategory(t,
		"Category not found",
		"Not found must be returned if unknown id is specified",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(models_om.CategoryDefaultPath()...),
			)
			psql.BulkInsert(self.rContext.Inserter.InsertCategory, categories...)
		})

		t.WithNewStep("Generate unknown id", func(sCtx provider.StepCtx) {
			id = uuidgen.Generate()

			for slices.ContainsFunc(categories, func(c models.Category) bool {
				return c.Id == id
			}) {
				id = uuidgen.Generate()
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		_, err = self.service.GetFullCategory(id)
	}, allure.NewParameter("leaf", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is not found")
}

func TestCategoryServiceIntegrationTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CategoryServiceIntegrationTestSuite))
}

