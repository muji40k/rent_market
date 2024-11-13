package storage_test

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/records"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/storage"
	"rent_service/internal/logic/services/interfaces/storage"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	"time"

	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"testing"

	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"

	defaccess "rent_service/builders/test/mock/defmisc/access"
	defauth "rent_service/builders/test/mock/defmisc/authenticator"
	defautho "rent_service/builders/test/mock/defmisc/authorizer"

	rpickuppoint "rent_service/internal/repository/implementation/mock/pickuppoint"
	rrole "rent_service/internal/repository/implementation/mock/role"
	rstorage "rent_service/internal/repository/implementation/mock/storage"
	ruser "rent_service/internal/repository/implementation/mock/user"

	storage_pmock "rent_service/internal/repository/context/mock/storage"

	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	records_om "rent_service/builders/mothers/domain/records"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

type ServiceBuilder struct {
	authenticator authenticator.IAuthenticator
	pickUpPoint   access.IPickUpPoint
	storage       *rstorage.MockIRepository
}

func NewBuilder(ctrl *gomock.Controller) *ServiceBuilder {
	return &ServiceBuilder{
		storage: rstorage.NewMockIRepository(ctrl),
	}
}

func (self *ServiceBuilder) WithAuthenticator(auth authenticator.IAuthenticator) *ServiceBuilder {
	self.authenticator = auth
	return self
}

func (self *ServiceBuilder) WithPickUpPointAccess(acc access.IPickUpPoint) *ServiceBuilder {
	self.pickUpPoint = acc
	return self
}

func (self *ServiceBuilder) WithStorageRepository(f func(repo *rstorage.MockIRepository)) *ServiceBuilder {
	f(self.storage)
	return self
}

func (self *ServiceBuilder) GetService() storage.IService {
	return service.New(
		self.authenticator,
		storage_pmock.New(self.storage),
		self.pickUpPoint,
	)
}

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

type StorageServiceTestSuite struct {
	suite.Suite
}

func (self *StorageServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Storage service",
	)
}

var describeListStoragesByPickUpPoint = testcommon.MethodDescriptor(
	"ListStoragesByPickUpPoint",
	"Get list of storages",
)

var describeGetStorageByInstance = testcommon.MethodDescriptor(
	"GetStorageByInstance",
	"Get storage by identifier",
)

func (self *StorageServiceTestSuite) TestListStoragesByPickUpPointPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service storage.IService

	var (
		skUser    models.User
		sk        models.Storekeeper
		pup       models.PickUpPoint
		storages  []records.Storage
		reference []storage.Storage
	)

	describeListStoragesByPickUpPoint(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("Test").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
			sk = testcommon.AssignParameter(sCtx, "role",
				models_om.StorekeeperWithUserId(skUser.Id, pup.Id).Build(),
			)
		})

		t.WithNewStep("Create reference storages", func(sCtx provider.StepCtx) {
			storages = testcommon.AssignParameter(sCtx, "storages",
				collection.Collect(
					collection.MapIterator(
						func(i *int) records.Storage {
							return records_om.StorageActive(
								pup.Id,
								models_om.InstanceRandomId().Build().Id,
								nullable.None[time.Time](),
							).Build()
						},
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)

			reference = collection.Collect(
				collection.MapIterator(
					MapStorage, collection.SliceIterator(storages),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewBuilder(ctrl).
				WithAuthenticator(defauth.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(skUser.Token).
							Return(skUser, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPickUpPointAccess(defaccess.NewPickUpPoint(ctrl).
					WithAuthorizer(defautho.New(ctrl).
						WithAdministratorRepository(func(repo *rrole.MockIAdministratorRepository) {
							repo.EXPECT().GetByUserId(skUser.Id).
								Return(models.Administrator{}, repo_errors.NotFound("administrator_user_id")).
								MinTimes(1)
						}).
						WithStorekeeperRepository(func(repo *rrole.MockIStorekeeperRepository) {
							repo.EXPECT().GetByUserId(skUser.Id).
								Return(sk, nil).
								MinTimes(1)
						}).
						Build(),
					).
					WithPickUpPointRepository(func(repo *rpickuppoint.MockIRepository) {
						repo.EXPECT().GetById(pup.Id).
							Return(pup, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithStorageRepository(func(repo *rstorage.MockIRepository) {
					repo.EXPECT().GetActiveByPickUpPointId(pup.Id).
						Return(collection.SliceCollection(storages), nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var result collection.Collection[storage.Storage]
	var err error

	t.WithNewStep(
		"Get active storages for pick up point",
		func(sCtx provider.StepCtx) {
			result, err = service.ListStoragesByPickUpPoint(token.Token(skUser.Token), pup.Id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickUpPointId", pup.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *StorageServiceTestSuite) TestListStoragesByPickUpPointsUnauthorizedUser(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service storage.IService

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
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointExample("Test").Build(),
			)
		})

		t.WithNewStep("Create plain user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewBuilder(ctrl).
				WithAuthenticator(defauth.New(ctrl).
					WithUserRepository(func(repo *ruser.MockIRepository) {
						repo.EXPECT().GetByToken(user.Token).
							Return(user, nil).
							MinTimes(1)
					}).
					Build(),
				).
				WithPickUpPointAccess(defaccess.NewPickUpPoint(ctrl).
					WithAuthorizer(defautho.New(ctrl).
						WithAdministratorRepository(func(repo *rrole.MockIAdministratorRepository) {
							repo.EXPECT().GetByUserId(user.Id).
								Return(models.Administrator{}, repo_errors.NotFound("administrator_user_id")).
								MinTimes(1)
						}).
						WithStorekeeperRepository(func(repo *rrole.MockIStorekeeperRepository) {
							repo.EXPECT().GetByUserId(user.Id).
								Return(models.Storekeeper{}, repo_errors.NotFound("storekeeper_user_id")).
								MinTimes(1)
						}).
						Build(),
					).
					WithPickUpPointRepository(func(repo *rpickuppoint.MockIRepository) {
						repo.EXPECT().GetById(pup.Id).
							Return(pup, nil).
							MinTimes(1)
					}).
					Build(),
				).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep(
		"Get active storages by pick up point",
		func(sCtx provider.StepCtx) {
			_, err = service.ListStoragesByPickUpPoint(token.Token(user.Token), pup.Id)
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

func (self *StorageServiceTestSuite) TestGetStorageByInstancePositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service storage.IService

	var (
		instance  models.Instance
		strg      records.Storage
		reference storage.Storage
	)

	describeGetStorageByInstance(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceExample("Test", uuidgen.Generate()).Build(),
			)
		})

		t.WithNewStep("Create storage", func(sCtx provider.StepCtx) {
			strg = testcommon.AssignParameter(sCtx, "storage",
				records_om.StorageActive(
					uuidgen.Generate(),
					instance.Id,
					nullable.None[time.Time](),
				).Build(),
			)
			reference = MapStorage(&strg)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewBuilder(ctrl).
				WithStorageRepository(func(repo *rstorage.MockIRepository) {
					repo.EXPECT().GetActiveByInstanceId(instance.Id).
						Return(strg, nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var result storage.Storage
	var err error

	t.WithNewStep(
		"Get active storage for instance",
		func(sCtx provider.StepCtx) {
			result, err = service.GetStorageByInstance(instance.Id)
		},
		allure.NewParameter("instanceId", instance.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value matches reference")
}

func (self *StorageServiceTestSuite) TestGetStorageByInstanceNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service storage.IService

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

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewBuilder(ctrl).
				WithStorageRepository(func(repo *rstorage.MockIRepository) {
					repo.EXPECT().GetActiveByInstanceId(id).
						Return(records.Storage{}, repo_errors.NotFound("storage_instance_id")).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep(
		"Get active storage for instance",
		func(sCtx provider.StepCtx) {
			_, err = service.GetStorageByInstance(id)
		},
		allure.NewParameter("instanceId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestStorageServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(StorageServiceTestSuite))
}

