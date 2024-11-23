package storage_test

import (
	"fmt"
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	records_om "rent_service/builders/mothers/domain/records"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/storage"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
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

func MapStorage(value *records.Storage) storage.Storage {
	return storage.Storage{
		Id:            value.Id,
		PickUpPointId: value.PickUpPointId,
		InstanceId:    value.InstanceId,
		InDate:        date.New(value.InDate),
		OutDate: nullcommon.CopyPtrIfSome(nullable.Map(
			nullable.FromPtr(value.OutDate),
			func(t *time.Time) date.Date {
				return date.New(*t)
			},
		)),
	}
}

func CompareStorage(e storage.Storage, a storage.Storage) bool {
	return e.Id == a.Id &&
		e.PickUpPointId == a.PickUpPointId &&
		e.InstanceId == a.InstanceId &&
		psqlcommon.CompareTimeMicro(e.InDate.Time, a.InDate.Time) &&
		psqlcommon.CompareTimePtrMicro(
			psqlcommon.UnwrapDate(e.OutDate),
			psqlcommon.UnwrapDate(a.OutDate),
		)
}

type StorageServiceIntegrationTestSuite struct {
	suite.Suite
	service  storage.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *StorageServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *StorageServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *StorageServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Storage service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateStorageService()
	})
}

var describeListStoragesByPickUpPoint = testcommon.MethodDescriptor(
	"ListStoragesByPickUpPoint",
	"Get list of storages",
)

var describeGetStorageByInstance = testcommon.MethodDescriptor(
	"GetStorageByInstance",
	"Get storage by identifier",
)

func (self *StorageServiceIntegrationTestSuite) TestListStoragesByPickUpPointPositive(t provider.T) {
	var (
		skUser    models.User
		pup       models.PickUpPoint
		reference []storage.Storage
	)

	describeListStoragesByPickUpPoint(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepValue(t, "pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				v := models_om.PickUpPointExample("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		skugen := psql.GeneratorStepValue(t, "user", &skUser,
			func() (models.User, uuid.UUID) {
				v := models_om.UserDefault(nullable.None[string]()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				v := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					pupgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		cgen := psql.GeneratorStepNewList(t, "categories",
			func(i uint) (models.Category, uuid.UUID) {
				v := models_om.CategoryRandomId().WithName(fmt.Sprint(i)).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewList(t, "products",
			func(i uint) (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		var instances []models.Instance
		igen := psql.GeneratorStepList(t, "instances", &instances,
			func(i uint) (models.Instance, uuid.UUID) {
				v := models_om.InstanceExample(
					fmt.Sprint(i),
					pgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		sgen := psql.GeneratorStepNewList(t, "storages",
			func(i uint) (records.Storage, uuid.UUID) {
				v := records_om.StorageActive(
					pupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				reference = append(reference, MapStorage(&v))
				return v, v.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 5)

		rugen := psql.GeneratorStepNewList(t, "renter users",
			func(i uint) (models.User, uuid.UUID) {
				v := models_om.UserDefault(nullable.Some(fmt.Sprint(i))).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		rgen := psql.GeneratorStepNewList(t, "renters",
			func(i uint) (models.Renter, uuid.UUID) {
				v := models_om.RenterWithUserId(rugen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		progen := psql.GeneratorStepNewList(t, "provision",
			func(i uint) (records.Provision, uuid.UUID) {
				v := records_om.ProvisionActive(
					rgen.Generate(),
					instances[i].Id,
					nullable.None[time.Time](),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 5)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[storage.Storage]
	var err error

	t.WithNewStep(
		"Get active storages for pick up point",
		func(sCtx provider.StepCtx) {
			result, err = self.service.ListStoragesByPickUpPoint(token.Token(skUser.Token), pup.Id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickUpPointId", pup.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[storage.Storage](t).ElementsMatchFunc(
		CompareStorage, reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *StorageServiceIntegrationTestSuite) TestListStoragesByPickUpPointsUnauthorizedUser(t provider.T) {
	var (
		user models.User
		pup  models.PickUpPoint
	)

	describeListStoragesByPickUpPoint(t,
		"Unauthorized access by plain user",
		"Check that Authorization:NoAccess error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepValue(t, "pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				v := models_om.PickUpPointExample("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		pugen := psql.GeneratorStepValue(t, "plain user", &user,
			func() (models.User, uuid.UUID) {
				v := models_om.UserDefault(nullable.None[string]()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.Add(pugen, 1)

		cgen := psql.GeneratorStepNewList(t, "categories",
			func(i uint) (models.Category, uuid.UUID) {
				v := models_om.CategoryRandomId().WithName(fmt.Sprint(i)).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewList(t, "products",
			func(i uint) (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		var instances []models.Instance
		igen := psql.GeneratorStepList(t, "instances", &instances,
			func(i uint) (models.Instance, uuid.UUID) {
				v := models_om.InstanceExample(
					fmt.Sprint(i),
					pgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		sgen := psql.GeneratorStepNewList(t, "storages",
			func(i uint) (records.Storage, uuid.UUID) {
				v := records_om.StorageActive(
					pupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 5)

		rugen := psql.GeneratorStepNewList(t, "renter users",
			func(i uint) (models.User, uuid.UUID) {
				v := models_om.UserDefault(nullable.Some(fmt.Sprint(i))).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		rgen := psql.GeneratorStepNewList(t, "renters",
			func(i uint) (models.Renter, uuid.UUID) {
				v := models_om.RenterWithUserId(rugen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		progen := psql.GeneratorStepNewList(t, "provision",
			func(i uint) (records.Provision, uuid.UUID) {
				v := records_om.ProvisionActive(
					rgen.Generate(),
					instances[i].Id,
					nullable.None[time.Time](),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 5)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var err error

	t.WithNewStep(
		"Get active storages by pick up point",
		func(sCtx provider.StepCtx) {
			_, err = self.service.ListStoragesByPickUpPoint(token.Token(user.Token), pup.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("pickUpPointId", pup.Id),
	)

	// Assert
	var aerr cmnerrors.ErrorAuthorization
	var naerr cmnerrors.ErrorNoAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &aerr, "Error is Authentication")
	t.Require().ErrorAs(aerr, &naerr, "Error is NoAccess")
}

func (self *StorageServiceIntegrationTestSuite) TestGetStorageByInstancePositive(t provider.T) {
	var (
		instance  models.Instance
		reference storage.Storage
	)

	describeGetStorageByInstance(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepNewValue(t, "pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				v := models_om.PickUpPointExample("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		cgen := psql.GeneratorStepNewValue(t, "category",
			func() (models.Category, uuid.UUID) {
				v := models_om.CategoryRandomId().WithName("test").Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				v := models_om.ProductExmaple("test", cgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepValue(t, "instances", &instance,
			func() (models.Instance, uuid.UUID) {
				v := models_om.InstanceExample(
					"test",
					pgen.Generate(),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		sgen := psql.GeneratorStepNewValue(t, "storage",
			func() (records.Storage, uuid.UUID) {
				v := records_om.StorageActive(
					pupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				reference = MapStorage(&v)
				return v, v.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 1)

		rugen := psql.GeneratorStepNewList(t, "renter users",
			func(i uint) (models.User, uuid.UUID) {
				v := models_om.UserDefault(nullable.Some(fmt.Sprint(i))).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		rgen := psql.GeneratorStepNewList(t, "renters",
			func(i uint) (models.Renter, uuid.UUID) {
				v := models_om.RenterWithUserId(rugen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		progen := psql.GeneratorStepNewList(t, "provision",
			func(i uint) (records.Provision, uuid.UUID) {
				v := records_om.ProvisionActive(
					rgen.Generate(),
					instance.Id,
					nullable.None[time.Time](),
				).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result storage.Storage
	var err error

	t.WithNewStep(
		"Get active storage for instance",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetStorageByInstance(instance.Id)
		},
		allure.NewParameter("instanceId", instance.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[storage.Storage](t).EqualFunc(
		CompareStorage, reference, result, "Returned value matches reference",
	)
}

func (self *StorageServiceIntegrationTestSuite) TestGetStorageByInstanceNotFound(t provider.T) {
	var id uuid.UUID

	describeGetStorageByInstance(t,
		"Storage not found",
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
		"Get active storage for instance",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetStorageByInstance(id)
		},
		allure.NewParameter("instanceId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestStorageServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(StorageServiceIntegrationTestSuite))
}

