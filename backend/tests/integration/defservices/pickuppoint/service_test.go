package pickuppoint_test

import (
	"fmt"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/pickuppoint"
	"rent_service/internal/logic/services/types/day"
	"rent_service/internal/logic/services/types/daytime"
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

func MapAddress(value *models.Address) pickuppoint.Address {
	return pickuppoint.Address{
		Country: value.Country,
		City:    value.City,
		Street:  value.Street,
		House:   value.House,
		Flat:    nullcommon.CopyPtrIfSome(nullable.FromPtr(value.Flat)),
	}
}

func MapPickUpPoint(value *models.PickUpPoint) pickuppoint.PickUpPoint {
	return pickuppoint.PickUpPoint{
		Id:       value.Id,
		Address:  MapAddress(&value.Address),
		Capacity: value.Capacity,
	}
}

type PickUpPointServiceIntegrationTestSuite struct {
	suite.Suite
	service  pickuppoint.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *PickUpPointServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *PickUpPointServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *PickUpPointServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Pick Up Point service",
	)
	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreatePickUpPointService()
	})
}

var describeListPickUpPoints = testcommon.MethodDescriptor(
	"ListPickUpPoints",
	"Get list of all pick up points",
)

var describeGetPickUpPointById = testcommon.MethodDescriptor(
	"GetPickUpPointById",
	"Get pick up point by identifier",
)

func (self *PickUpPointServiceIntegrationTestSuite) TestListPicUpPointsPositive(t provider.T) {
	var reference []pickuppoint.PickUpPoint

	describeListPickUpPoints(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepNewList(t, "pick up point",
				func(i uint) (models.PickUpPoint, uuid.UUID) {
					p := models_om.PickUpPointExample(fmt.Sprint(i)).Build()
					reference = append(reference, MapPickUpPoint(&p))
					return p, p.Id
				},
				self.rContext.Inserter.InsertPickUpPoint,
			), 5).
			Generate().
			Finish()

	})

	// Act
	var result collection.Collection[pickuppoint.PickUpPoint]
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		result, err = self.service.ListPickUpPoints()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointServiceIntegrationTestSuite) TestListPickUpPointsEmpty(t provider.T) {
	describeListPickUpPoints(t,
		"Empty pick up point set",
		"Check that empty set can by returned without error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[pickuppoint.PickUpPoint]
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		result, err = self.service.ListPickUpPoints()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()), "Collection is empty")
}

func (self *PickUpPointServiceIntegrationTestSuite) TestGetPickUpPointByIdPositive(t provider.T) {
	var reference pickuppoint.PickUpPoint

	describeGetPickUpPointById(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepNewValue(t, "pick up point",
				func() (models.PickUpPoint, uuid.UUID) {
					p := models_om.PickUpPointExample("test").Build()
					reference = MapPickUpPoint(&p)
					return p, p.Id
				},
				self.rContext.Inserter.InsertPickUpPoint,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var result pickuppoint.PickUpPoint
	var err error

	t.WithNewStep(
		"Get pick up point",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetPickUpPointById(reference.Id)
		},
		allure.NewParameter("pickUpPointId", reference.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value matches reference")
}

func (self *PickUpPointServiceIntegrationTestSuite) TestGetPickUpPointByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetPickUpPointById(t,
		"Pick up point not found",
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
		"Get pick up point",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetPickUpPointById(id)
		},
		allure.NewParameter("pickUpPointId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type PickUpPointPhotoServiceIntegrationTestSuite struct {
	suite.Suite
	service  pickuppoint.IPhotoService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *PickUpPointPhotoServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *PickUpPointPhotoServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *PickUpPointPhotoServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Pick Up Point photo service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreatePickUpPointPhotoService()
	})
}

var describeListPickUpPointPhotos = testcommon.MethodDescriptor(
	"ListPickUpPointPhotos",
	"Get list of all pick up point photos",
)

func (self *PickUpPointPhotoServiceIntegrationTestSuite) TestListPickUpPointPhotosPositive(t provider.T) {
	var (
		pup       models.PickUpPoint
		reference []uuid.UUID
	)

	describeListPickUpPointPhotos(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepValue(t, "pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				p := models_om.PickUpPointExample("test").Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.Add(pupgen, 1)

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
				self.rContext.Inserter.InsertPickUpPointPhoto(psql.NewPhoto(
					nullable.None[uuid.UUID](),
					pup.Id,
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

	t.WithNewStep("Get all pick up point photos", func(sCtx provider.StepCtx) {
		result, err = self.service.ListPickUpPointPhotos(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointPhotoServiceIntegrationTestSuite) TestListPickUpPointPhotosNotFound(t provider.T) {
	var id uuid.UUID

	describeListPickUpPointPhotos(t,
		"Pick up point not found",
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

	t.WithNewStep("Get all pick up point photos", func(sCtx provider.StepCtx) {
		_, err = self.service.ListPickUpPointPhotos(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func MapWorkingHours(value *models.PickUpPointWorkingHours) []pickuppoint.WorkingHours {
	return collection.Collect(
		collection.MapIterator(
			func(wh *collection.KV[time.Weekday, models.WorkingHours]) pickuppoint.WorkingHours {
				return pickuppoint.WorkingHours{
					Id:        wh.Value.Id,
					Day:       day.New(wh.Value.Day),
					StartHour: daytime.NewDuration(wh.Value.Begin),
					EndHour:   daytime.NewDuration(wh.Value.End),
				}
			},
			collection.HashMapIterator(value.Map),
		),
	)
}

type PickUpPointWorkingHoursServiceIntegrationTestSuite struct {
	suite.Suite
	service  pickuppoint.IWorkingHoursService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *PickUpPointWorkingHoursServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *PickUpPointWorkingHoursServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *PickUpPointWorkingHoursServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Pick Up Point working hours service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreatePickUpPointWorkingHoursService()
	})
}

var describeListPickUpPointWorkingHours = testcommon.MethodDescriptor(
	"ListPickUpPointPhotos",
	"Get pick up point working hours",
)

func (self *PickUpPointWorkingHoursServiceIntegrationTestSuite) TestListPickUpPointWorkingHoursPositive(t provider.T) {
	var (
		pup       models.PickUpPoint
		reference []pickuppoint.WorkingHours
	)

	describeListPickUpPointWorkingHours(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		gg := generator.NewGeneratorGroup()

		pupgen := psql.GeneratorStepValue(t, "pick up point", &pup,
			func() (models.PickUpPoint, uuid.UUID) {
				p := models_om.PickUpPointExample("test").Build()
				return p, p.Id
			},
			self.rContext.Inserter.InsertPickUpPoint,
		)
		gg.AddFinish(pupgen)

		whgen := psql.GeneratorStepNewValue(t, "working hours",
			func() (models.PickUpPointWorkingHours, uuid.UUID) {
				v := models_om.PickUpPointWorkingHours(
					pupgen.Generate(), collect.Do(
						models_om.WorkingHoursWeek(8*time.Hour, 20*time.Hour)...,
					)...,
				).Build()
				reference = MapWorkingHours(&v)
				return v, v.PickUpPointId
			},
			self.rContext.Inserter.InsertPickUpPointWorkingHours,
		)
		gg.Add(whgen, 1)

		gg.Generate()
		gg.Finish()
	})

	// Act
	var result collection.Collection[pickuppoint.WorkingHours]
	var err error

	t.WithNewStep("Get all pick up point working hours", func(sCtx provider.StepCtx) {
		result, err = self.service.ListPickUpPointWorkingHours(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointWorkingHoursServiceIntegrationTestSuite) TestListPickUpPointWorkingHoursNotFound(t provider.T) {
	var id uuid.UUID

	describeListPickUpPointWorkingHours(t,
		"Pick up point not found",
		"Check that error is returned and mapped to NotFound",
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

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		_, err = self.service.ListPickUpPointWorkingHours(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPickUpPointServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointServiceIntegrationTestSuite))
}

func TestPickUpPointPhotoServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointPhotoServiceIntegrationTestSuite))
}

func TestPickUpPointWorkingHoursServiceIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointWorkingHoursServiceIntegrationTestSuite))
}

