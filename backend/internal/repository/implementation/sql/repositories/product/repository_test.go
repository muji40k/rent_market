package product_test

import (
	"math/rand/v2"
	models_b "rent_service/builders/domain/models"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/product"
	"time"

	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	repo product.IRepository
	psqlcommon.Context
}

func (self *ProductRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *ProductRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *ProductRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Delivery Company repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateProductRepository()
	})
}

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get product by id",
)

var describeGetWithFilter = testcommon.MethodDescriptor(
	"GetWithFilter",
	"Get products with filter",
)

func (self *ProductRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var (
		category  models.Category
		reference models.Product
	)

	describeGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("1").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})
		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&reference)
		})
	})

	// Act
	var result models.Product
	var err error

	t.WithNewStep("Get product by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("productId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same product value")
}

func (self *ProductRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Product not found",
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

	t.WithNewStep("Get product by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("productId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *ProductRepositoryTestSuite) TestGetWithFilterPositive(t provider.T) {
	var (
		category  models.Category
		chars     []models.ProductCharacteristics
		all       []models.Product
		odd       []models.Product
		reference []models.Product
		filter    product.Filter
		sort      product.Sort
	)

	describeGetWithFilter(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("Test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert products", func(sCtx provider.StepCtx) {
			all = testcommon.AssignParameter(sCtx, "products",
				collect.DoN(10, collect.FmtWrap(func(p string) *models_b.ProductBuilder {
					return models_om.ProductExmaple(p, category.Id)
				})),
			)
			odd = all[:5]
			reference = all[5:]
			psql.BulkInsert(self.Inserter.InsertProduct, all...)
		})

		t.WithNewStep("Create and insert product characteristics", func(sCtx provider.StepCtx) {
			chars = testcommon.AssignParameter(sCtx, "characteristics",
				collection.Collect(
					collection.ChainIterator(
						collection.MapIterator(
							func(p *models.Product) models.ProductCharacteristics {
								return models_om.ProductCharacteristics(
									p.Id,
									models_om.CharacteristicExample("key1", "acceptable").Build(),
									models_om.CharacteristicExampleNumeric("key2", rand.Float64()).Build(),
								).Build()
							},
							collection.SliceIterator(reference),
						),
						collection.MapIterator(
							func(p *models.Product) models.ProductCharacteristics {
								return models_om.ProductCharacteristics(
									p.Id,
									models_om.CharacteristicExample("key1", "odd").Build(),
									models_om.CharacteristicExampleNumeric("key2", 1+3*rand.Float64()).Build(),
								).Build()
							},
							collection.SliceIterator(odd),
						),
					),
				),
			)
			psql.BulkInsert(self.Inserter.InsertProductCharacteristics, chars...)
		})

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			filter = product.Filter{
				CategoryId: category.Id,
				Query:      nil,
				Ranges: []product.Range{
					{
						Characteristic: product.Characteristic{
							Key: "key2",
						},
						Min: 0,
						Max: 1,
					},
				},
				Selectors: []product.Selector{
					{
						Characteristic: product.Characteristic{
							Key: "key1",
						},
						Values: []string{"acceptable"},
					},
				},
			}
			sort = product.SORT_NONE
		})
	})

	// Act
	var result collection.Collection[models.Product]
	var err error

	t.WithNewStep("Get products with filter", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetWithFilter(filter, sort)
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(reference, collection.Collect(result.Iter()),
		"Same company values")
}

func (self *ProductRepositoryTestSuite) TestGetWithFilterNotFound(t provider.T) {
	var (
		id     uuid.UUID
		filter product.Filter
		sort   product.Sort
	)

	describeGetWithFilter(t,
		"No category found",
		"Checks that method return empty collection withour error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknonw id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			filter = product.Filter{
				CategoryId: id,
				Query:      nil,
				Ranges:     nil,
				Selectors:  nil,
			}
			sort = product.SORT_NONE
		})
	})

	// Act
	var err error

	t.WithNewStep("Get products", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetWithFilter(filter, sort)
	})

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type ProductCharacteristicsRepositoryTestSuite struct {
	suite.Suite
	repo product.ICharacteristicsRepository
	psqlcommon.Context
}

func (self *ProductCharacteristicsRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *ProductCharacteristicsRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *ProductCharacteristicsRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Product Characteristics repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateProductCharacteristicsRepository()
	})
}

var describeCharacteristicsGetByProductId = testcommon.MethodDescriptor(
	"GetByProductId",
	"Get product characteristics by id",
)

func (self *ProductCharacteristicsRepositoryTestSuite) TestGetByProductIdPositive(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		reference models.ProductCharacteristics
	)

	describeCharacteristicsGetByProductId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("1").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert characteristics", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "characteristeics",
				models_om.ProductCharacteristics(
					product.Id,
					collect.Do(
						models_om.CharacteristicExample("key1", "value1"),
						models_om.CharacteristicExample("key2", "value2"),
						models_om.CharacteristicExample("key3", "value3"),
						models_om.CharacteristicExampleNumeric("key4", rand.Float64()),
						models_om.CharacteristicExampleNumeric("key5", rand.Float64()),
						models_om.CharacteristicExampleNumeric("key6", rand.Float64()),
					)...,
				).Build(),
			)
			self.Inserter.InsertProductCharacteristics(&reference)
		})
	})

	// Act
	var result models.ProductCharacteristics
	var err error

	t.WithNewStep("Get product charcteristics", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetByProductId(product.Id)
	}, allure.NewParameter("productId", product.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same delivery value")
}

func (self *ProductCharacteristicsRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeCharacteristicsGetByProductId(t,
		"Product not found",
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

	t.WithNewStep("Get product characteristics", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByProductId(id)
	}, allure.NewParameter("productId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type ProductPhotoRepositoryTestSuite struct {
	suite.Suite
	repo product.IPhotoRepository
	psqlcommon.Context
}

func (self *ProductPhotoRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *ProductPhotoRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *ProductPhotoRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Product photo repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateProductPhotoRepository()
	})
}

var describePhotoGetByProductId = testcommon.MethodDescriptor(
	"GetByProductId",
	"Get product photos by id",
)

func (self *ProductPhotoRepositoryTestSuite) TestGetByProductIdPositive(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		reference []models.Photo
	)

	describePhotoGetByProductId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert photos", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photos",
				collect.DoN(5, collect.FmtWrap(func(p string) *models_b.PhotoBuilder {
					return models_om.PhotoExample(p, nullable.None[time.Time]())
				})),
			)
			psql.BulkInsert(self.Inserter.InsertPhoto, reference...)
		})

		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("1").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert product photos", func(sCtx provider.StepCtx) {
			psql.BulkInsert(self.Inserter.InsertProductPhoto,
				collection.Collect(
					collection.MapIterator(
						func(photo *models.Photo) psql.Photo {
							return *psql.NewPhoto(
								nullable.None[uuid.UUID](),
								product.Id,
								photo.Id,
							)
						},
						collection.SliceIterator(reference),
					),
				)...,
			)
		})
	})

	// Act
	var result collection.Collection[uuid.UUID]
	var err error

	t.WithNewStep("Get product photos by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetByProductId(product.Id)
	}, allure.NewParameter("productId", product.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(
		collection.Collect(
			collection.MapIterator(
				func(photo *models.Photo) uuid.UUID {
					return photo.Id
				},
				collection.SliceIterator(reference),
			),
		),
		collection.Collect(result.Iter()),
		"Same photo ids",
	)
}

func (self *ProductPhotoRepositoryTestSuite) TestGetByProductIdNotFound(t provider.T) {
	var id uuid.UUID

	describePhotoGetByProductId(t,
		"Product not found",
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

	t.WithNewStep("Get product photos by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetByProductId(id)
	}, allure.NewParameter("productId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestProductRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ProductRepositoryTestSuite))
}

func TestProductCharacteristicsRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ProductCharacteristicsRepositoryTestSuite))
}

func TestProductPhotoRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ProductPhotoRepositoryTestSuite))
}

