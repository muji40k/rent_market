package product_test

import (
	"errors"
	"reflect"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/product"
	"rent_service/internal/logic/services/interfaces/product"
	"rent_service/internal/misc/types/collection"

	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"rent_service/misc/testcommon"
	"testing"

	rproduct "rent_service/internal/repository/implementation/mock/product"
	repo "rent_service/internal/repository/interfaces/product"

	product_pmock "rent_service/internal/repository/context/mock/product"

	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

func GetService(ctrl *gomock.Controller, f func(repo *rproduct.MockIRepository)) product.IService {
	repo := rproduct.NewMockIRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.New(product_pmock.New(repo))
}

func MapProduct(value *models.Product) product.Product {
	return product.Product{
		Id:          value.Id,
		Name:        value.Name,
		CategoryId:  value.CategoryId,
		Description: value.Description,
	}
}

func MapFilter(value *product.Filter) repo.Filter {
	return repo.Filter{
		CategoryId: value.CategoryId,
		Query:      value.Query,
		Ranges: collection.Collect(
			collection.MapIterator(
				func(f *product.FilterCharachteristic) repo.Range {
					return repo.Range{
						Characteristic: repo.Characteristic{
							Key: f.Key,
						},
						Min: f.Range.Min,
						Max: f.Range.Max,
					}
				},
				collection.FilterIterator(
					func(f *product.FilterCharachteristic) bool {
						return nil != f.Range && nil == f.Values
					},
					collection.SliceIterator(value.Characteristics),
				),
			),
		),
		Selectors: collection.Collect(
			collection.MapIterator(
				func(f *product.FilterCharachteristic) repo.Selector {
					return repo.Selector{
						Characteristic: repo.Characteristic{
							Key: f.Key,
						},
						Values: f.Values,
					}
				},
				collection.FilterIterator(
					func(f *product.FilterCharachteristic) bool {
						return nil == f.Range && nil != f.Values
					},
					collection.SliceIterator(value.Characteristics),
				),
			),
		),
	}
}

func MapSort(value *product.Sort) repo.Sort {
	switch *value {
	case product.SORT_NONE:
		return repo.SORT_NONE
	case product.SORT_OFFERS_ASC:
		return repo.SORT_OFFERS_ASC
	case product.SORT_OFFERS_DSC:
		return repo.SORT_OFFERS_DSC
	}

	panic("Unknown sort value")
}

type FilterMatcher struct {
	value *repo.Filter
}

func (self FilterMatcher) Matches(x any) bool {
	if reflect.TypeOf(repo.Filter{}) != reflect.TypeOf(x) {
		return false
	}

	xc := reflect.ValueOf(x).Interface().(repo.Filter)

	return gomock.Eq(self.value.CategoryId).Matches(xc.CategoryId) &&
		gomock.Eq(self.value.Query).Matches(xc.Query) &&
		gomock.InAnyOrder(self.value.Ranges).Matches(xc.Ranges) &&
		gomock.InAnyOrder(self.value.Selectors).Matches(xc.Selectors)
}

func (self FilterMatcher) String() string {
	return "Filter elements matches"
}

type ProductServiceTestSuite struct {
	suite.Suite
}

func (self *ProductServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Product service",
	)
}

var describeListProducts = testcommon.MethodDescriptor(
	"ListProducts",
	"Get list of products with filter",
)

var describeGetProductById = testcommon.MethodDescriptor(
	"GetProductById",
	"Get product by identifier",
)

func (self *ProductServiceTestSuite) TestListPicUpPointsPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IService

	var (
		category  uuid.UUID
		products  []models.Product
		reference []product.Product
		filter    product.Filter
		sort      product.Sort
	)

	describeListProducts(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create category id", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				uuidgen.Generate(),
			)
		})

		t.WithNewStep("Create filter", func(sCtx provider.StepCtx) {
			var query = "example"
			filter = testcommon.AssignParameter(sCtx, "filter",
				product.Filter{
					CategoryId: category,
					Query:      &query,
					Characteristics: []product.FilterCharachteristic{
						{Key: "key1", Values: []string{"value1", "value2", "value3"}},
						{Key: "key2", Range: &struct {
							Min float64
							Max float64
						}{Min: 1, Max: 2}},
						{Key: "key3", Range: &struct {
							Min float64
							Max float64
						}{Min: -4, Max: 23}},
						{Key: "key4", Values: []string{"value1", "value2", "value3"}},
					},
				},
			)
			sort = testcommon.AssignParameter(sCtx, "sort", product.SORT_OFFERS_ASC)
		})

		t.WithNewStep("Create reference products", func(sCtx provider.StepCtx) {
			products = testcommon.AssignParameter(sCtx, "products",
				collect.Do(
					models_om.ProductExmaple("1", category),
					models_om.ProductExmaple("2", category),
					models_om.ProductExmaple("3", category),
					models_om.ProductExmaple("4", category),
					models_om.ProductExmaple("5", category),
				),
			)

			reference = collection.Collect(
				collection.MapIterator(
					MapProduct, collection.SliceIterator(products),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			f := MapFilter(&filter)
			s := MapSort(&sort)
			service = GetService(ctrl, func(repo *rproduct.MockIRepository) {
				repo.EXPECT().GetWithFilter(FilterMatcher{&f}, s).
					Return(collection.SliceCollection(products), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[product.Product]
	var err error

	t.WithNewStep(
		"Get all products",
		func(sCtx provider.StepCtx) {
			result, err = service.ListProducts(filter, sort)
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

func (self *ProductServiceTestSuite) TestListPickUpPointsInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IService

	var (
		filter product.Filter
		sort   product.Sort
	)

	describeListProducts(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rproduct.MockIRepository) {
				repo.EXPECT().GetWithFilter(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all products", func(sCtx provider.StepCtx) {
		_, err = service.ListProducts(filter, sort)
	})

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func (self *ProductServiceTestSuite) TestGetProductByIdPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IService

	var (
		category  uuid.UUID
		prod      models.Product
		reference product.Product
	)

	describeGetProductById(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create category id", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				uuidgen.Generate(),
			)
		})

		t.WithNewStep("Create reference period", func(sCtx provider.StepCtx) {
			prod = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category).Build(),
			)
			reference = MapProduct(&prod)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rproduct.MockIRepository) {
				repo.EXPECT().GetById(prod.Id).
					Return(prod, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result product.Product
	var err error

	t.WithNewStep(
		"Get product",
		func(sCtx provider.StepCtx) {
			result, err = service.GetProductById(prod.Id)
		},
		allure.NewParameter("productId", prod.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value matches reference")
}

func (self *ProductServiceTestSuite) TestGetProductByIdNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IService

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

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rproduct.MockIRepository) {
				repo.EXPECT().GetById(id).
					Return(models.Product{}, repo_errors.NotFound("product_id")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep(
		"Get product",
		func(sCtx provider.StepCtx) {
			_, err = service.GetProductById(id)
		},
		allure.NewParameter("productId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func GetCharacteristicsService(ctrl *gomock.Controller, f func(repo *rproduct.MockICharacteristicsRepository)) product.ICharacteristicsService {
	repo := rproduct.NewMockICharacteristicsRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.NewCharacteristics(product_pmock.NewCharacteristics(repo))
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

type ProductCharacteristicsServiceTestSuite struct {
	suite.Suite
}

func (self *ProductCharacteristicsServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Product characteristics service",
	)
}

var describeListProductCharacteristics = testcommon.MethodDescriptor(
	"ListProductCharacteristics",
	"Get product characteristics",
)

func (self *ProductCharacteristicsServiceTestSuite) TestListProductCharacteristicsPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.ICharacteristicsService

	var (
		prod      models.Product
		chars     models.ProductCharacteristics
		reference []product.Charachteristic
	)

	describeListProductCharacteristics(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create product", func(sCtx provider.StepCtx) {
			prod = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductRandomId().Build(),
			)
		})

		t.WithNewStep("Create reference characteristics", func(sCtx provider.StepCtx) {
			chars = testcommon.AssignParameter(sCtx, "Working hours",
				models_om.ProductCharacteristics(
					prod.Id, collect.Do(
						models_om.CharacteristicExample("key1", "value1"),
						models_om.CharacteristicExample("key2", "value2"),
						models_om.CharacteristicExampleNumeric("key3", 1.6),
						models_om.CharacteristicExampleNumeric("key4", -51.215),
					)...,
				).Build(),
			)
			reference = MapCharacteristics(&chars)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetCharacteristicsService(ctrl, func(repo *rproduct.MockICharacteristicsRepository) {
				repo.EXPECT().GetByProductId(prod.Id).
					Return(chars, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[product.Charachteristic]
	var err error

	t.WithNewStep("Get all product characteristics", func(sCtx provider.StepCtx) {
		result, err = service.ListProductCharacteristics(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *ProductCharacteristicsServiceTestSuite) TestListProductCharacteristicsInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.ICharacteristicsService

	var prod models.Product

	describeListProductCharacteristics(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create product", func(sCtx provider.StepCtx) {
			prod = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductRandomId().Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetCharacteristicsService(ctrl, func(repo *rproduct.MockICharacteristicsRepository) {
				repo.EXPECT().GetByProductId(prod.Id).
					Return(models.ProductCharacteristics{}, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all product characteristics", func(sCtx provider.StepCtx) {
		_, err = service.ListProductCharacteristics(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func GetPhotoService(ctrl *gomock.Controller, f func(repo *rproduct.MockIPhotoRepository)) product.IPhotoService {
	repo := rproduct.NewMockIPhotoRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.NewPhoto(product_pmock.NewPhoto(repo))
}

type ProductPhotoServiceTestSuite struct {
	suite.Suite
}

func (self *ProductPhotoServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Product photo service",
	)
}

var describeListProductPhotos = testcommon.MethodDescriptor(
	"ListProductPhotos",
	"Get list of all product photos",
)

func (self *ProductPhotoServiceTestSuite) TestListProductPhotosPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IPhotoService

	var (
		prod models.Product
		ids  []uuid.UUID
	)

	describeListProductPhotos(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create product", func(sCtx provider.StepCtx) {
			prod = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductRandomId().Build(),
			)
		})

		t.WithNewStep("Create reference photo ids", func(sCtx provider.StepCtx) {
			ids = testcommon.AssignParameter(sCtx, "photoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID {
							return uuidgen.Generate()
						},
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetPhotoService(ctrl, func(repo *rproduct.MockIPhotoRepository) {
				repo.EXPECT().GetByProductId(prod.Id).
					Return(collection.SliceCollection(ids), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[uuid.UUID]
	var err error

	t.WithNewStep("Get all product photos", func(sCtx provider.StepCtx) {
		result, err = service.ListProductPhotos(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(ids, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *ProductPhotoServiceTestSuite) TestListProductPhotosInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service product.IPhotoService

	var prod models.Product

	describeListProductPhotos(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create product", func(sCtx provider.StepCtx) {
			prod = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductRandomId().Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetPhotoService(ctrl, func(repo *rproduct.MockIPhotoRepository) {
				repo.EXPECT().GetByProductId(prod.Id).
					Return(nil, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all product photos", func(sCtx provider.StepCtx) {
		_, err = service.ListProductPhotos(prod.Id)
	}, allure.NewParameter("productId", prod.Id))

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductServiceTestSuite))
}

func TestProductCharacteristicsServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductCharacteristicsServiceTestSuite))
}

func TestProductPhotoServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProductPhotoServiceTestSuite))
}

