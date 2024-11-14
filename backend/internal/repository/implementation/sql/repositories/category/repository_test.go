package category_test

import (
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/category"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	repo category.IRepository
	psqlcommon.Context
}

func (self *CategoryRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *CategoryRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *CategoryRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Category repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateCategoryRepository()
	})
}

var describeGetAll = testcommon.MethodDescriptor(
	"GetAll",
	"Get all categories",
)

var describeGetPath = testcommon.MethodDescriptor(
	"GetPath",
	"Get path to category by id",
)

func (self *CategoryRepositoryTestSuite) TestGetAllPositive(t provider.T) {
	var reference []models.Category

	describeGetAll(t,
		"Simple return all test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert reference categories", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(models_om.CategoryDefaultPath()...),
			)
			psql.BulkInsert(self.Inserter.InsertCategory, reference...)
		})
	})

	// Act
	var result collection.Collection[models.Category]
	var err error

	t.WithNewStep("Get all", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, collection.Collect(result.Iter()),
		"Same values")
}

func (self *CategoryRepositoryTestSuite) TestGetAllEmpty(t provider.T) {
	describeGetAll(t,
		"No values to return",
		"Checks that method return empty collection without error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[models.Category]
	var err error

	t.WithNewStep("Get all addresses", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()), "Collection is empty")
}

func (self *CategoryRepositoryTestSuite) TestGetPathPositive(t provider.T) {
	var reference []models.Category
	var last *models.Category

	describeGetPath(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert reference categories", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(models_om.CategoryDefaultPath()...),
			)
			last = &reference[len(reference)-1]
			psql.BulkInsert(self.Inserter.InsertCategory, reference...)
		})
	})

	// Act
	var result collection.Collection[models.Category]
	var err error

	t.WithNewStep("Get path to leaf", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetPath(last.Id)
	}, allure.NewParameter("leafId", last.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, collection.Collect(result.Iter()),
		"Same values")
}

func (self *CategoryRepositoryTestSuite) TestGetPathNotFound(t provider.T) {
	var id uuid.UUID

	describeGetPath(t,
		"Category not found",
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

	t.WithNewStep("Get all addresses", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetPath(id)
	}, allure.NewParameter("leafId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestCategoryRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CategoryRepositoryTestSuite))
}

