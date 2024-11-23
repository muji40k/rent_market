package delivery_test

import (
	"fmt"
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	records_om "rent_service/builders/mothers/domain/records"
	requests_om "rent_service/builders/mothers/domain/requests"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/delivery/implementations/dummy"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/delivery"
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

func MapDelivery(value *requests.Delivery) delivery.Delivery {
	return delivery.Delivery{
		Id:         value.Id,
		CompanyId:  value.CompanyId,
		InstanceId: value.InstanceId,
		FromId:     value.FromId,
		ToId:       value.ToId,
		BeginDate: delivery.Dates{
			Scheduled: date.New(value.ScheduledBeginDate),
			Actual: nullcommon.CopyPtrIfSome(
				nullable.Map(
					nullable.FromPtr(value.ActualBeginDate),
					func(value *time.Time) date.Date {
						return date.New(*value)
					},
				),
			),
		},
		EndDate: delivery.Dates{
			Scheduled: date.New(value.ScheduledEndDate),
			Actual: nullcommon.CopyPtrIfSome(
				nullable.Map(
					nullable.FromPtr(value.ActualEndDate),
					func(value *time.Time) date.Date {
						return date.New(*value)
					},
				),
			),
		},
		VerificationCode: value.VerificationCode,
		CreateDate:       date.New(value.CreateDate),
	}
}

func CompareDelivery(reference delivery.Delivery, actual delivery.Delivery) bool {
	return reference.Id == actual.Id &&
		reference.CompanyId == actual.CompanyId &&
		reference.InstanceId == actual.InstanceId &&
		reference.FromId == actual.FromId &&
		reference.ToId == actual.ToId &&
		psqlcommon.CompareTimePtrMicro(
			psqlcommon.UnwrapDate(reference.BeginDate.Actual),
			psqlcommon.UnwrapDate(actual.BeginDate.Actual),
		) &&
		psqlcommon.CompareTimeMicro(reference.BeginDate.Scheduled.Time, actual.BeginDate.Scheduled.Time) &&
		psqlcommon.CompareTimePtrMicro(
			psqlcommon.UnwrapDate(reference.EndDate.Actual),
			psqlcommon.UnwrapDate(actual.EndDate.Actual),
		) &&
		psqlcommon.CompareTimeMicro(reference.EndDate.Scheduled.Time, actual.EndDate.Scheduled.Time) &&
		reference.VerificationCode == actual.VerificationCode &&
		psqlcommon.CompareTimeMicro(reference.CreateDate.Time, actual.CreateDate.Time)
}

type DeliveryServiceIntegrationTestSuite struct {
	suite.Suite
	service  delivery.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *DeliveryServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *DeliveryServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *DeliveryServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Delivery service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateDeliveryService()
	})
}

var describeAcceptDelivery = testcommon.MethodDescriptor(
	"AcceptDelivery",
	"Accept delivery",
)

var describeCreateDelivery = testcommon.MethodDescriptor(
	"CreateDelivery",
	"Create delivery",
)

var describeSendDelivery = testcommon.MethodDescriptor(
	"SendDelivery",
	"Send delivery",
)

var describeListDeliveriesByPickUpPoint = testcommon.MethodDescriptor(
	"ListDeliveriesByPickUpPoint",
	"Get list of deliveries by pick up point id",
)

var describeGetDeliveryByInstance = testcommon.MethodDescriptor(
	"GetDeliveryByInstance",
	"Get active delivery by instance id",
)

func (self *DeliveryServiceIntegrationTestSuite) TestAcceptDeliveryPositive(t provider.T) {
	var COMMENT string = "Update"

	var (
		skUser       models.User
		refdelivery  requests.Delivery
		tempPhotoIds []uuid.UUID
		form         delivery.AcceptForm
	)

	describeAcceptDelivery(t,
		"Accept existed delivery",
		"Accept existed delivery by storekeeper",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					tpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("test").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepValue(t, "delivery", &refdelivery,
			func() (requests.Delivery, uuid.UUID) {
				d := requests_om.DeliveryRandomSent(
					dcgen.Generate(),
					igen.Generate(),
					fpupgen.Generate(),
					tpupgen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 1)

		photogen := psql.GeneratorStepNewList(t, "temp photos",
			func(i uint) (models.TempPhoto, uuid.UUID) {
				path := self.sContext.PhotoRegistry.SaveTempPhoto(
					models_om.ImagePNGContent(nullable.None[int]()),
				)
				photo := models_om.TempPhotoExampleUploaded(
					fmt.Sprint(i),
					nullable.None[time.Time](),
				).WithPath(
					nullable.Some(path),
				).Build()
				tempPhotoIds = append(tempPhotoIds, photo.Id)
				return photo, photo.Id
			},
			self.rContext.Inserter.InsertTempPhoto,
		)
		gg.Add(photogen, 5)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.AcceptForm{
				DeliveryId:       refdelivery.Id,
				Comment:          &COMMENT,
				VerificationCode: refdelivery.VerificationCode,
				TempPhotos:       tempPhotoIds,
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			err = self.service.AcceptDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *DeliveryServiceIntegrationTestSuite) TestAcceptDeliveryConflict(t provider.T) {
	var COMMENT string = "Update"

	var (
		skUser       models.User
		refdelivery  requests.Delivery
		tempPhotoIds []uuid.UUID
		form         delivery.AcceptForm
	)

	describeAcceptDelivery(t,
		"Attemp to accept delivery in wrong state",
		"Check that Conflict error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					tpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		sgen := psql.GeneratorStepNewValue(t, "storage",
			func() (records.Storage, uuid.UUID) {
				s := records_om.StorageActive(
					tpupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return s, s.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("test").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepValue(t, "delivery", &refdelivery,
			func() (requests.Delivery, uuid.UUID) {
				d := requests_om.DeliveryRandomAccepted(
					dcgen.Generate(),
					igen.Generate(),
					fpupgen.Generate(),
					tpupgen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 1)

		photogen := psql.GeneratorStepNewList(t, "temp photos",
			func(i uint) (models.TempPhoto, uuid.UUID) {
				path := self.sContext.PhotoRegistry.SaveTempPhoto(
					models_om.ImagePNGContent(nullable.None[int]()),
				)
				photo := models_om.TempPhotoExampleUploaded(
					fmt.Sprint(i),
					nullable.None[time.Time](),
				).WithPath(
					nullable.Some(path),
				).Build()
				tempPhotoIds = append(tempPhotoIds, photo.Id)
				return photo, photo.Id
			},
			self.rContext.Inserter.InsertTempPhoto,
		)
		gg.Add(photogen, 5)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.AcceptForm{
				DeliveryId:       refdelivery.Id,
				Comment:          &COMMENT,
				VerificationCode: refdelivery.VerificationCode,
				TempPhotos:       tempPhotoIds,
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Try to accept already accepted delivery",
		func(sCtx provider.StepCtx) {
			err = self.service.AcceptDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var cerr cmnerrors.ErrorConflict

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &cerr, "Error Conflict must be returned")
}

func (self *DeliveryServiceIntegrationTestSuite) TestCreateDeliveryPositive(t provider.T) {
	var (
		skUser    models.User
		reference delivery.Delivery
		form      delivery.CreateForm
	)

	describeCreateDelivery(t,
		"Create delivery",
		"Must return created delivery without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					fpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		sgen := psql.GeneratorStepNewValue(t, "storage",
			func() (records.Storage, uuid.UUID) {
				s := records_om.StorageActive(
					fpupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return s, s.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("dummy").
					WithId(dummy.Id).
					Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			mdelivery := testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomId().
					WithId(uuid.UUID{}).
					WithCompanyId(dcgen.Generate()).
					WithInstanceId(igen.Generate()).
					WithFromId(fpupgen.Generate()).
					WithToId(tpupgen.Generate()).
					Build(),
			)
			reference = MapDelivery(&mdelivery)
		})

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.CreateForm{
				InstanceId: igen.Generate(),
				From:       fpupgen.Generate(),
				To:         tpupgen.Generate(),
			}
		})

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result delivery.Delivery
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = self.service.CreateDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[delivery.Delivery](t).EqualFunc(
		func(r delivery.Delivery, a delivery.Delivery) bool {
			return uuid.UUID{} != a.Id &&
				r.CompanyId == a.CompanyId &&
				r.InstanceId == a.InstanceId &&
				r.FromId == a.FromId &&
				r.ToId == a.ToId
		},
		reference, result, "Result matches expected reference",
	)
}

func (self *DeliveryServiceIntegrationTestSuite) TestCreateDeliveryConflict(t provider.T) {
	var (
		skUser models.User
		form   delivery.CreateForm
	)

	describeCreateDelivery(t,
		"Instance is not available for delivery (stored in other pup point)",
		"Must return ErrorConflic if instance isn't stored in pick up point",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		ugen := psql.GeneratorStepNewValue(t, "user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserDefault(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(ugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					fpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		periodGen := psql.GeneratorStepNewValue(t, "period",
			func() (models.Period, uuid.UUID) {
				s := models_om.PeriodDay().Build()
				return s, s.Id
			},
			self.rContext.Inserter.InsertPeriod,
		)
		gg.AddFinish(periodGen)

		rentGen := psql.GeneratorStepNewValue(t, "rent",
			func() (records.Rent, uuid.UUID) {
				s := records_om.RentActive(
					ugen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
					periodGen.Generate(),
				).Build()
				return s, s.Id
			},
			self.rContext.Inserter.InsertRent,
		)
		gg.Add(rentGen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("dummy").
					WithId(dummy.Id).
					Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.CreateForm{
				InstanceId: igen.Generate(),
				From:       fpupgen.Generate(),
				To:         tpupgen.Generate(),
			}
		})

		gg.Generate()
		gg.Finish()
	})

	// Act
	var err error

	t.WithNewStep("Create delivery",
		func(sCtx provider.StepCtx) {
			_, err = self.service.CreateDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var aerr cmnerrors.ErrorAuthorization
	var naerr cmnerrors.ErrorNoAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &aerr, "Error is Authorization")
	t.Require().ErrorAs(err, &naerr, "Error is NoAccess")
}

func (self *DeliveryServiceIntegrationTestSuite) TestSendDeliveryPositive(t provider.T) {
	var (
		skUser       models.User
		refdelivery  requests.Delivery
		tempPhotoIds []uuid.UUID
		form         delivery.SendForm
	)

	describeSendDelivery(t,
		"Send created delivery",
		"Send created delivery by storekeeper",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					fpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		sgen := psql.GeneratorStepNewValue(t, "storage",
			func() (records.Storage, uuid.UUID) {
				p := records_om.StorageActive(
					fpupgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertStorage,
		)
		gg.Add(sgen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("test").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepValue(t, "delivery", &refdelivery,
			func() (requests.Delivery, uuid.UUID) {
				d := requests_om.DeliveryRandomCreated(
					dcgen.Generate(),
					igen.Generate(),
					fpupgen.Generate(),
					tpupgen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 1)

		photogen := psql.GeneratorStepNewList(t, "temp photos",
			func(i uint) (models.TempPhoto, uuid.UUID) {
				path := self.sContext.PhotoRegistry.SaveTempPhoto(
					models_om.ImagePNGContent(nullable.None[int]()),
				)
				photo := models_om.TempPhotoExampleUploaded(
					fmt.Sprint(i),
					nullable.None[time.Time](),
				).WithPath(
					nullable.Some(path),
				).Build()
				tempPhotoIds = append(tempPhotoIds, photo.Id)
				return photo, photo.Id
			},
			self.rContext.Inserter.InsertTempPhoto,
		)
		gg.Add(photogen, 5)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.SendForm{
				DeliveryId:       refdelivery.Id,
				VerificationCode: refdelivery.VerificationCode,
				TempPhotos:       tempPhotoIds,
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			err = self.service.SendDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *DeliveryServiceIntegrationTestSuite) TestSednDeliveryConflict(t provider.T) {
	var (
		skUser       models.User
		refdelivery  requests.Delivery
		tempPhotoIds []uuid.UUID
		form         delivery.SendForm
	)

	describeSendDelivery(t,
		"Attemp to send delivery that already sent",
		"Check that Conflict error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					fpupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepNewValue(t, "instance",
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("test").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepValue(t, "delivery", &refdelivery,
			func() (requests.Delivery, uuid.UUID) {
				d := requests_om.DeliveryRandomSent(
					dcgen.Generate(),
					igen.Generate(),
					fpupgen.Generate(),
					tpupgen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 1)

		photogen := psql.GeneratorStepNewList(t, "temp photos",
			func(i uint) (models.TempPhoto, uuid.UUID) {
				path := self.sContext.PhotoRegistry.SaveTempPhoto(
					models_om.ImagePNGContent(nullable.None[int]()),
				)
				photo := models_om.TempPhotoExampleUploaded(
					fmt.Sprint(i),
					nullable.None[time.Time](),
				).WithPath(
					nullable.Some(path),
				).Build()
				tempPhotoIds = append(tempPhotoIds, photo.Id)
				return photo, photo.Id
			},
			self.rContext.Inserter.InsertTempPhoto,
		)
		gg.Add(photogen, 5)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = delivery.SendForm{
				DeliveryId:       refdelivery.Id,
				VerificationCode: refdelivery.VerificationCode,
				TempPhotos:       tempPhotoIds,
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Try to accept already sent delivery",
		func(sCtx provider.StepCtx) {
			err = self.service.SendDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var cerr cmnerrors.ErrorConflict

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &cerr, "Error Conflict must be returned")
}

func (self *DeliveryServiceIntegrationTestSuite) TestListDeliveriesByPickUpPointPositive(t provider.T) {
	var (
		pup       models.PickUpPoint
		skUser    models.User
		reference []delivery.Delivery
	)

	describeListDeliveriesByPickUpPoint(t,
		"List deliveries by pick up point (storekeeper)",
		"All values must be returned wihtout error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepValue(t, "target pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("Target").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		opupgen := psql.GeneratorStepNewList(t, "other pick up point",
			func(i uint) (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample(fmt.Sprint(i)).Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(opupgen)

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		rugen := psql.GeneratorStepNewList(t, "renter users",
			func(i uint) (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.Some(fmt.Sprint(i))).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					pupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		rgen := psql.GeneratorStepNewList(t, "renters",
			func(i uint) (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewList(t, "categroies",
			func(i uint) (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName(fmt.Sprint(i)).Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewList(t, "products",
			func(i uint) (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple(fmt.Sprint(i), cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		var instances []models.Instance
		igen := psql.GeneratorStepList(t, "instances", &instances,
			func(i uint) (models.Instance, uuid.UUID) {
				v := models_om.InstanceExample(fmt.Sprint(i), pgen.Generate()).Build()
				return v, v.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewList(t, "provisions",
			func(_ uint) (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 10)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("dummy").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepNewList(t, "deliveries",
			func(i uint) (requests.Delivery, uuid.UUID) {
				var d requests.Delivery

				if 0 == i%2 {
					d = requests_om.DeliveryRandomSent(
						dcgen.Generate(),
						instances[i].Id,
						pupgen.Generate(),
						opupgen.Generate(),
						nullable.None[string](),
						nullable.None[time.Time](),
						nullable.None[time.Time](),
						nullable.None[time.Time](),
						nullable.None[string](),
						nullable.None[time.Time](),
					).Build()
				} else {
					d = requests_om.DeliveryRandomCreated(
						dcgen.Generate(),
						instances[i].Id,
						opupgen.Generate(),
						pupgen.Generate(),
						nullable.None[string](),
						nullable.None[time.Time](),
						nullable.None[time.Time](),
						nullable.None[string](),
						nullable.None[time.Time](),
					).Build()
				}

				reference = append(reference, MapDelivery(&d))

				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 10)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var err error
	var result collection.Collection[delivery.Delivery]

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = self.service.ListDeliveriesByPickUpPoint(token.Token(skUser.Token), pup.Id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickIpPointId", pup.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[delivery.Delivery](t).ElementsMatchFunc(
		CompareDelivery, reference, collection.Collect(result.Iter()),
		"Elements must match",
	)
}

func (self *DeliveryServiceIntegrationTestSuite) TestListDeliveriesByPickUpPointNotFound(t provider.T) {
	var (
		skUser models.User
		id     uuid.UUID
	)

	describeListDeliveriesByPickUpPoint(t,
		"Pick up point not found",
		"Error NotFound must be returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		pupgen := psql.GeneratorStepNewValue(t, "pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("Target").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					pupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create id to fetch", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = self.service.ListDeliveriesByPickUpPoint(token.Token(skUser.Token), id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickIpPointId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *DeliveryServiceIntegrationTestSuite) TestListDeliveriesByInstancePositive(t provider.T) {
	var (
		user      models.User
		instance  models.Instance
		reference delivery.Delivery
	)

	describeGetDeliveryByInstance(t,
		"List deliveries by instance (user)",
		"All values must be returned wihtout error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		fpupgen := psql.GeneratorStepNewValue(t, "from pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("From").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(fpupgen)

		tpupgen := psql.GeneratorStepNewValue(t, "to pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("To").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(tpupgen)

		ugen := psql.GeneratorStepValue(t, "user", &user,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(ugen)

		rugen := psql.GeneratorStepNewValue(t, "renter user",
			func() (models.User, uuid.UUID) {
				user := models_om.UserRenter(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(rugen)

		rgen := psql.GeneratorStepNewValue(t, "renter",
			func() (models.Renter, uuid.UUID) {
				r := models_om.RenterWithUserId(rugen.Generate()).Build()
				return r, r.Id
			},
			self.rContext.Inserter.InsertRenter,
		)
		gg.AddFinish(rgen)

		cgen := psql.GeneratorStepNewValue(t, "categroy",
			func() (models.Category, uuid.UUID) {
				c := models_om.CategoryRandomId().WithName("Test").Build()
				return c, c.Id
			},
			self.rContext.Inserter.InsertCategory,
		)
		gg.AddFinish(cgen)

		pgen := psql.GeneratorStepNewValue(t, "product",
			func() (models.Product, uuid.UUID) {
				p := models_om.ProductExmaple("1", cgen.Generate()).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProduct,
		)
		gg.AddFinish(pgen)

		igen := psql.GeneratorStepValue(t, "instance", &instance,
			func() (models.Instance, uuid.UUID) {
				i := models_om.InstanceExample("1", pgen.Generate()).Build()
				return i, i.Id
			},
			self.rContext.Inserter.InsertInstance,
		)
		gg.AddFinish(igen)

		progen := psql.GeneratorStepNewValue(t, "provision",
			func() (records.Provision, uuid.UUID) {
				p := records_om.ProvisionActive(
					rgen.Generate(),
					igen.Generate(),
					nullable.None[time.Time](),
				).Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertProvision,
		)
		gg.Add(progen, 1)

		pergen := psql.GeneratorStepNewValue(t, "period",
			func() (models.Period, uuid.UUID) {
				p := models_om.PeriodMonth().Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertPeriod,
		)
		gg.AddFinish(pergen)

		rrgen := psql.GeneratorStepNewValue(t, "rent request",
			func() (requests.Rent, uuid.UUID) {
				dc := requests_om.Rent(
					igen.Generate(),
					ugen.Generate(),
					tpupgen.Generate(),
					pergen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertRentRequest,
		)
		gg.Add(rrgen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery company",
			func() (models.DeliveryCompany, uuid.UUID) {
				dc := models_om.DeliveryCompanyExample("test").Build()
				return dc, dc.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.AddFinish(dcgen)

		dgen := psql.GeneratorStepNewValue(t, "delivery",
			func() (requests.Delivery, uuid.UUID) {
				d := requests_om.DeliveryRandomSent(
					dcgen.Generate(),
					igen.Generate(),
					fpupgen.Generate(),
					tpupgen.Generate(),
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build()
				reference = MapDelivery(&d)
				return d, d.Id
			},
			self.rContext.Inserter.InsertDelivery,
		)
		gg.Add(dgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var err error
	var result delivery.Delivery

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetDeliveryByInstance(token.Token(user.Token), instance.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("instanceId", instance.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[delivery.Delivery](t).EqualFunc(
		CompareDelivery, reference, result, "Same value returned",
	)
}

func (self *DeliveryServiceIntegrationTestSuite) TestListDeliveriesByInstanceNotFound(t provider.T) {
	var (
		skUser models.User
		id     uuid.UUID
	)

	describeGetDeliveryByInstance(t,
		"Instance not found",
		"Error NotFound must be returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		skugen := psql.GeneratorStepValue(t, "storekeeper user", &skUser,
			func() (models.User, uuid.UUID) {
				user := models_om.UserStorekeeper(nullable.None[string]()).Build()
				return user, user.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.AddFinish(skugen)

		pupgen := psql.GeneratorStepNewValue(t, "pick up point",
			func() (models.PickUpPoint, uuid.UUID) {
				pup := models_om.PickUpPointExample("Target").Build()
				return pup, pup.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		skgen := psql.GeneratorStepNewValue(t, "storekeeper",
			func() (models.Storekeeper, uuid.UUID) {
				sk := models_om.StorekeeperWithUserId(
					skugen.Generate(),
					pupgen.Generate(),
				).Build()
				return sk, sk.Id
			},
			self.rContext.Inserter.InsertStorekeeper,
		)
		gg.Add(skgen, 1)

		gg.Generate()
		gg.Finish()

		t.WithNewStep("Create id to fetch", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})
	})

	// Act
	var err error

	t.WithNewStep("Get delivery by instance",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetDeliveryByInstance(token.Token(skUser.Token), id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("instanceId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is NotFound")
}

func MapDeliveryCompany(value *models.DeliveryCompany) delivery.DeliveryCompany {
	return delivery.DeliveryCompany{
		Id:          value.Id,
		Name:        value.Name,
		Site:        value.Site,
		PhoneNumber: value.PhoneNumber,
		Description: value.Description,
	}
}

type DeliveryCompanyServiceIntegrationTestSuite struct {
	suite.Suite
	service  delivery.ICompanyService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Delivery company service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreateDeliveryCompanyService()
	})
}

var describeGetDeliveryCompanyById = testcommon.MethodDescriptor(
	"GetDeliveryCompanyById",
	"Get delivery company by id",
)

var describeListDeliveryCompanies = testcommon.MethodDescriptor(
	"ListDeliveryCompanies",
	"List all delivery companies",
)

func (self *DeliveryCompanyServiceIntegrationTestSuite) TestListDeliveryCompaniesPositive(t provider.T) {
	var (
		user      models.User
		reference []delivery.DeliveryCompany
	)

	describeListDeliveryCompanies(t,
		"List all delivery companies",
		"All compnanies must be returned withou error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		ugen := psql.GeneratorStepValue(t, "user", &user,
			func() (models.User, uuid.UUID) {
				u := models_om.UserDefault(nullable.None[string]()).Build()
				return u, u.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.Add(ugen, 1)

		dcgen := psql.GeneratorStepNewList(t, "delivery companies",
			func(i uint) (models.DeliveryCompany, uuid.UUID) {
				u := models_om.DeliveryCompanyExample(fmt.Sprint(i)).Build()
				reference = append(reference, MapDeliveryCompany(&u))
				return u, u.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.Add(dcgen, 10)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[delivery.DeliveryCompany]
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = self.service.ListDeliveryCompanies(token.Token(user.Token))
		},
		allure.NewParameter("token", user.Token),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All companies returned",
	)
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) TestListDeliveryCompaniesUnknownUser(t provider.T) {
	var user models.User

	describeListDeliveryCompanies(t,
		"Unauthorized user can't access delivery companies",
		"Error must be returned and mapped to Authentication",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepNewList(t, "delivery companies",
				func(i uint) (models.DeliveryCompany, uuid.UUID) {
					u := models_om.DeliveryCompanyExample(fmt.Sprint(i)).Build()
					return u, u.Id
				},
				self.rContext.Inserter.InsertDeliveryCompany,
			), 10).
			Generate().
			Finish()
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = self.service.ListDeliveryCompanies(token.Token(user.Token))
		},
		allure.NewParameter("token", user.Token),
	)

	// Assert
	var aerr cmnerrors.ErrorAuthentication

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &aerr, "Error is authentication")
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) TestGetDeliveryCompanyByIdPositive(t provider.T) {
	var (
		user      models.User
		reference delivery.DeliveryCompany
	)

	describeGetDeliveryCompanyById(t,
		"Get delivery company by id",
		"Value must be returned without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		ugen := psql.GeneratorStepValue(t, "user", &user,
			func() (models.User, uuid.UUID) {
				u := models_om.UserDefault(nullable.None[string]()).Build()
				return u, u.Id
			},
			self.rContext.Inserter.InsertUser,
		)
		gg.Add(ugen, 1)

		dcgen := psql.GeneratorStepNewValue(t, "delivery companies",
			func() (models.DeliveryCompany, uuid.UUID) {
				u := models_om.DeliveryCompanyExample("1").Build()
				reference = MapDeliveryCompany(&u)
				return u, u.Id
			},
			self.rContext.Inserter.InsertDeliveryCompany,
		)
		gg.Add(dcgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result delivery.DeliveryCompany
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetDeliveryCompanyById(token.Token(user.Token), reference.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("companyId", reference.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Company returned")
}

func (self *DeliveryCompanyServiceIntegrationTestSuite) TestGetDeliveryCompanyByIdNotFound(t provider.T) {
	var (
		user models.User
		id   uuid.UUID
	)

	describeGetDeliveryCompanyById(t,
		"Delivery company not found",
		"Error NotFound must be returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepValue(t, "user", &user,
				func() (models.User, uuid.UUID) {
					u := models_om.UserDefault(nullable.None[string]()).Build()
					return u, u.Id
				},
				self.rContext.Inserter.InsertUser,
			), 1).
			Generate().
			Finish()

		t.WithNewStep("Create random id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetDeliveryCompanyById(token.Token(user.Token), id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("companyId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestDeliveryServiceIntegrationTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryServiceIntegrationTestSuite))
}

func TestDeliveryCompanyServiceIntegrationTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryCompanyServiceIntegrationTestSuite))
}

