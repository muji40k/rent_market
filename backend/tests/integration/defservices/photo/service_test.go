package photo_test

import (
	"rent_service/builders/misc/generator"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	"rent_service/builders/mothers/test/repository/psql"
	mdefservices "rent_service/builders/mothers/test/service/defservices"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/photo"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
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

func MapPhoto(value *models.Photo, f func(string) string) photo.Photo {
	return photo.Photo{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Href:        f(value.Path),
		Date:        date.New(value.Date),
	}
}

func MapTempPhoto(value *models.TempPhoto) photo.TempPhoto {
	return photo.TempPhoto{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Date:        date.New(value.Create),
	}
}

func ComparePhoto(e photo.Photo, a photo.Photo) bool {
	return e.Id == a.Id &&
		e.Mime == a.Mime &&
		e.Placeholder == a.Placeholder &&
		e.Description == a.Description &&
		e.Href == a.Href &&
		psqlcommon.CompareTimeMicro(e.Date.Time, a.Date.Time)
}

func CompareTempPhoto(e photo.TempPhoto, a photo.TempPhoto) bool {
	return e.Id == a.Id &&
		e.Mime == a.Mime &&
		e.Placeholder == a.Placeholder &&
		e.Description == a.Description &&
		psqlcommon.CompareTimeMicro(e.Date.Time, a.Date.Time)
}

type PhotoServiceIntegrationTestSuite struct {
	suite.Suite
	service  photo.IService
	sContext defservices.Context
	rContext psqlcommon.Context
}

func (self *PhotoServiceIntegrationTestSuite) BeforeAll(t provider.T) {
	self.rContext.SetUp(t)
	self.sContext.SetUp(t, self.rContext.Factory.ToFactories())
}

func (self *PhotoServiceIntegrationTestSuite) AfterAll(t provider.T) {
	self.sContext.TearDown(t)
	self.rContext.TearDown(t)
}

func (self *PhotoServiceIntegrationTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"Integration tests",
		"Default services with PSQL repository",
		"Photo service",
	)

	t.WithNewStep("Clear database", func(sCtx provider.StepCtx) {
		self.rContext.Inserter.ClearDB()
	})

	t.WithNewStep("Clear photo registry", func(sCtx provider.StepCtx) {
		self.sContext.PhotoRegistry.Clear()
	})

	t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
		self.service = self.sContext.Factory.CreatePhotoService()
	})
}

var describeCreateTempPhoto = testcommon.MethodDescriptor(
	"CreateTempPhoto",
	"Create entry for temporary photo",
)

var describeUploadTempPhoto = testcommon.MethodDescriptor(
	"UploadTempPhoto",
	"Upload temporary photo",
)

var describeGetTempPhoto = testcommon.MethodDescriptor(
	"GetTempPhoto",
	"Get meta information of the temporary photo",
)

var describeGetPhoto = testcommon.MethodDescriptor(
	"GetPhoto",
	"Get meta information of the photo",
)

func (self *PhotoServiceIntegrationTestSuite) TestCreateTempPhotoPositive(t provider.T) {
	var (
		user models.User
		form photo.Description
	)

	describeCreateTempPhoto(t,
		"Simple create temp photo test",
		"Check that new photo is created and id is returned",
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

		var refphoto models.TempPhoto
		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = photo.Description{
				Mime:        refphoto.Mime,
				Placeholder: refphoto.Placeholder,
				Description: refphoto.Description,
			}
		})
	})

	// Act
	var result uuid.UUID
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = self.service.CreateTempPhoto(token.Token(user.Token), form)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().NotEqual(uuid.UUID{}, result, "Id is not empty")
}

func (self *PhotoServiceIntegrationTestSuite) TestCreateTempPhotoUnknownMime(t provider.T) {
	var (
		user models.User
		form photo.Description
	)

	describeCreateTempPhoto(t,
		"Unknown mime type supplied",
		"Check that user can't create temporaray photo with unknown (unsupported) mime type",
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

		var refphoto models.TempPhoto
		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty(
					"Test",
					nullable.None[time.Time](),
				).WithMime("unknown/mime/for/test").Build(),
			)
		})

		t.WithNewStep("Create form", func(sCtx provider.StepCtx) {
			form = photo.Description{
				Mime:        refphoto.Mime,
				Placeholder: refphoto.Placeholder,
				Description: refphoto.Description,
			}
		})
	})

	// Act
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			_, err = self.service.CreateTempPhoto(token.Token(user.Token), form)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var umerr photo.ErrorUnsupportedMime

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &umerr, "Error is UnsupportedMime")
}

func (self *PhotoServiceIntegrationTestSuite) TestUploadTempPhotoPositive(t provider.T) {
	var (
		id      uuid.UUID
		user    models.User
		content []byte
	)

	describeUploadTempPhoto(t,
		"Simple upload temp photo test",
		"Check that photo is uploaded",
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
			Add(psql.GeneratorStepNewValue(t, "photo",
				func() (models.TempPhoto, uuid.UUID) {
					p := models_om.TempPhotoExampleEmpty(
						"Test",
						nullable.None[time.Time](),
					).Build()
					id = p.Id
					return p, p.Id
				},
				self.rContext.Inserter.InsertTempPhoto,
			), 1).
			Generate().
			Finish()

		t.WithNewStep("Create photo content", func(sCtx provider.StepCtx) {
			content = models_om.ImagePNGContent(nullable.None[int]())
			sCtx.WithNewAttachment("content", allure.Png, content)
		})
	})

	// Act
	var err error

	t.WithNewStep("Upload temp photo",
		func(sCtx provider.StepCtx) {
			err = self.service.UploadTempPhoto(
				token.Token(user.Token),
				id,
				content,
			)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *PhotoServiceIntegrationTestSuite) TestUploadTempNotFound(t provider.T) {
	var (
		id      uuid.UUID
		user    models.User
		content []byte
	)

	describeUploadTempPhoto(t,
		"Photo not found",
		"Check that user can't upload content of unknown photo (NotFound retuned)",
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

		t.WithNewStep("Create unknown id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "photoId", uuidgen.Generate())
		})

		t.WithNewStep("Create photo content", func(sCtx provider.StepCtx) {
			content = models_om.ImagePNGContent(nullable.Some(100))
			sCtx.WithNewAttachment("content", allure.Png, content)
		})
	})

	// Act
	var err error

	t.WithNewStep("Upload temp photo",
		func(sCtx provider.StepCtx) {
			err = self.service.UploadTempPhoto(
				token.Token(user.Token),
				id,
				content,
			)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *PhotoServiceIntegrationTestSuite) TestGetTempPhotoPositive(t provider.T) {
	var (
		reference photo.TempPhoto
		user      models.User
	)

	describeGetTempPhoto(t,
		"Get temporary photo by id",
		"Check that photo is returned without error",
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
			Add(psql.GeneratorStepNewValue(t, "photo",
				func() (models.TempPhoto, uuid.UUID) {
					p := models_om.TempPhotoExampleEmpty(
						"Test",
						nullable.None[time.Time](),
					).Build()

					reference = MapTempPhoto(&p)

					return p, p.Id
				},
				self.rContext.Inserter.InsertTempPhoto,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var result photo.TempPhoto
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetTempPhoto(
				token.Token(user.Token),
				reference.Id,
			)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", reference.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[photo.TempPhoto](t).EqualFunc(
		CompareTempPhoto, reference, result,
		"Returned value is equal to reference",
	)
}

func (self *PhotoServiceIntegrationTestSuite) TestGetTempPhotoNotFound(t provider.T) {
	var (
		id   uuid.UUID
		user models.User
	)

	describeGetTempPhoto(t,
		"Photo not found",
		"Check that error NotFound is returned",
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

		t.WithNewStep("Create uknown identifier", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id",
				uuidgen.Generate(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetTempPhoto(token.Token(user.Token), id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *PhotoServiceIntegrationTestSuite) TestGetPhotoPositive(t provider.T) {
	var reference photo.Photo

	describeGetPhoto(t,
		"Get photo by id",
		"Check that photo is returned without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		generator.NewGeneratorGroup().
			Add(psql.GeneratorStepNewValue(t, "photo",
				func() (models.Photo, uuid.UUID) {
					path := self.sContext.PhotoRegistry.SavePhoto(
						models_om.ImagePNGContent(nullable.None[int]()),
					)
					p := models_om.PhotoExample(
						"Test",
						nullable.None[time.Time](),
					).WithPath(path).Build()

					reference = MapPhoto(&p, mdefservices.DefaultPathConverter)

					return p, p.Id
				},
				self.rContext.Inserter.InsertPhoto,
			), 1).
			Generate().
			Finish()
	})

	// Act
	var result photo.Photo
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = self.service.GetPhoto(reference.Id)
		},
		allure.NewParameter("photoId", reference.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	testcommon.Require[photo.Photo](t).EqualFunc(
		ComparePhoto, reference, result,
		"Returned value is equal to reference",
	)
}

func (self *PhotoServiceIntegrationTestSuite) TestGetPhotoNotFound(t provider.T) {
	var id uuid.UUID

	describeGetPhoto(t,
		"Photo not found",
		"Check that error NotFound is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create uknown identifier", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id",
				uuidgen.Generate(),
			)
		})
	})

	// Act
	var err error

	t.WithNewStep("Get photo",
		func(sCtx provider.StepCtx) {
			_, err = self.service.GetPhoto(id)
		},
		allure.NewParameter("photoId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPhotoServiceIntegrationTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PhotoServiceIntegrationTestSuite))
}

