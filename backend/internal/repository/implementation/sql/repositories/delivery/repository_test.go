package delivery_test

import (
	"math/rand/v2"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	requests_om "rent_service/builders/mothers/domain/requests"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/delivery"
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

func CompareDelivery(reference requests.Delivery, actual requests.Delivery) bool {
	return reference.Id == actual.Id &&
		reference.CompanyId == actual.CompanyId &&
		reference.InstanceId == actual.InstanceId &&
		reference.FromId == actual.FromId &&
		reference.ToId == actual.ToId &&
		reference.DeliveryId == actual.DeliveryId &&
		psqlcommon.CompareTimeMicro(reference.ScheduledBeginDate, actual.ScheduledBeginDate) &&
		psqlcommon.CompareTimePtrMicro(reference.ActualBeginDate, actual.ActualBeginDate) &&
		psqlcommon.CompareTimeMicro(reference.ScheduledEndDate, actual.ScheduledEndDate) &&
		psqlcommon.CompareTimePtrMicro(reference.ActualEndDate, actual.ActualEndDate) &&
		reference.VerificationCode == actual.VerificationCode &&
		psqlcommon.CompareTimeMicro(reference.CreateDate, actual.CreateDate)
}

type DeliveryRepositoryTestSuite struct {
	suite.Suite
	repo delivery.IRepository
	psqlcommon.Context
}

func (self *DeliveryRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *DeliveryRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *DeliveryRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Delivery repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateDeliveryRepository()
	})
}

var describeCreate = testcommon.MethodDescriptor(
	"Create",
	"Create delivery",
)

var describeUpdate = testcommon.MethodDescriptor(
	"Update",
	"Update delivery",
)

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get delivery by id",
)

var describeGetActiveByPickUpPointId = testcommon.MethodDescriptor(
	"GetActiveByPickUpPointId",
	"Get active deliveries for pick up point",
)

var describeGetActiveByInstanceId = testcommon.MethodDescriptor(
	"GetActiveByInstanceId",
	"Get active deliveries for instance",
)

func CheckCreated(reference requests.Delivery, actual requests.Delivery) bool {
	return uuid.UUID{} != actual.Id &&
		reference.InstanceId == actual.InstanceId &&
		reference.FromId == actual.FromId &&
		reference.ToId == actual.ToId &&
		reference.DeliveryId == actual.DeliveryId &&
		psqlcommon.CompareTimeMicro(reference.ScheduledBeginDate, actual.ScheduledBeginDate) &&
		psqlcommon.CompareTimePtrMicro(reference.ActualBeginDate, actual.ActualBeginDate) &&
		psqlcommon.CompareTimeMicro(reference.ScheduledEndDate, actual.ScheduledEndDate) &&
		psqlcommon.CompareTimePtrMicro(reference.ActualEndDate, actual.ActualEndDate) &&
		reference.VerificationCode == actual.VerificationCode &&
		psqlcommon.CompareTimeMicro(reference.CreateDate, actual.CreateDate)
}

func (self *DeliveryRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeCreate(t,
		"Simple create test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&company)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomCreated(company.Id, instance.Id,
					fromPuP.Id, toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result requests.Delivery
	var err error

	t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(reference)
	}, allure.NewParameter("delivery", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[requests.Delivery](t).EqualFunc(
		CheckCreated, reference, result, "Same delivery with non null uuid",
	)
}

func (self *DeliveryRepositoryTestSuite) TestCreateNotFound(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeCreate(t,
		"Delivery company not exists",
		"Checks that error returned and mapped to NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and DON'T insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomCreated(company.Id, instance.Id,
					fromPuP.Id, toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(reference)
	}, allure.NewParameter("delivery", reference))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *DeliveryRepositoryTestSuite) TestUpdatePositive(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeUpdate(t,
		"Simple update test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&company)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create and insert new delivery", func(sCtx provider.StepCtx) {
			builder := requests_om.DeliveryRandomSent(company.Id, instance.Id,
				fromPuP.Id, toPuP.Id,
				nullable.None[string](),
				nullable.None[time.Time](),
				nullable.None[time.Time](),
				nullable.None[time.Time](),
				nullable.None[string](),
				nullable.None[time.Time](),
			)
			reference = testcommon.AssignParameter(sCtx, "delivery",
				builder.Build())
			created := builder.
				WithActualBeginDate(nullable.None[time.Time]()).
				Build()
			self.Inserter.InsertDelivery(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update delivery (send)", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("update", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *DeliveryRepositoryTestSuite) TestUpdateNotFound(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeUpdate(t,
		"Delivery update to conflicting state",
		"Checks that error returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&company)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create and insert delivery", func(sCtx provider.StepCtx) {
			builder := requests_om.DeliveryRandomCreated(
				company.Id, instance.Id,
				fromPuP.Id, toPuP.Id,
				nullable.None[string](),
				nullable.None[time.Time](),
				nullable.None[time.Time](),
				nullable.None[string](),
				nullable.None[time.Time](),
			)
			created := builder.Build()
			reference = testcommon.AssignParameter(sCtx, "delivery",
				builder.WithActualEndDate(nullable.Some(time.Now())).Build())
			self.Inserter.InsertDelivery(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update delivery", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("delivery", reference))

	// Assert
	t.Require().NotNil(err, "Error must be returned")
}

func (self *DeliveryRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeGetById(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&company)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create and insert delivery", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomAccepted(company.Id, instance.Id,
					fromPuP.Id, toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
			self.Inserter.InsertDelivery(&reference)
		})
	})

	// Act
	var result requests.Delivery
	var err error

	t.WithNewStep("Get delivery by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("deliveryId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[requests.Delivery](t).EqualFunc(
		CompareDelivery, reference, result, "Same delivery value",
	)
}

func (self *DeliveryRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Delivery not found",
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

	t.WithNewStep("Get delivery by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("deliveryId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func GenerateRandomActive(
	company models.DeliveryCompany,
	instance models.Instance,
	from models.PickUpPoint,
	to models.PickUpPoint,
) requests.Delivery {
	v := rand.Float64()

	if v < 0.5 {
		return requests_om.DeliveryRandomSent(
			company.Id, instance.Id,
			from.Id, to.Id,
			nullable.None[string](),
			nullable.None[time.Time](),
			nullable.None[time.Time](),
			nullable.None[time.Time](),
			nullable.None[string](),
			nullable.None[time.Time](),
		).Build()
	} else {
		return requests_om.DeliveryRandomCreated(
			company.Id, instance.Id,
			from.Id, to.Id,
			nullable.None[string](),
			nullable.None[time.Time](),
			nullable.None[time.Time](),
			nullable.None[string](),
			nullable.None[time.Time](),
		).Build()
	}
}

func WrapGenerateRandomActive(
	companies []models.DeliveryCompany,
	instances []models.Instance,
	target *models.PickUpPoint,
	others []models.PickUpPoint,
) func(i *int) requests.Delivery {
	get := func(i int) (models.PickUpPoint, models.PickUpPoint) {
		v := rand.Float64()

		if v < 0.5 {
			return *target, others[i]
		} else {
			return others[i], *target
		}
	}

	return func(i *int) requests.Delivery {
		from, to := get(*i)
		return GenerateRandomActive(
			companies[*i], instances[*i], from, to,
		)
	}
}

func (self *DeliveryRepositoryTestSuite) TestGetActiveByPickUpPointIdPositive(t provider.T) {
	var (
		companies  []models.DeliveryCompany
		categories []models.Category
		products   []models.Product
		instances  []models.Instance
		target     models.PickUpPoint
		others     []models.PickUpPoint
		references []requests.Delivery
	)

	describeGetActiveByPickUpPointId(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			companies = testcommon.AssignParameter(sCtx, "companies",
				collect.Do(
					models_om.DeliveryCompanyExample("1"),
					models_om.DeliveryCompanyExample("2"),
					models_om.DeliveryCompanyExample("3"),
					models_om.DeliveryCompanyExample("4"),
					models_om.DeliveryCompanyExample("5"),
				),
			)
			psql.BulkInsert(self.Inserter.InsertDeliveryCompany, companies...)
		})

		t.WithNewStep("Create and insert categories", func(sCtx provider.StepCtx) {
			categories = testcommon.AssignParameter(sCtx, "categories",
				collect.Do(
					models_om.CategoryRandomId().WithName("1"),
					models_om.CategoryRandomId().WithName("2"),
					models_om.CategoryRandomId().WithName("3"),
					models_om.CategoryRandomId().WithName("4"),
					models_om.CategoryRandomId().WithName("5"),
				),
			)
			psql.BulkInsert(self.Inserter.InsertCategory, categories...)
		})

		t.WithNewStep("Create and insert products", func(sCtx provider.StepCtx) {
			products = testcommon.AssignParameter(sCtx, "products",
				collect.Do(
					models_om.ProductExmaple("1", categories[0].Id),
					models_om.ProductExmaple("2", categories[1].Id),
					models_om.ProductExmaple("3", categories[2].Id),
					models_om.ProductExmaple("4", categories[3].Id),
					models_om.ProductExmaple("5", categories[4].Id),
				),
			)
			psql.BulkInsert(self.Inserter.InsertProduct, products...)
		})

		t.WithNewStep("Create and insert instances", func(sCtx provider.StepCtx) {
			instances = testcommon.AssignParameter(sCtx, "instances",
				collect.Do(
					models_om.InstanceExample("1", products[0].Id),
					models_om.InstanceExample("2", products[1].Id),
					models_om.InstanceExample("3", products[2].Id),
					models_om.InstanceExample("4", products[3].Id),
					models_om.InstanceExample("5", products[4].Id),
				),
			)
			psql.BulkInsert(self.Inserter.InsertInstance, instances...)
		})

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			target = testcommon.AssignParameter(sCtx, "target",
				models_om.PickUpPointExample("Target").Build(),
			)
			others = testcommon.AssignParameter(sCtx, "others",
				collect.Do(
					models_om.PickUpPointExample("1"),
					models_om.PickUpPointExample("2"),
					models_om.PickUpPointExample("3"),
					models_om.PickUpPointExample("4"),
					models_om.PickUpPointExample("5"),
				),
			)
			self.Inserter.InsertPickUpPoint(&target)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, others...)
		})

		t.WithNewStep("Create and insert deliveries", func(sCtx provider.StepCtx) {
			references = testcommon.AssignParameter(sCtx, "deliveries",
				collection.Collect(
					collection.MapIterator(
						WrapGenerateRandomActive(companies, instances, &target, others),
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
			psql.BulkInsert(self.Inserter.InsertDelivery, references...)
		})
	})

	// Act
	var result collection.Collection[requests.Delivery]
	var err error

	t.WithNewStep("Get active delivery by pick up point id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetActiveByPickUpPointId(target.Id)
	}, allure.NewParameter("pickUpPointId", target.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[requests.Delivery](t).ElementsMatchFunc(
		CompareDelivery, references, collection.Collect(result.Iter()), "Same collection",
	)
}

func (self *DeliveryRepositoryTestSuite) TestGetActiveByPickUpPointIdNotFound(t provider.T) {
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

	t.WithNewStep("Get active delivery by pick up point id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetActiveByPickUpPointId(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *DeliveryRepositoryTestSuite) TestGetActiveByInstanceIdPositive(t provider.T) {
	var (
		company   models.DeliveryCompany
		category  models.Category
		product   models.Product
		instance  models.Instance
		fromPuP   models.PickUpPoint
		toPuP     models.PickUpPoint
		reference requests.Delivery
	)

	describeGetActiveByInstanceId(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			company = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&company)
		})

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

		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "from",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "to",
				models_om.PickUpPointExample("To").Build(),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, fromPuP, toPuP)
		})

		t.WithNewStep("Create and insert delivery", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "delivery",
				GenerateRandomActive(company, instance, fromPuP, toPuP),
			)
			self.Inserter.InsertDelivery(&reference)
		})
	})

	// Act
	var result requests.Delivery
	var err error

	t.WithNewStep("Get active delivery by instance id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetActiveByInstanceId(instance.Id)
	}, allure.NewParameter("instanceId", instance.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[requests.Delivery](t).EqualFunc(
		CompareDelivery, reference, result, "Same delivery value",
	)
}

func (self *DeliveryRepositoryTestSuite) TestGetActiveByInstanceIdNotFound(t provider.T) {
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

	t.WithNewStep("Get active delivery by instance id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetActiveByInstanceId(id)
	}, allure.NewParameter("instanceId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type DeliveryCompanyRepositoryTestSuite struct {
	suite.Suite
	repo delivery.ICompanyRepository
	psqlcommon.Context
}

func (self *DeliveryCompanyRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *DeliveryCompanyRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *DeliveryCompanyRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Delivery Company repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreateDeliveryCompanyRepository()
	})
}

var describeCompanyGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get delivery company by id",
)

var describeCompanyGetAll = testcommon.MethodDescriptor(
	"GetAll",
	"Get all delivery comapnies",
)

func (self *DeliveryCompanyRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var reference models.DeliveryCompany

	describeCompanyGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert company", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "company",
				models_om.DeliveryCompanyExample("1").Build(),
			)
			self.Inserter.InsertDeliveryCompany(&reference)
		})
	})

	// Act
	var result models.DeliveryCompany
	var err error

	t.WithNewStep("Get delivery company by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("companyId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same delivery value")
}

func (self *DeliveryCompanyRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Delivery company not found",
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

	t.WithNewStep("Get delivery company by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("companyId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *DeliveryCompanyRepositoryTestSuite) TestGetAllPositive(t provider.T) {
	var reference []models.DeliveryCompany

	describeCompanyGetAll(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert companies", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "company",
				collect.Do(
					models_om.DeliveryCompanyExample("1"),
					models_om.DeliveryCompanyExample("2"),
					models_om.DeliveryCompanyExample("3"),
					models_om.DeliveryCompanyExample("4"),
					models_om.DeliveryCompanyExample("5"),
				),
			)
			psql.BulkInsert(self.Inserter.InsertDeliveryCompany, reference...)
		})
	})

	// Act
	var result collection.Collection[models.DeliveryCompany]
	var err error

	t.WithNewStep("Get all delivery company", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(reference, collection.Collect(result.Iter()),
		"Same company values")
}

func (self *DeliveryCompanyRepositoryTestSuite) TestGetAllEmpty(t provider.T) {
	describeCompanyGetAll(t,
		"No delivery companies",
		"Checks that method return empty collection withour error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[models.DeliveryCompany]
	var err error

	t.WithNewStep("Get all delivery companies", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()),
		"Collection is empty")
}

func TestDeliveryRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryRepositoryTestSuite))
}

func TestDeliveryCompanyRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryCompanyRepositoryTestSuite))
}

