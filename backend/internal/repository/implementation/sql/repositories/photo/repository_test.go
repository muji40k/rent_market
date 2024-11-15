package photo_test

import (
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/errors/cmnerrors"
	"rent_service/internal/repository/interfaces/photo"
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

type PhotoRepositoryTestSuite struct {
	suite.Suite
	repo photo.IRepository
	psqlcommon.Context
}

func (self *PhotoRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *PhotoRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *PhotoRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Photo repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePhotoRepository()
	})
}

var describeCreate = testcommon.MethodDescriptor(
	"Create",
	"Create photo",
)

var describeGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get photo by id",
)

func ComparePhoto(expected models.Photo, actual models.Photo) bool {
	return expected.Id == actual.Id &&
		expected.Path == actual.Path &&
		expected.Mime == actual.Mime &&
		expected.Placeholder == actual.Placeholder &&
		expected.Description == actual.Description &&
		psqlcommon.CompareTimeMicro(expected.Date, actual.Date)
}

func CheckCreated(expected models.Photo, actual models.Photo) bool {
	return uuid.UUID{} != actual.Id &&
		expected.Path == actual.Path &&
		expected.Mime == actual.Mime &&
		expected.Placeholder == actual.Placeholder &&
		expected.Description == actual.Description &&
		psqlcommon.CompareTimeMicro(expected.Date, actual.Date)
}

func (self *PhotoRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var reference models.Photo

	describeCreate(t,
		"Simple create test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.PhotoExample("Test", nullable.None[time.Time]()).
					WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result models.Photo
	var err error

	t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.Photo](t).EqualFunc(
		CheckCreated, reference, result, "Same photo with non null uuid",
	)
}

func (self *PhotoRepositoryTestSuite) TestCreateDuplicate(t provider.T) {
	var reference models.Photo

	describeCreate(t,
		"Photo with path already exists",
		"Checks that error is returned and is Duplicate",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert photo", func(sCtx provider.StepCtx) {
			builder := models_om.PhotoExample("Test", nullable.None[time.Time]())
			created := testcommon.AssignParameter(sCtx, "created",
				builder.Build(),
			)
			reference = testcommon.AssignParameter(sCtx, "photo",
				builder.WithId(uuid.UUID{}).Build(),
			)
			self.Inserter.InsertPhoto(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	var derr cmnerrors.ErrorDuplicate

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &derr, "Error is Duplicate")
}

func (self *PhotoRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var reference models.Photo

	describeGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.PhotoExample("Test", nullable.None[time.Time]()).
					Build(),
			)
			self.Inserter.InsertPhoto(&reference)
		})
	})

	// Act
	var result models.Photo
	var err error

	t.WithNewStep("Get photo by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("photoId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.Photo](t).EqualFunc(
		ComparePhoto, reference, result, "Same photo value",
	)
}

func (self *PhotoRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeGetById(t,
		"Photo not found",
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

	t.WithNewStep("Get photo by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("photoId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

type TempPhotoRepositoryTestSuite struct {
	suite.Suite
	repo photo.ITempRepository
	psqlcommon.Context
}

func (self *TempPhotoRepositoryTestSuite) BeforeAll(t provider.T) {
	self.Context.SetUp(t)
}

func (self *TempPhotoRepositoryTestSuite) AfterAll(t provider.T) {
	self.Context.TearDown(t)
}

func (self *TempPhotoRepositoryTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"PSQLRepositories",
		"PSQL repository implementation",
		"Temp Photo repository",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.Inserter.ClearDB()
	})

	t.WithNewStep("Create repository", func(sCtx provider.StepCtx) {
		self.repo = self.Factory.CreatePhotoTempRepository()
	})
}

var describeTempCreate = testcommon.MethodDescriptor(
	"Create",
	"Create temporary photo",
)

var describeTempUpdate = testcommon.MethodDescriptor(
	"Update",
	"Update temporary photo",
)

var describeTempGetById = testcommon.MethodDescriptor(
	"GetById",
	"Get temporary photo by id",
)

var describeTempRemove = testcommon.MethodDescriptor(
	"Remove",
	"Remove temporary photo",
)

func CompareTempPhoto(expected models.TempPhoto, actual models.TempPhoto) bool {
	return expected.Id == actual.Id &&
		((nil == expected.Path && nil == actual.Path) ||
			(*expected.Path == *actual.Path)) &&
		expected.Mime == actual.Mime &&
		expected.Placeholder == actual.Placeholder &&
		expected.Description == actual.Description &&
		psqlcommon.CompareTimeMicro(expected.Create, actual.Create)
}

func CheckTempCreated(expected models.TempPhoto, actual models.TempPhoto) bool {
	return uuid.UUID{} != actual.Id &&
		((nil == expected.Path && nil == actual.Path) ||
			(*expected.Path == *actual.Path)) &&
		expected.Mime == actual.Mime &&
		expected.Placeholder == actual.Placeholder &&
		expected.Description == actual.Description &&
		psqlcommon.CompareTimeMicro(expected.Create, actual.Create)
}

func (self *TempPhotoRepositoryTestSuite) TestCreatePositive(t provider.T) {
	var reference models.TempPhoto

	describeTempCreate(t,
		"Simple create test",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create temporary photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty("Test", nullable.None[time.Time]()).
					WithId(uuid.UUID{}).Build(),
			)
		})
	})

	// Act
	var result models.TempPhoto
	var err error

	t.WithNewStep("Create temporary photo", func(sCtx provider.StepCtx) {
		result, err = self.repo.Create(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.TempPhoto](t).EqualFunc(
		CheckTempCreated, reference, result, "Same photo with non null uuid",
	)
}

func (self *TempPhotoRepositoryTestSuite) TestCreateDuplicate(t provider.T) {
	var reference models.TempPhoto

	describeCreate(t,
		"Temporary Photo with path already exists",
		"Checks that error is returned and is Duplicate",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert photo", func(sCtx provider.StepCtx) {
			builder := models_om.TempPhotoExampleUploaded("Test", nullable.None[time.Time]())
			created := testcommon.AssignParameter(sCtx, "created",
				builder.Build(),
			)
			reference = testcommon.AssignParameter(sCtx, "photo",
				builder.WithId(uuid.UUID{}).Build(),
			)
			self.Inserter.InsertTempPhoto(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
		_, err = self.repo.Create(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	var derr cmnerrors.ErrorDuplicate

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &derr, "Error is Duplicate")
}

func (self *TempPhotoRepositoryTestSuite) TestUpdatePositive(t provider.T) {
	var reference models.TempPhoto

	describeTempUpdate(t,
		"Simple update test (photo upload)",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert temporary photo", func(sCtx provider.StepCtx) {
			builder := models_om.TempPhotoExampleEmpty("Test", nullable.None[time.Time]())
			created := builder.Build()
			reference = testcommon.AssignParameter(sCtx, "photo",
				builder.WithPath(nullable.Some("/test/path")).Build(),
			)
			self.Inserter.InsertTempPhoto(&created)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update temporary photo", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *TempPhotoRepositoryTestSuite) TestUpdateNotFound(t provider.T) {
	var reference models.TempPhoto

	describeTempUpdate(t,
		"Temp photo not found during update",
		"Checks that error is returned and is NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create temporary photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleUploaded("Test", nullable.None[time.Time]()).
					Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Update temporary photo", func(sCtx provider.StepCtx) {
		err = self.repo.Update(reference)
	}, allure.NewParameter("photo", reference))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *TempPhotoRepositoryTestSuite) TestGetByIdPositive(t provider.T) {
	var reference models.TempPhoto

	describeTempGetById(t,
		"Simple return test",
		"Checks that method returns value without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty("Test", nullable.None[time.Time]()).
					Build(),
			)
			self.Inserter.InsertTempPhoto(&reference)
		})
	})

	// Act
	var result models.TempPhoto
	var err error

	t.WithNewStep("Get photo by id", func(sCtx provider.StepCtx) {
		result, err = self.repo.GetById(reference.Id)
	}, allure.NewParameter("photoId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Assert[models.TempPhoto](t).EqualFunc(
		CompareTempPhoto, reference, result, "Same photo value",
	)
}

func (self *TempPhotoRepositoryTestSuite) TestGetByIdNotFound(t provider.T) {
	var id uuid.UUID

	describeTempGetById(t,
		"Photo not found",
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

	t.WithNewStep("Get photo by id", func(sCtx provider.StepCtx) {
		_, err = self.repo.GetById(id)
	}, allure.NewParameter("photoId", id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *TempPhotoRepositoryTestSuite) TestRemovePositive(t provider.T) {
	var reference models.TempPhoto

	describeTempRemove(t,
		"Simple remove test (photo moved to persistent storage)",
		"Checks that no error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create and insert temporary photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleUploaded("Test", nullable.None[time.Time]()).
					Build(),
			)
			self.Inserter.InsertTempPhoto(&reference)
		})
	})

	// Act
	var err error

	t.WithNewStep("Remove temporary photo", func(sCtx provider.StepCtx) {
		err = self.repo.Remove(reference.Id)
	}, allure.NewParameter("photoId", reference.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *TempPhotoRepositoryTestSuite) TestRemoveNotFound(t provider.T) {
	var reference models.TempPhoto

	describeTempRemove(t,
		"Temp photo not found during remove",
		"Checks that error is returned and is NotFound",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create temporary photo", func(sCtx provider.StepCtx) {
			reference = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleUploaded("Test", nullable.None[time.Time]()).
					Build(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Remove temporary photo", func(sCtx provider.StepCtx) {
		err = self.repo.Remove(reference.Id)
	}, allure.NewParameter("photoId", reference.Id))

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPhotoRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PhotoRepositoryTestSuite))
}

func TestTempPhotoRepositoryTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TempPhotoRepositoryTestSuite))
}

