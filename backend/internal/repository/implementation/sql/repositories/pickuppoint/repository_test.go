package pickuppoint_test

import (
	models_b "rent_service/builders/domain/models"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/pickuppoint"
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

type PickUpPointRepositoryTestSuite struct {
	suite.Suite
	repo pickuppoint.IRepository
	psqlcommon.Context
}

func (self *PickUpPointRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PickUpPointRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PickUpPointRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Pick Up Point repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePickUpPointRepository()
	})
}

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get pick up point by id",
)

var describeGetAll = testcommon.MethodDescriptor(
	"GetAll",
	"Get all pick up points",
)

func (self *PickUpPointRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var reference models.PickUpPoint

	describeGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "company",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&reference)
		})
	})

	// Act
	var result models.PickUpPoint
	var err error

	t.WithNewStep("Get pick up point by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("pickUpPointId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same pick up point value")
}

func (self *PickUpPointRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Pick up point not found",
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

	t.WithNewStep("Get pick up point by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *PickUpPointRepositoryTestSuite) TestGetAllPositive(t provider.T) {
	var reference []models.PickUpPoint

	describeGetAll(t,
		"Simple return test",
		"Checks that method returns collection without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert pick up points", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "pickUpPoints",
				collect.DoN(5, collect.FmtWrap(models_om.PickUpPointExample)),
			)
			psql.BulkInsert(self.Inserter.InsertPickUpPoint, reference...)
		})
	})

	// Act
	var result collection.Collection[models.PickUpPoint]
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().ElementsMatch(reference, collection.Collect(result.Iter()),
		"Same company values")
}

func (self *PickUpPointRepositoryTestSuite) TestGetAllEmpty(t provider.T) {
	describeGetAll(t,
		"No pick up points",
		"Checks that method return empty collection withour error",
	)

	// Arrange
	// Empty

	// Act
	var result collection.Collection[models.PickUpPoint]
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetAll()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(uint(0), collection.Count(result.Iter()),
		"Collection is empty")
}

type PickUpPointPhotoRepositoryTestSuite struct {
	suite.Suite
	repo pickuppoint.IPhotoRepository
	psqlcommon.Context
}

func (self *PickUpPointPhotoRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PickUpPointPhotoRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PickUpPointPhotoRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Pick Up Point photo repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePickUpPointPhotoRepository()
	})
}

var describePhotoGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get pick up point photos by id",
)

func (self *PickUpPointPhotoRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var (
		reference []models.Photo
		pup       models.PickUpPoint
	)

	describePhotoGetById(t,
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

		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create and insert pick up point photos", func(sCtx provider.StepCtx) {
			psql.BulkInsert(self.Inserter.InsertPickUpPointPhoto,
				collection.Collect(
					collection.MapIterator(
						func(photo *models.Photo) psql.Photo {
							return *psql.NewPhoto(
								nullable.None[uuid.UUID](),
								pup.Id,
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

	t.WithNewStep("Get pick up point photos by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

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

func (self *PickUpPointPhotoRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describePhotoGetById(t,
		"Pick up point not found",
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

	t.WithNewStep("Get pick up point photos by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type PickUpPointWorkingHoursRepositoryTestSuite struct {
	suite.Suite
	repo pickuppoint.IWorkingHoursRepository
	psqlcommon.Context
}

func (self *PickUpPointWorkingHoursRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PickUpPointWorkingHoursRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PickUpPointWorkingHoursRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Pick Up Point working hours repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePickUpPointWorkingHoursRepository()
	})
}

var describeWorkingHoursGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get pick up point working hours by id",
)

func (self *PickUpPointWorkingHoursRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var (
		pup       models.PickUpPoint
		reference models.PickUpPointWorkingHours
	)

	describeWorkingHoursGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("1").Build(),
			)
			self.Inserter.InsertPickUpPoint(&pup)
		})

		t.WithNewStep("Create and insert pick up point working hours", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "workingHours",
				models_om.PickUpPointWorkingHours(
					pup.Id,
					collect.Do(
						models_om.WorkingHoursWeek(8*time.Hour, 20*time.Hour)...,
					)...,
				).Build(),
			)
			self.Inserter.InsertPickUpPointWorkingHours(&reference)
		})
	})

	// Act
	var result models.PickUpPointWorkingHours
	var err error

	t.WithNewStep("Get pick up point working hours by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Assert().Equal(reference, result, "Same photo ids")
}

func (self *PickUpPointWorkingHoursRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeWorkingHoursGetById(t,
		"Pick up point not found",
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

	t.WithNewStep("Get pick up point photos by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("pickUpPointId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPickUpPointRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointRepositoryTestSuite))
}

func TestPickUpPointPhotoRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointPhotoRepositoryTestSuite))
}

func TestPickUpPointWorkingHoursRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointWorkingHoursRepositoryTestSuite))
}

