package pickuppoint_test

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/pickuppoint"
	"rent_service/internal/logic/services/interfaces/pickuppoint"
	"rent_service/internal/logic/services/types/day"
	"rent_service/internal/logic/services/types/daytime"
	"rent_service/internal/misc/types/collection"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"testing"
	"time"

	rpickuppoint "rent_service/internal/repository/implementation/mock/pickuppoint"

	pickuppoint_pmock "rent_service/internal/repository/context/mock/pickuppoint"

	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

func GetService(ctrl *gomock.Controller, f func(repo *rpickuppoint.MockIRepository)) pickuppoint.IService {
	repo := rpickuppoint.NewMockIRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.New(pickuppoint_pmock.New(repo))
}

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

type PickUpPointServiceTestSuite struct {
	suite.Suite
}

func (self *PickUpPointServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Pick Up Point service",
	)
}

var describeListPickUpPoints = testcommon.MethodDescriptor(
	"ListPickUpPoints",
	"Get list of all pick up points",
)

var describeGetPickUpPointById = testcommon.MethodDescriptor(
	"GetPickUpPointById",
	"Get pick up point by identifier",
)

func (self *PickUpPointServiceTestSuite) TestListPicUpPointsPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IService

	var pups []models.PickUpPoint
	var reference []pickuppoint.PickUpPoint

	describeListPickUpPoints(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create reference periods", func(sCtx provider.StepCtx) {
			pups = testcommon.AssignParameter(sCtx, "Pick up points",
				collect.Do(
					models_om.PickUpPointExample("1"),
					models_om.PickUpPointExample("2"),
					models_om.PickUpPointExample("3"),
					models_om.PickUpPointExample("4"),
					models_om.PickUpPointExample("5"),
				),
			)

			reference = collection.Collect(
				collection.MapIterator(
					MapPickUpPoint, collection.SliceIterator(pups),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rpickuppoint.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(collection.SliceCollection(pups), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[pickuppoint.PickUpPoint]
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		result, err = service.ListPickUpPoints()
	})

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointServiceTestSuite) TestListPickUpPointsInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IService

	describeListPickUpPoints(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rpickuppoint.MockIRepository) {
				repo.EXPECT().GetAll().
					Return(nil, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		_, err = service.ListPickUpPoints()
	})

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "No error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func (self *PickUpPointServiceTestSuite) TestGetPickUpPointByIdPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IService

	var pup models.PickUpPoint
	var reference pickuppoint.PickUpPoint

	describeGetPickUpPointById(t,
		"Simple return by id test",
		"Checks that get method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create reference period", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "Pick up point",
				models_om.PickUpPointExample("1").Build(),
			)
			reference = MapPickUpPoint(&pup)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rpickuppoint.MockIRepository) {
				repo.EXPECT().GetById(pup.Id).
					Return(pup, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result pickuppoint.PickUpPoint
	var err error

	t.WithNewStep(
		"Get pick up point",
		func(sCtx provider.StepCtx) {
			result, err = service.GetPickUpPointById(pup.Id)
		},
		allure.NewParameter("pickUpPointId", pup.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Returned value matches reference")
}

func (self *PickUpPointServiceTestSuite) TestGetPickUpPointByIdNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IService

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

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetService(ctrl, func(repo *rpickuppoint.MockIRepository) {
				repo.EXPECT().GetById(id).
					Return(models.PickUpPoint{}, repo_errors.NotFound("pick_up_point_id")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep(
		"Get pick up point",
		func(sCtx provider.StepCtx) {
			_, err = service.GetPickUpPointById(id)
		},
		allure.NewParameter("pickUpPointId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func GetPhotoService(ctrl *gomock.Controller, f func(repo *rpickuppoint.MockIPhotoRepository)) pickuppoint.IPhotoService {
	repo := rpickuppoint.NewMockIPhotoRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.NewPhoto(pickuppoint_pmock.NewPhoto(repo))
}

type PickUpPointPhotoServiceTestSuite struct {
	suite.Suite
}

func (self *PickUpPointPhotoServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Pick Up Point photo service",
	)
}

var describeListPickUpPointPhotos = testcommon.MethodDescriptor(
	"ListPickUpPointPhotos",
	"Get list of all pick up point photos",
)

func (self *PickUpPointPhotoServiceTestSuite) TestListPickUpPointPhotosPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IPhotoService

	var (
		pup models.PickUpPoint
		ids []uuid.UUID
	)

	describeListPickUpPointPhotos(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointRandomId().Build(),
			)
		})

		t.WithNewStep("Create reference photo ids", func(sCtx provider.StepCtx) {
			ids = testcommon.AssignParameter(sCtx, "photoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID {
							return uuidgen.Generate()
						},
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetPhotoService(ctrl, func(repo *rpickuppoint.MockIPhotoRepository) {
				repo.EXPECT().GetById(pup.Id).
					Return(collection.SliceCollection(ids), nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[uuid.UUID]
	var err error

	t.WithNewStep("Get all pick up point photos", func(sCtx provider.StepCtx) {
		result, err = service.ListPickUpPointPhotos(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(ids, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointPhotoServiceTestSuite) TestListPickUpPointPhotosInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IPhotoService

	var pup models.PickUpPoint

	describeListPickUpPointPhotos(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "Pick up point",
				models_om.PickUpPointRandomId().Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetPhotoService(ctrl, func(repo *rpickuppoint.MockIPhotoRepository) {
				repo.EXPECT().GetById(pup.Id).
					Return(nil, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		_, err = service.ListPickUpPointPhotos(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func GetWorkingHoursService(ctrl *gomock.Controller, f func(repo *rpickuppoint.MockIWorkingHoursRepository)) pickuppoint.IWorkingHoursService {
	repo := rpickuppoint.NewMockIWorkingHoursRepository(ctrl)

	if nil != f {
		f(repo)
	}

	return service.NewWorkingHours(pickuppoint_pmock.NewWorkingHours(repo))
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

type PickUpPointWorkingHoursServiceTestSuite struct {
	suite.Suite
}

func (self *PickUpPointWorkingHoursServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Pick Up Point working hours service",
	)
}

var describeListPickUpPointWorkingHours = testcommon.MethodDescriptor(
	"ListPickUpPointPhotos",
	"Get pick up point working hours",
)

func (self *PickUpPointWorkingHoursServiceTestSuite) TestListPickUpPointWorkingHoursPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IWorkingHoursService

	var (
		pup       models.PickUpPoint
		wh        models.PickUpPointWorkingHours
		reference []pickuppoint.WorkingHours
	)

	describeListPickUpPointWorkingHours(t,
		"Simple return all test",
		"Checks that list method calls repository and doesn't return error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "pickUpPoint",
				models_om.PickUpPointRandomId().Build(),
			)
		})

		t.WithNewStep("Create reference working hours", func(sCtx provider.StepCtx) {
			wh = testcommon.AssignParameter(sCtx, "Working hours",
				models_om.PickUpPointWorkingHours(
					pup.Id, collect.Do(
						models_om.WorkingHoursWeek(8*time.Hour, 20*time.Hour)...,
					)...,
				).Build(),
			)
			reference = MapWorkingHours(&wh)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetWorkingHoursService(ctrl, func(repo *rpickuppoint.MockIWorkingHoursRepository) {
				repo.EXPECT().GetById(pup.Id).
					Return(wh, nil).
					MinTimes(1)
			})
		})
	})

	// Act
	var result collection.Collection[pickuppoint.WorkingHours]
	var err error

	t.WithNewStep("Get all pick up point working hours", func(sCtx provider.StepCtx) {
		result, err = service.ListPickUpPointWorkingHours(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All values must be returned",
	)
}

func (self *PickUpPointWorkingHoursServiceTestSuite) TestListPickUpPointWorkingHoursInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service pickuppoint.IWorkingHoursService

	var pup models.PickUpPoint

	describeListPickUpPointWorkingHours(t,
		"Internal error mapping",
		"Check that error is returned and mapped to Internal:DataAccess",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up point", func(sCtx provider.StepCtx) {
			pup = testcommon.AssignParameter(sCtx, "Pick up point",
				models_om.PickUpPointRandomId().Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = GetWorkingHoursService(ctrl, func(repo *rpickuppoint.MockIWorkingHoursRepository) {
				repo.EXPECT().GetById(pup.Id).
					Return(models.PickUpPointWorkingHours{}, errors.New("Some internal error")).
					MinTimes(1)
			})
		})
	})

	// Act
	var err error

	t.WithNewStep("Get all pick up points", func(sCtx provider.StepCtx) {
		_, err = service.ListPickUpPointWorkingHours(pup.Id)
	}, allure.NewParameter("pickUpPointId", pup.Id))

	// Assert
	var ierr cmnerrors.ErrorInternal
	var daerr cmnerrors.ErrorDataAccess

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &ierr, "Error is Internal")
	t.Require().ErrorAs(ierr, &daerr, "Error is DataAccess")
}

func TestPickUpPointServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointServiceTestSuite))
}

func TestPickUpPointPhotoServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointPhotoServiceTestSuite))
}

func TestPickUpPointWorkingHoursServiceTestSuite(t *testing.T) {
	suite.RunSuite(t, new(PickUpPointWorkingHoursServiceTestSuite))
}

