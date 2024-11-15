package photo_test

import (
	"errors"
	"reflect"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	auth_builder "rent_service/builders/test/mock/defmisc/authenticator"
	registry_builder "rent_service/builders/test/mock/defmisc/photoregistry"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry/storages/mock"
	service "rent_service/internal/logic/services/implementations/defservices/services/photo"
	"rent_service/internal/logic/services/interfaces/photo"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	photo_pmock "rent_service/internal/repository/context/mock/photo"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	rphoto "rent_service/internal/repository/implementation/mock/photo"
	ruser "rent_service/internal/repository/implementation/mock/user"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

type ServiceBuilder struct {
	authenticator authenticator.IAuthenticator
	registry      photoregistry.IRegistry
	photo         *rphoto.MockIRepository
	temp          *rphoto.MockITempRepository
}

func NewServiceBuilder(ctrl *gomock.Controller) *ServiceBuilder {
	return &ServiceBuilder{
		photo: rphoto.NewMockIRepository(ctrl),
		temp:  rphoto.NewMockITempRepository(ctrl),
	}
}

func (self *ServiceBuilder) WithAuthenticator(authenticator authenticator.IAuthenticator) *ServiceBuilder {
	self.authenticator = authenticator
	return self
}

func (self *ServiceBuilder) WithPhotoRegistry(registry photoregistry.IRegistry) *ServiceBuilder {
	self.registry = registry
	return self
}

func (self *ServiceBuilder) WithPhotoRepository(f func(repo *rphoto.MockIRepository)) *ServiceBuilder {
	f(self.photo)
	return self
}

func (self *ServiceBuilder) WithPhotoTempRepository(f func(repo *rphoto.MockITempRepository)) *ServiceBuilder {
	f(self.temp)
	return self
}

func (self *ServiceBuilder) GetService() photo.IService {
	return service.New(
		self.authenticator,
		self.registry,
		photo_pmock.New(self.photo),
		photo_pmock.NewTemp(self.temp),
	)
}

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

type PhotoServiceTestSuite struct {
	suite.Suite
}

func (self *PhotoServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Photo service",
	)
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

type TempPhotoMatcher struct {
	value *models.TempPhoto
}

func (self TempPhotoMatcher) Matches(x any) bool {
	if reflect.TypeOf(models.TempPhoto{}) != reflect.TypeOf(x) {
		return false
	}

	xc := reflect.ValueOf(x).Interface().(models.TempPhoto)

	return gomock.Eq(self.value.Id).Matches(xc.Id) &&
		gomock.Eq(self.value.Path).Matches(xc.Path) &&
		gomock.Eq(self.value.Mime).Matches(xc.Mime) &&
		gomock.Eq(self.value.Placeholder).Matches(xc.Placeholder) &&
		gomock.Eq(self.value.Description).Matches(xc.Description) &&
		gomock.Any().Matches(xc.Create)
}

func (self TempPhotoMatcher) String() string {
	return "Temp photo match with any create date"
}

func (self *PhotoServiceTestSuite) TestCreateTempPhotoPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto models.TempPhoto
		user     models.User
		form     photo.Description
	)

	describeCreateTempPhoto(t,
		"Simple create temp photo test",
		"Check that new photo is created and id is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
					toBeCreated := models.TempPhoto{
						Mime:        refphoto.Mime,
						Placeholder: refphoto.Placeholder,
						Description: refphoto.Description,
					}
					repo.EXPECT().Create(TempPhotoMatcher{&toBeCreated}).
						Return(refphoto, nil).
						Times(1)
				}).
				GetService()
		})
	})

	form = photo.Description{
		Mime:        refphoto.Mime,
		Placeholder: refphoto.Placeholder,
		Description: refphoto.Description,
	}

	// Act
	var result uuid.UUID
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = service.CreateTempPhoto(token.Token(user.Token), form)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(refphoto.Id, result, "Id is equal to reference")
}

func (self *PhotoServiceTestSuite) TestCreateTempPhotoUnknownMime(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto models.TempPhoto
		user     models.User
		form     photo.Description
	)

	describeCreateTempPhoto(t,
		"Unknown mime type supplied",
		"Check that user can't create temporaray photo with unknown (unsupported) mime type",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty(
					"Test",
					nullable.None[time.Time](),
				).WithMime("unknown/mime/for/test").Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							AnyTimes()
					}).
					Build(),
				).
				WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
					toBeCreated := models.TempPhoto{
						Mime:        refphoto.Mime,
						Placeholder: refphoto.Placeholder,
						Description: refphoto.Description,
					}
					repo.EXPECT().Create(TempPhotoMatcher{&toBeCreated}).
						Return(refphoto, nil).
						Times(0)
				}).
				GetService()
		})
	})

	form = photo.Description{
		Mime:        refphoto.Mime,
		Placeholder: refphoto.Placeholder,
		Description: refphoto.Description,
	}

	// Act
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			_, err = service.CreateTempPhoto(token.Token(user.Token), form)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var umerr photo.ErrorUnsupportedMime

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &umerr, "Error is UnsupportedMime")
}

func (self *PhotoServiceTestSuite) TestUploadTempPhotoPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto models.TempPhoto
		user     models.User
		content  []byte
	)

	describeUploadTempPhoto(t,
		"Simple upload temp photo test",
		"Check that photo is uploaded",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleUploaded(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create photo content", func(sCtx provider.StepCtx) {
			content = models_om.ImagePNGContent(nullable.None[int]())
			sCtx.WithNewAttachment("content", allure.Png, content)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoRegistry(registry_builder.New(ctrl).
					WithStorage(func(storage *mock_defregistry.MockIStorage) {
						storage.EXPECT().WriteTempData(content).
							Return(*refphoto.Path, nil).
							Times(1)
					}).
					WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
						var toBeSavedPhoto = models.TempPhoto{
							Id:          refphoto.Id,
							Path:        nil,
							Mime:        refphoto.Mime,
							Placeholder: refphoto.Placeholder,
							Description: refphoto.Description,
							Create:      refphoto.Create,
						}
						repo.EXPECT().GetById(refphoto.Id).
							Return(toBeSavedPhoto, nil).
							MinTimes(1)
						repo.EXPECT().Update(refphoto).
							Return(nil).
							Times(1)
					}).
					Build(),
				).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Upload temp photo",
		func(sCtx provider.StepCtx) {
			err = service.UploadTempPhoto(
				token.Token(user.Token),
				refphoto.Id,
				content,
			)
		},
		allure.NewParameter("photoId", user.Token),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *PhotoServiceTestSuite) TestUploadTempInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto models.TempPhoto
		user     models.User
		content  []byte
	)

	describeUploadTempPhoto(t,
		"Error during content write",
		"Check that internal error is returned and no records changed",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleUploaded(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create photo content", func(sCtx provider.StepCtx) {
			content = models_om.ImagePNGContent(nullable.Some(100))
			sCtx.WithNewAttachment("content", allure.Png, content)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoRegistry(registry_builder.New(ctrl).
					WithStorage(func(storage *mock_defregistry.MockIStorage) {
						storage.EXPECT().WriteTempData(content).
							Return("", errors.New("Some internale error")).
							Times(1)
					}).
					WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
						var toBeSavedPhoto = models.TempPhoto{
							Id:          refphoto.Id,
							Path:        nil,
							Mime:        refphoto.Mime,
							Placeholder: refphoto.Placeholder,
							Description: refphoto.Description,
							Create:      refphoto.Create,
						}
						repo.EXPECT().GetById(refphoto.Id).
							Return(toBeSavedPhoto, nil).
							MinTimes(1)
						repo.EXPECT().Update(gomock.Any()).Times(0)
					}).
					Build(),
				).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Upload temp photo",
		func(sCtx provider.StepCtx) {
			err = service.UploadTempPhoto(
				token.Token(user.Token),
				refphoto.Id,
				content,
			)
		},
		allure.NewParameter("photoId", user.Token),
	)

	// Assert
	var ierr cmnerrors.ErrorInternal

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
}

func (self *PhotoServiceTestSuite) TestGetTempPhotoPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto  models.TempPhoto
		reference photo.TempPhoto
		user      models.User
	)

	describeGetTempPhoto(t,
		"Get temporary photo by id",
		"Check that photo is returned without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.TempPhotoExampleEmpty(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
			reference = MapTempPhoto(&refphoto)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
					repo.EXPECT().GetById(refphoto.Id).
						Return(refphoto, nil).
						Times(1)
				}).
				GetService()
		})
	})

	// Act
	var result photo.TempPhoto
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = service.GetTempPhoto(token.Token(user.Token), refphoto.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", refphoto.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value is equal to reference")
}

func (self *PhotoServiceTestSuite) TestGetTempPhotoNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

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
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create uknown identifier", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id",
				uuidgen.Generate(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(auth_builder.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoTempRepository(func(repo *rphoto.MockITempRepository) {
					repo.EXPECT().GetById(id).
						Return(models.TempPhoto{}, repo_errors.NotFound("temp_photo_id")).
						Times(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			_, err = service.GetTempPhoto(token.Token(user.Token), id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("photoId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func convertPath(path string) string {
	return "http://some.storage.url.com/remarkable/path/for/substitute"
}

func (self *PhotoServiceTestSuite) TestGetPhotoPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

	var (
		refphoto  models.Photo
		reference photo.Photo
	)

	describeGetPhoto(t,
		"Get photo by id",
		"Check that photo is returned without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create photo", func(sCtx provider.StepCtx) {
			refphoto = testcommon.AssignParameter(sCtx, "photo",
				models_om.PhotoExample(
					"Test",
					nullable.None[time.Time](),
				).Build(),
			)
			reference = MapPhoto(&refphoto, convertPath)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithPhotoRegistry(registry_builder.New(ctrl).
					WithStorage(func(storage *mock_defregistry.MockIStorage) {
						storage.EXPECT().ConvertPath(refphoto.Path).
							DoAndReturn(convertPath).
							MinTimes(1)
					}).
					Build(),
				).
				WithPhotoRepository(func(repo *rphoto.MockIRepository) {
					repo.EXPECT().GetById(refphoto.Id).
						Return(refphoto, nil).
						Times(1)
				}).
				GetService()
		})
	})

	// Act
	var result photo.Photo
	var err error

	t.WithNewStep("Create temp photo",
		func(sCtx provider.StepCtx) {
			result, err = service.GetPhoto(refphoto.Id)
		},
		allure.NewParameter("photoId", refphoto.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value is equal to reference")
}

func (self *PhotoServiceTestSuite) TestGetPhotoNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service photo.IService

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

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithPhotoRegistry(registry_builder.New(ctrl).
					WithStorage(func(storage *mock_defregistry.MockIStorage) {
						storage.EXPECT().ConvertPath(gomock.Any()).
							DoAndReturn(convertPath).
							AnyTimes()
					}).
					Build(),
				).
				WithPhotoRepository(func(repo *rphoto.MockIRepository) {
					repo.EXPECT().GetById(id).
						Return(models.Photo{}, repo_errors.NotFound("photo_id")).
						Times(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Get photo",
		func(sCtx provider.StepCtx) {
			_, err = service.GetPhoto(id)
		},
		allure.NewParameter("photoId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestPhotoServiceTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(PhotoServiceTestSuite))
}

