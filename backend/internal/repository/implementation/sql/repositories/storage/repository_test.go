package storage_test

import (
	"fmt"
	models_b "rent_service/builders/domain/models"
	records_b "rent_service/builders/domain/records"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	records_om "rent_service/builders/mothers/domain/records"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"

	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/storage"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	psqlcommon "rent_service/misc/testcommon/psql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func CompareStorage(ref records.Storage, act records.Storage) bool {
	return ref.Id == act.Id &&
		ref.PickUpPointId == act.PickUpPointId &&
		ref.InstanceId == act.InstanceId &&
		psqlcommon.CompareTimeMicro(ref.InDate, act.InDate) &&
		psqlcommon.CompareTimePtrMicro(ref.OutDate, act.OutDate)
}

type StorageRepositoryTestSuite struct {
	suite.Suite
	repo storage.IRepository
	psqlcommon.Context
}

func (self *StorageRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *StorageRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *StorageRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Storage repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateStorageRepository()
	})
}

var describeCreate = testcommon.MethodDescriptor(
	"Create",
	"Create sorage",
)

var describeUpdate = testcommon.MethodDescriptor(
	"Update",
	"Update sorage",
)

var describeGetActiveByPickUpPointId = testcommon.MethodDescriptor(
	"GetActiveByPickUpPointId",
	"Get active storages by pick up point",
)

var describeGetActiveByInstanceId = testcommon.MethodDescriptor(
	"GetActiveByInstanceId",
	"Get active storages by instance",
)

func CheckCreated(ref records.Storage, act records.Storage) bool {
	return uuid.UUID{} != act.Id &&
		ref.PickUpPointId == act.PickUpPointId &&
		ref.InstanceId == act.InstanceId &&
		psqlcommon.CompareTimeMicro(ref.InDate, act.InDate) &&
		psqlcommon.CompareTimePtrMicro(ref.OutDate, act.OutDate)
}

func (self *StorageRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		instance  models.Instance
		pup       models.PickUpPoint
		reference records.Storage
	)

	describeCreate(t,
		"Simple create test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storage",
				records_om.StorageActive(
					pup.Id,
					instance.Id,
					nullable.None[time.Time](),
				).WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result records.Storage
	var err error

	t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(reference)
	}, allure.NewParameter("storage", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[records.Storage](t).EqualFunc(
		CheckCreated, reference, result, "Same storage with non null uuid",
	)
}

func (self *StorageRepositoryTestSuite) TestCreateNotFound(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		instance  models.Instance
		pup       models.PickUpPoint
		reference records.Storage
	)

	describeCreate(t,
		"Instance not exists",
		"Checks that error returned and mapped to NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and DON'T insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
		})

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storage",
				records_om.StorageActive(
					pup.Id,
					instance.Id,
					nullable.None[time.Time](),
				).WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(reference)
	}, allure.NewParameter("storage", reference))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *StorageRepositoryTestSuite) TestUpdatePositive(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		instance  models.Instance
		pup       models.PickUpPoint
		created   records.Storage
		reference records.Storage
	)

	describeUpdate(t,
		"Simple update test (storage finished)",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create and insert storage", func(sCtx provider.StepCtx) {
			builder := records_om.StorageActive(
				pup.Id,
				instance.Id,
				nullable.None[time.Time](),
			)
			created = testcommon.AssignParameter(sCtx, "created",
				builder.Build(),
			)
			reference = testcommon.AssignParameter(sCtx, "storage",
				builder.WithOutDate(nullable.Some(time.Now())).Build(),
			)
			self.Inserter.InsertStorage(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update storage", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("storage", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *StorageRepositoryTestSuite) TestUpdateNotFound(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		instance  models.Instance
		pup       models.PickUpPoint
		reference records.Storage
	)

	describeUpdate(t,
		"Storage not exists",
		"Checks that error returned and mapped to NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert category", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("test").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert product", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and DON'T insert instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
		})

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storage",
				records_om.StorageFinished(
					pup.Id,
					instance.Id,
					nullable.None[time.Time](),
					nullable.None[time.Time](),
				).Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update storage", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("storage", reference))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *StorageRepositoryTestSuite) TestGetActiveByPickUpPointIdPositive(t provider.T) {
	var (
		categories []models.Category
		products   []models.Product
		instances  []models.Instance
		pup        models.PickUpPoint
		reference  []records.Storage
	)

	describeGetActiveByPickUpPointId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				collect.DoN(5, collect.FmtWrap(func(p string) *models_b.CategoryBuilder {
					return models_om.CategoryRandomId().WithName(p)
				})),
			)
			psql.BulkInsert(self.Inserter.InsertCategory, categories...)
		})

		t.WithNewStep("Create and insert products", func(sCtx provider.StepCtx) {
			products = testcommon.AssignParameter(sCtx, "products",
				collect.DoN(5, func(i uint) *models_b.ProductBuilder {
					return models_om.ProductExmaple(fmt.Sprint(i), categories[i].Id)
				}),
			)
			psql.BulkInsert(self.Inserter.InsertProduct, products...)
		})

		t.WithNewStep("Create and insert instances", func(sCtx provider.StepCtx) {
			instances = testcommon.AssignParameter(sCtx, "instances",
				collect.DoN(5, func(i uint) *models_b.InstanceBuilder {
					return models_om.InstanceExample(fmt.Sprint(i), products[i].Id)
				}),
			)
			psql.BulkInsert(self.Inserter.InsertInstance, instances...)
		})

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create storages", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storages",
				collect.DoN(5, func(i uint) *records_b.StorageBuilder {
					return records_om.StorageActive(
						pup.Id,
						instances[i].Id,
						nullable.None[time.Time](),
					)
				}),
			)
			psql.BulkInsert(self.Inserter.InsertStorage, reference...)
		})
	})

	// Act
	var result collection.Collection[records.Storage]
	var err error

	t.WithNewStep("Get active storages by pick up point id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetActiveByPickUpPointId(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[records.Storage](t).ElementsMatchFunc(
		CompareStorage, reference, collection.Collect(result.Iter()),
		"Same storages",
	)
}

func (self *StorageRepositoryTestSuite) TestGetActiveByPickUpPointIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetActiveByPickUpPointId(t,
		"PickUpPoint not found",
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

	t.WithNewStep("Get active storages by pick up point id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetActiveByPickUpPointId(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *StorageRepositoryTestSuite) TestGetActiveByInstanceIdPositive(t provider.T) {
	var (
		category  models.Category
		product   models.Product
		instance  models.Instance
		pup       models.PickUpPoint
		reference records.Storage
	)

	describeGetActiveByInstanceId(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			category = testcommon.AssignParameter(sCtx, "category",
				models_om.CategoryRandomId().WithName("1").Build(),
			)
			self.Inserter.InsertCategory(&category)
		})

		t.WithNewStep("Create and insert products", func(sCtx provider.StepCtx) {
			product = testcommon.AssignParameter(sCtx, "product",
				models_om.ProductExmaple("1", category.Id).Build(),
			)
			self.Inserter.InsertProduct(&product)
		})

		t.WithNewStep("Create and insert instances", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("1", product.Id).Build(),
			)
			self.Inserter.InsertInstance(&instance)
		})

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create storages", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "storage",
				records_om.StorageActive(
					pup.Id,
					instance.Id,
					nullable.None[time.Time](),
				).Build(),
			)
			self.Inserter.InsertStorage(&reference)
		})
	})

	// Act
	var result records.Storage
	var err error

	t.WithNewStep("Get active storages by instance id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetActiveByInstanceId(instance.Id)
	}, allure.NewParameter("instanceId", instance.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[records.Storage](t).EqualFunc(
		CompareStorage, reference, result, "Same storage",
	)
}

func (self *StorageRepositoryTestSuite) TestGetActiveByInstanceIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetActiveByInstanceId(t,
		"Instance not found",
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

	t.WithNewStep("Get active storages by instance id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetActiveByInstanceId(id)
	}, allure.NewParameter("instanceId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestStorageRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(StorageRepositoryTestSuite))
}

