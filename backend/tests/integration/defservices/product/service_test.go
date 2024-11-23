package product_test

import (
	"fmt"
	"math/rand/v2"
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/product"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"rent_service/misc/testcommon/defservices"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func MapProduct(value *models.Product) product.Product {
	return product.Product{
		Id:          value.Id,
		Name:        value.Name,
		CategoryId:  value.CategoryId,
		Description: value.Description,
	}
}

type ProductServiceIntegrationTestSuite struct {
	suite.Suite
	service  product.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *ProductServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *ProductServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *ProductServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Product service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateProductService()
	})
}

var describeListProducts = testcommon.MethodDescriptor(
	"ListProducts",
	"Get list of products with filter",
)

var describeGetProductById = testcommon.MethodDescriptor(
	"GetProductById",
	"Get product by identifier",
)

func (self *ProductServiceIntegrationTestSuite) TestListProductsPositive(t provider.T) {
	var (
		reference []product.Product
		filter    product.Filter
		sort      product.Sort
	)

	describeListProducts(t,
		"Simple return with filter",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewList(t, "products",
			func(i uint) (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple(
					fmt.Sprintf("%v", i), cgen.Generate(),
				).Build()
				reference = append(reference, MapProduct(&v))
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		opgen := psql.GeneratorStepNewList(t, "odd products",
			func(i uint) (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple(
					fmt.Sprintf("%v-odd", i), cgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(opgen)

		chgen := psql.GeneratorStep(t, "characteristics",
			func(spy *nullable.Nullable[generator.Spy]) generator.IGenerator {
				var chars []models.ProductCharacteristics
				var i uint = 0

				return generator.NewFunc(
					func() uuid.UUID {
						var v models.ProductCharacteristics
						if 5 > i {
							v = models_om.ProductCharacteristics(
								pgen.Generate(),
								models_om.CharacteristicExample(
									"key1",
									fmt.Sprintf("value%v", rand.Uint64()%3+1),
								).Build(),
								models_om.CharacteristicExampleNumeric(
									"key2",
									1+rand.Float64(),
								).Build(),
							).Build()
						} else {
							v = models_om.ProductCharacteristics(
								opgen.Generate(),
								models_om.CharacteristicExample(
									"key1",
									fmt.Sprintf("value%v", rand.Uint64()%3+4),
								).Build(),
								models_om.CharacteristicExampleNumeric(
									"key2",
									rand.Float64(),
								).Build(),
							).Build()
						}

						chars = append(chars, v)

						return v.ProductId
					},
					func() {
						nullable.IfSome(spy, func(spy *generator.Spy) {
							spy.SniffValue("characteristics", chars)
						})
						psql.BulkInsert(
							self.rContext.Inserter.InsertProductCharacteristics,
							chars...,
						)
					},
				)
			},
		)
		gg.Add(chgen, 10)

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			var query = "example"
			filter = testcommon.AssignParameter(sCtx, "filter",
				product.Filter{
					CategoryId: cgen.Generate(),
					Query:      &query,
					Characteristics: []product.FilterCharachteristic{
						{Key: "key1", Values: []string{"value1", "value2", "value3"}},
						{Key: "key2", Range: &struct {
							Min float64
							Max float64
						}{Min: 1, Max: 2}},
					},
				},
			)
			sort = testcommon.AssignParameter(sCtx, "sort", product.SORT_NONE)
		})

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[product.Product]
	var err error

	t.WithNewStep(
		"Get all products",
		func(sCtx provider.StepCtx) {
			result, err = self.service.ListProducts(filter, sort)
		},
		allure.NewParameter("filter", filter),
		allure.NewParameter("sort", sort),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *ProductServiceIntegrationTestSuite) TestListProductsEmptyCollection(t provider.T) {
	var (
		filter product.Filter
		sort   product.Sort
	)

	describeListProducts(t,
		"No products found",
		"Check that empty collection is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewList(t, "products",
			func(i uint) (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple(
					fmt.Sprintf("%v", i), cgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		chgen := psql.GeneratorStepNewList(t, "characteristics",
			func(i uint) (models.ProductCharacteristics, uuid.UUID) {
				v := models_om.ProductCharacteristics(
					pgen.Generate(),
					models_om.CharacteristicExample(
						"key1",
						fmt.Sprintf("value%v", rand.Uint64()%3+4),
					).Build(),
					models_om.CharacteristicExampleNumeric(
						"key2",
						rand.Float64(),
					).Build(),
				).Build()
				return v, v.ProductId
			},
			self.rContext.Inserter.InsertProductCharacteristics,
		)
		gg.Add(chgen, 5)

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			var query = "example"
			filter = testcommon.AssignParameter(sCtx, "filter",
				product.Filter{
					CategoryId: cgen.Generate(),
					Query:      &query,
					Characteristics: []product.FilterCharachteristic{
						{Key: "key1", Values: []string{"value1", "value2", "value3"}},
						{Key: "key2", Range: &struct {
							Min float64
							Max float64
						}{Min: 1, Max: 2}},
					},
				},
			)
			sort = testcommon.AssignParameter(sCtx, "sort", product.SORT_NONE)
		})

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[product.Product]
	var err error

	t.WithNewStep("Get all products", func(sCtx provider.StepCtx) {
		result, err = self.service.ListProducts(filter, sort)
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()),
		"Collection is empty")
}

func (self *ProductServiceIntegrationTestSuite) TestGetProductByIdPositive(t provider.T) {
	var reference product.Product

	describeGetProductById(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "products",
			func() (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				reference = MapProduct(&v)
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.Add(pgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result product.Product
	var err error

	t.WithNewStep(
		"Get product",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetProductById(reference.Id)
		},
		allure.NewParameter("productId", reference.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value matches reference")
}

func (self *ProductServiceIntegrationTestSuite) TestGetProductByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetProductById(t,
		"Product not found",
		"Checks that get method calls repository and NotFound error is retuned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create unknown id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id",
				uuidgen.Generate(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep(
		"Get product",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetProductById(id)
		},
		allure.NewParameter("productId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func MapCharacteristics(value *models.ProductCharacteristics) []product.Charachteristic {
	return collection.Collect(
		collection.MapIterator(
			func(kv *collection.KV[string, models.Charachteristic]) product.Charachteristic {
				return product.Charachteristic{
					Id:    kv.Value.Id,
					Name:  kv.Value.Name,
					Value: kv.Value.Value,
				}
			},
			collection.HashMapIterator(value.Map),
		),
	)
}

type ProductCharacteristicsServiceIntegrationTestSuite struct {
	suite.Suite
	service  product.ICharacteristicsService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *ProductCharacteristicsServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *ProductCharacteristicsServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *ProductCharacteristicsServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Product characteristics service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateProductCharacteristicsService()
	})
}

var describeListProductCharacteristics = testcommon.MethodDescriptor(
	"ListProductCharacteristics",
	"Get product characteristics",
)

func (self *ProductCharacteristicsServiceIntegrationTestSuite) TestListProductCharacteristicsPositive(t provider.T) {
	var (
		prod      models.Product
		reference []product.Charachteristic
	)

	describeListProductCharacteristics(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepValue(t, "product", &prod,
			func() (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		chgen := psql.GeneratorStepNewValue(t, "characteristics",
			func() (models.ProductCharacteristics, uuid.UUID) {
				v := models_om.ProductCharacteristics(
					pgen.Generate(),
					models_om.CharacteristicExample("key1", "value1").Build(),
					models_om.CharacteristicExample("key2", "value2").Build(),
					models_om.CharacteristicExampleNumeric("key3", 1.6).Build(),
					models_om.CharacteristicExampleNumeric("key4", -51.215).Build(),
				).Build()
				reference = MapCharacteristics(&v)
				return v, v.ProductId
			},
			self.rContext.Inserter.InsertProductCharacteristics,
		)
		gg.Add(chgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[product.Charachteristic]
	var err error

	t.WithNewStep("Get all product characteristics", func(sCtx provider.StepCtx) {
		result, err = self.service.ListProductCharacteristics(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *ProductCharacteristicsServiceIntegrationTestSuite) TestListProductCharacteristicsNotFound(t provider.T) {
	var id uuid.UUID

	describeListProductCharacteristics(t,
		"Product not found",
		"Check that error is retuned and mapped to NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create unknown id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id",
				uuidgen.Generate(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all product characteristics", func(sCtx provider.StepCtx) {
		_, err = self.service.ListProductCharacteristics(id)
	}, allure.NewParameter("productId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type ProductPhotoServiceIntegrationTestSuite struct {
	suite.Suite
	service  product.IPhotoService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *ProductPhotoServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *ProductPhotoServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *ProductPhotoServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Product photo service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateProductPhotoService()
	})
}

var describeListProductPhotos = testcommon.MethodDescriptor(
	"ListProductPhotos",
	"Get list of all product photos",
)

func (self *ProductPhotoServiceIntegrationTestSuite) TestListProductPhotosPositive(t provider.T) {
	var (
		prod      models.Product
		reference []uuid.UUID
	)

	describeListProductPhotos(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		prodgen := psql.GeneratorStepValue(t, "pick up point", &prod,
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple(
					"test",
					cgen.Generate(),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.Add(prodgen, 1)

		pgen := psql.GeneratorStepNewList(t, "photos",
			func(i uint) (models.Photo, uuid.UUID) {
				path := self.sContext.PhotoRegistry.SavePhoto(
					models_om.ImagePNGContent(nullable.None[int]()),
				)
				p := models_om.PhotoExample(
					fmt.Sprint(i),
					nullable.None[time.Time](),
				).WithPath(path).Build()
				reference = append(reference, p.Id)
				return p, p.Id
			},
			func(photo *models.Photo) {
				self.rContext.Inserter.InsertPhoto(photo)
				self.rContext.Inserter.InsertProductPhoto(psql.NewPhoto(
					nullable.None[uuid.UUID](),
					prod.Id,
					photo.Id,
				))
			},
		)
		gg.Add(pgen, 5)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[uuid.UUID]
	var err error

	t.WithNewStep("Get all product photos", func(sCtx provider.StepCtx) {
		result, err = self.service.ListProductPhotos(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *ProductPhotoServiceIntegrationTestSuite) TestListProductPhotosNotFound(t provider.T) {
	var id uuid.UUID

	describeListProductPhotos(t,
		"Product not found",
		"Check that erro is returned and is NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Generate unknown id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all product photos", func(sCtx provider.StepCtx) {
		_, err = self.service.ListProductPhotos(id)
	}, allure.NewParameter("productId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestProductServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductServiceIntegrationTestSuite))
}

func TestProductCharacteristicsServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductCharacteristicsServiceIntegrationTestSuite))
}

func TestProductPhotoServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductPhotoServiceIntegrationTestSuite))
}

