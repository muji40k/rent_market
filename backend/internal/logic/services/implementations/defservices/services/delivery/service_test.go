package delivery_test

import (
	"errors"
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/logic/services/errors/cmnerrors"
	service "rent_service/internal/logic/services/implementations/defservices/services/delivery"
	"rent_service/internal/logic/services/interfaces/delivery"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/types/collection"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	"rent_service/misc/nullable"
	"rent_service/misc/testcommon"
	"testing"
	"time"

	access "rent_service/internal/logic/services/implementations/defservices/access/implementations/mock"
	authenticator "rent_service/internal/logic/services/implementations/defservices/authenticator/implementations/mock"
	photoregistry "rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/mock"
	istates "rent_service/internal/logic/services/implementations/defservices/states"
	states "rent_service/internal/logic/services/implementations/defservices/states/implementations/mock"

	rdelivery "rent_service/internal/repository/implementation/mock/delivery"
	rinstance "rent_service/internal/repository/implementation/mock/instance"
	rrent "rent_service/internal/repository/implementation/mock/rent"

	delivery_pmock "rent_service/internal/repository/context/mock/delivery"
	instance_pmock "rent_service/internal/repository/context/mock/instance"
	rent_pmock "rent_service/internal/repository/context/mock/rent"

	"rent_service/builders/misc/nullcommon"
	"rent_service/builders/misc/uuidgen"
	models_om "rent_service/builders/mothers/domain/models"
	requests_om "rent_service/builders/mothers/domain/requests"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
)

type ServiceBuilder struct {
	states         *states.MockIInstanceStateMachine
	authenticator  *authenticator.MockIAuthenticator
	registry       *photoregistry.MockIRegistry
	delivery       *rdelivery.MockIRepository
	photo          *rinstance.MockIPhotoRepository
	rentReq        *rrent.MockIRequestRepository
	instanceAcc    *access.MockIInstance
	pickUpPointAcc *access.MockIPickUpPoint
}

func NewServiceBuilder(ctrl *gomock.Controller) *ServiceBuilder {
	return &ServiceBuilder{
		states.NewMockIInstanceStateMachine(ctrl),
		authenticator.NewMockIAuthenticator(ctrl),
		photoregistry.NewMockIRegistry(ctrl),
		rdelivery.NewMockIRepository(ctrl),
		rinstance.NewMockIPhotoRepository(ctrl),
		rrent.NewMockIRequestRepository(ctrl),
		access.NewMockIInstance(ctrl),
		access.NewMockIPickUpPoint(ctrl),
	}
}

func (self *ServiceBuilder) WithInstanceStateMachine(f func(machine *states.MockIInstanceStateMachine)) *ServiceBuilder {
	f(self.states)
	return self
}

func (self *ServiceBuilder) WithAuthenticator(f func(authenticator *authenticator.MockIAuthenticator)) *ServiceBuilder {
	f(self.authenticator)
	return self
}

func (self *ServiceBuilder) WithPhotoRegistry(f func(registry *photoregistry.MockIRegistry)) *ServiceBuilder {
	f(self.registry)
	return self
}

func (self *ServiceBuilder) WithDeliveryRepository(f func(delivery *rdelivery.MockIRepository)) *ServiceBuilder {
	f(self.delivery)
	return self
}

func (self *ServiceBuilder) WithInstancePhotoRepository(f func(photo *rinstance.MockIPhotoRepository)) *ServiceBuilder {
	f(self.photo)
	return self
}

func (self *ServiceBuilder) WithRentRequestRepository(f func(rentReq *rrent.MockIRequestRepository)) *ServiceBuilder {
	f(self.rentReq)
	return self
}

func (self *ServiceBuilder) WithInstanceAccess(f func(instanceAcc *access.MockIInstance)) *ServiceBuilder {
	f(self.instanceAcc)
	return self
}

func (self *ServiceBuilder) WithPickUpPointAccess(f func(pickUpPointAcc *access.MockIPickUpPoint)) *ServiceBuilder {
	f(self.pickUpPointAcc)
	return self
}

func (self *ServiceBuilder) GetService() delivery.IService {
	return service.New(
		self.states,
		self.authenticator,
		self.registry,
		delivery_pmock.New(self.delivery),
		instance_pmock.NewPhoto(self.photo),
		rent_pmock.NewRequest(self.rentReq),
		self.instanceAcc,
		self.pickUpPointAcc,
	)
}

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

type DeliveryServiceTestSuite struct {
	suite.Suite
}

func (self *DeliveryServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Delivery service",
	)
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

func (self *DeliveryServiceTestSuite) TestAcceptDeliveryPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService
	var COMMENT string = "Update"

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		skUser          models.User
		refdelivery     requests.Delivery
		deliveryCompany models.DeliveryCompany
		tempPhotoIds    []uuid.UUID
		savedPhotoIds   []uuid.UUID
		anyTemp         []any
		anySaved        []any
		form            delivery.AcceptForm
	)

	describeAcceptDelivery(t,
		"Accept existed delivery",
		"Accept existed delivery by storekeeper",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create photo ids", func(sCtx provider.StepCtx) {
			tempPhotoIds = testcommon.AssignParameter(sCtx, "tempPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)

			savedPhotoIds = testcommon.AssignParameter(sCtx, "savedPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)

			anyTemp = collection.Collect(
				collection.MapIterator(
					func(v *uuid.UUID) any { return *v },
					collection.SliceIterator(tempPhotoIds),
				),
			)

			anySaved = collection.Collect(
				collection.MapIterator(
					func(v *uuid.UUID) any { return *v },
					collection.SliceIterator(savedPhotoIds),
				),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			refdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomSent(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetById(refdelivery.Id).
						Return(refdelivery, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, toPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().AcceptDelivery(
						instance.Id,
						refdelivery.Id,
						&COMMENT,
						refdelivery.VerificationCode,
					).Return(nil).Times(1)
				}).
				WithPhotoRegistry(func(registry *photoregistry.MockIRegistry) {
					i := 0
					registry.EXPECT().MoveFromTemp(gomock.AnyOf(anyTemp...)).
						DoAndReturn(func(_ uuid.UUID) (uuid.UUID, error) {
							j := i % len(savedPhotoIds)
							i++
							return savedPhotoIds[j], nil
						}).
						Times(len(tempPhotoIds))
				}).
				WithInstancePhotoRepository(func(photo *rinstance.MockIPhotoRepository) {
					photo.EXPECT().Create(instance.Id, gomock.AnyOf(anySaved...)).
						Return(nil).
						Times(len(savedPhotoIds))
				}).
				GetService()
		})
	})

	form = delivery.AcceptForm{
		DeliveryId:       refdelivery.Id,
		Comment:          &COMMENT,
		VerificationCode: refdelivery.VerificationCode,
		TempPhotos:       tempPhotoIds,
	}

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			err = service.AcceptDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *DeliveryServiceTestSuite) TestAcceptDeliveryConflict(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService
	var COMMENT string = "Update"

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		skUser          models.User
		refdelivery     requests.Delivery
		deliveryCompany models.DeliveryCompany
		tempPhotoIds    []uuid.UUID
		form            delivery.AcceptForm
	)

	describeAcceptDelivery(t,
		"Attemp to accept delivery in wrong state",
		"Check that Conflict error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create photo ids", func(sCtx provider.StepCtx) {
			tempPhotoIds = testcommon.AssignParameter(sCtx, "tempPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			refdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomAccepted(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetById(refdelivery.Id).
						Return(refdelivery, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, toPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().AcceptDelivery(
						instance.Id,
						refdelivery.Id,
						&COMMENT,
						refdelivery.VerificationCode,
					).Return(
						istates.ErrorForbiddenMethod{
							Err: errors.New("Instance isn't being delivered"),
						},
					).Times(1)
				}).
				GetService()
		})
	})

	form = delivery.AcceptForm{
		DeliveryId:       refdelivery.Id,
		Comment:          &COMMENT,
		VerificationCode: refdelivery.VerificationCode,
		TempPhotos:       tempPhotoIds,
	}

	// Act
	var err error

	t.WithNewStep("Try to accept already accepted delivery",
		func(sCtx provider.StepCtx) {
			err = service.AcceptDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var cerr cmnerrors.ErrorConflict

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &cerr, "Error Conflict must be returned")
}

func (self *DeliveryServiceTestSuite) TestCreateDeliveryPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		skUser          models.User
		mdelivery       requests.Delivery
		refdelivery     delivery.Delivery
		deliveryCompany models.DeliveryCompany
		form            delivery.CreateForm
	)

	describeCreateDelivery(t,
		"Create delivery",
		"Must return created delivery without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			mdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomCreated(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
			refdelivery = MapDelivery(&mdelivery)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, fromPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceAccess(func(acc *access.MockIInstance) {
					acc.EXPECT().Access(skUser.Id, instance.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().CreateDelivery(
						instance.Id,
						fromPuP.Id,
						toPuP.Id,
					).Return(mdelivery, nil).Times(1)
				}).
				GetService()
		})
	})

	form = delivery.CreateForm{
		InstanceId: instance.Id,
		From:       fromPuP.Id,
		To:         toPuP.Id,
	}

	// Act
	var result delivery.Delivery
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = service.CreateDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(refdelivery, result, "Result matches expected reference")
}

func (self *DeliveryServiceTestSuite) TestCreateDeliveryConflict(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instance models.Instance
		fromPuP  models.PickUpPoint
		toPuP    models.PickUpPoint
		skUser   models.User
		form     delivery.CreateForm
	)

	describeCreateDelivery(t,
		"Instance is not available for rent",
		"Must return ErrorConflic if instance isn't stored in pick up point",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, fromPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceAccess(func(acc *access.MockIInstance) {
					acc.EXPECT().Access(skUser.Id, instance.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().CreateDelivery(
						instance.Id,
						fromPuP.Id,
						toPuP.Id,
					).Return(requests.Delivery{}, istates.ErrorForbiddenMethod{
						Err: errors.New("Instance is rented"),
					}).Times(1)
				}).
				GetService()
		})
	})

	form = delivery.CreateForm{
		InstanceId: instance.Id,
		From:       fromPuP.Id,
		To:         toPuP.Id,
	}

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = service.CreateDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var cerr cmnerrors.ErrorConflict

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &cerr, "Error is ErrorConflict")
}

func (self *DeliveryServiceTestSuite) TestSendDeliveryPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		skUser          models.User
		refdelivery     requests.Delivery
		deliveryCompany models.DeliveryCompany
		tempPhotoIds    []uuid.UUID
		savedPhotoIds   []uuid.UUID
		anyTemp         []any
		anySaved        []any
		form            delivery.SendForm
	)

	describeSendDelivery(t,
		"Send created delivery",
		"Send created delivery by storekeeper",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create photo ids", func(sCtx provider.StepCtx) {
			tempPhotoIds = testcommon.AssignParameter(sCtx, "tempPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)

			savedPhotoIds = testcommon.AssignParameter(sCtx, "savedPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)

			anyTemp = collection.Collect(
				collection.MapIterator(
					func(v *uuid.UUID) any { return *v },
					collection.SliceIterator(tempPhotoIds),
				),
			)

			anySaved = collection.Collect(
				collection.MapIterator(
					func(v *uuid.UUID) any { return *v },
					collection.SliceIterator(savedPhotoIds),
				),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			refdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomCreated(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetById(refdelivery.Id).
						Return(refdelivery, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, fromPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().SendDelivery(
						instance.Id,
						refdelivery.Id,
						refdelivery.VerificationCode,
					).Return(nil).Times(1)
				}).
				WithPhotoRegistry(func(registry *photoregistry.MockIRegistry) {
					i := 0
					registry.EXPECT().MoveFromTemp(gomock.AnyOf(anyTemp...)).
						DoAndReturn(func(_ uuid.UUID) (uuid.UUID, error) {
							j := i % len(savedPhotoIds)
							i++
							return savedPhotoIds[j], nil
						}).
						Times(len(tempPhotoIds))
				}).
				WithInstancePhotoRepository(func(photo *rinstance.MockIPhotoRepository) {
					photo.EXPECT().Create(instance.Id, gomock.AnyOf(anySaved...)).
						Return(nil).
						Times(len(savedPhotoIds))
				}).
				GetService()
		})
	})

	form = delivery.SendForm{
		DeliveryId:       refdelivery.Id,
		VerificationCode: refdelivery.VerificationCode,
		TempPhotos:       tempPhotoIds,
	}

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			err = service.SendDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
}

func (self *DeliveryServiceTestSuite) TestSednDeliveryConflict(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		skUser          models.User
		refdelivery     requests.Delivery
		deliveryCompany models.DeliveryCompany
		tempPhotoIds    []uuid.UUID
		form            delivery.SendForm
	)

	describeSendDelivery(t,
		"Attemp to send delivery that already sent",
		"Check that Conflict error is returned",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create photo ids", func(sCtx provider.StepCtx) {
			tempPhotoIds = testcommon.AssignParameter(sCtx, "tempPhotoIds",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) uuid.UUID { return uuidgen.Generate() },
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			refdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomSent(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetById(refdelivery.Id).
						Return(refdelivery, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, fromPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				WithInstanceStateMachine(func(machine *states.MockIInstanceStateMachine) {
					machine.EXPECT().SendDelivery(
						instance.Id,
						refdelivery.Id,
						refdelivery.VerificationCode,
					).Return(
						istates.ErrorForbiddenMethod{
							Err: errors.New("Instance is already sent"),
						},
					).Times(1)
				}).
				GetService()
		})
	})

	form = delivery.SendForm{
		DeliveryId:       refdelivery.Id,
		VerificationCode: refdelivery.VerificationCode,
		TempPhotos:       tempPhotoIds,
	}

	// Act
	var err error

	t.WithNewStep("Try to accept already accepted delivery",
		func(sCtx provider.StepCtx) {
			err = service.SendDelivery(token.Token(skUser.Token), form)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("form", form),
	)

	// Assert
	var cerr cmnerrors.ErrorConflict

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &cerr, "Error Conflict must be returned")
}

func (self *DeliveryServiceTestSuite) TestListDeliveriesByPickUpPointPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instances         []models.Instance
		fromPuP           models.PickUpPoint
		toPuP             models.PickUpPoint
		skUser            models.User
		deliveries        []requests.Delivery
		refdeliveries     []delivery.Delivery
		deliveryCompanies []models.DeliveryCompany
	)

	describeListDeliveriesByPickUpPoint(t,
		"List deliveries by pick up point (storekeeper)",
		"All values must be returned wihtout error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instances", func(sCtx provider.StepCtx) {
			instances = testcommon.AssignParameter(sCtx, "instances",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) models.Instance {
							return models_om.InstanceRandomId().Build()
						},
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create delivery companis", func(sCtx provider.StepCtx) {
			deliveryCompanies = testcommon.AssignParameter(sCtx, "deliveryCompanies",
				collection.Collect(
					collection.MapIterator(
						func(_ *int) models.DeliveryCompany {
							return models_om.DeliveryCompanyRandomId().Build()
						},
						collection.RangeIterator(collection.RangeEnd(5)),
					),
				),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			deliveries = testcommon.AssignParameter(sCtx, "deliveries",
				collection.Collect(
					collection.MapIterator(
						func(values *collection.Pair[models.DeliveryCompany, models.Instance]) requests.Delivery {
							return requests_om.DeliveryRandomSent(
								values.A.Id,
								values.B.Id,
								fromPuP.Id,
								toPuP.Id,
								nullable.None[string](),
								nullable.None[time.Time](),
								nullable.None[time.Time](),
								nullable.None[time.Time](),
								nullable.None[string](),
								nullable.None[time.Time](),
							).Build()
						},
						collection.ZipIterator(
							collection.SliceIterator(deliveryCompanies),
							collection.SliceIterator(instances),
						),
					),
				),
			)

			refdeliveries = collection.Collect(
				collection.MapIterator(
					MapDelivery,
					collection.SliceIterator(deliveries),
				),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetActiveByPickUpPointId(toPuP.Id).
						Return(collection.SliceCollection(deliveries), nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, toPuP.Id).
						Return(nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error
	var result collection.Collection[delivery.Delivery]

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = service.ListDeliveriesByPickUpPoint(token.Token(skUser.Token), toPuP.Id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickIpPointId", toPuP.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(refdeliveries, collection.Collect(result.Iter()),
		"Elements must match",
	)
}

func (self *DeliveryServiceTestSuite) TestListDeliveriesByPickUpPointNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

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
		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create id to fetch", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithPickUpPointAccess(func(acc *access.MockIPickUpPoint) {
					acc.EXPECT().Access(skUser.Id, id).
						Return(cmnerrors.NotFound("pick_up_point_id")).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = service.ListDeliveriesByPickUpPoint(token.Token(skUser.Token), id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("pickIpPointId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is NotFound")
}

func (self *DeliveryServiceTestSuite) TestListDeliveriesByInstancePositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

	var (
		instance        models.Instance
		fromPuP         models.PickUpPoint
		toPuP           models.PickUpPoint
		user            models.User
		mdelivery       requests.Delivery
		refdelivery     delivery.Delivery
		deliveryCompany models.DeliveryCompany
		rentRequest     requests.Rent
	)

	describeGetDeliveryByInstance(t,
		"List deliveries by instance (user)",
		"All values must be returned wihtout error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create pick up points", func(sCtx provider.StepCtx) {
			fromPuP = testcommon.AssignParameter(sCtx, "fromPuP",
				models_om.PickUpPointExample("From").Build(),
			)
			toPuP = testcommon.AssignParameter(sCtx, "toPuP",
				models_om.PickUpPointExample("To").Build(),
			)
		})

		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
			instance = testcommon.AssignParameter(sCtx, "instance",
				models_om.InstanceRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyRandomId().Build(),
			)
		})

		t.WithNewStep("Create delivery", func(sCtx provider.StepCtx) {
			mdelivery = testcommon.AssignParameter(sCtx, "delivery",
				requests_om.DeliveryRandomSent(
					deliveryCompany.Id,
					instance.Id,
					fromPuP.Id,
					toPuP.Id,
					nullable.None[string](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[time.Time](),
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)

			refdelivery = MapDelivery(&mdelivery)
		})

		t.WithNewStep("Create rent request", func(sCtx provider.StepCtx) {
			var id uuid.UUID

			sCtx.WithNewStep("Create period id", func(sCtx provider.StepCtx) {
				id = uuidgen.Generate()
			})

			rentRequest = testcommon.AssignParameter(sCtx, "rentRequest",
				requests_om.Rent(
					instance.Id,
					user.Id,
					toPuP.Id,
					id,
					nullable.None[string](),
					nullable.None[time.Time](),
				).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(user.Token)).
						Return(user, nil).
						MinTimes(1)
				}).
				WithInstanceAccess(func(acc *access.MockIInstance) {
					acc.EXPECT().Access(user.Id, instance.Id).
						Return(cmnerrors.NoAccess("No access to isntance")).
						MinTimes(1)
				}).
				WithRentRequestRepository(func(rentReq *rrent.MockIRequestRepository) {
					rentReq.EXPECT().GetByInstanceId(instance.Id).
						Return(rentRequest, nil).
						MinTimes(1)
				}).
				WithDeliveryRepository(func(delivery *rdelivery.MockIRepository) {
					delivery.EXPECT().GetActiveByInstanceId(instance.Id).
						Return(mdelivery, nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error
	var result delivery.Delivery

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = service.GetDeliveryByInstance(token.Token(user.Token), instance.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("instanceId", instance.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(refdelivery, result, "Same value returned")
}

func (self *DeliveryServiceTestSuite) TestListDeliveriesByInstanceNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.IService

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
		t.WithNewStep("Create storekeeper", func(sCtx provider.StepCtx) {
			skUser = testcommon.AssignParameter(sCtx, "user",
				models_om.UserStorekeeper(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create id to fetch", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(skUser.Token)).
						Return(skUser, nil).
						MinTimes(1)
				}).
				WithInstanceAccess(func(acc *access.MockIInstance) {
					acc.EXPECT().Access(skUser.Id, id).
						Return(cmnerrors.NotFound("instance_id")).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = service.GetDeliveryByInstance(token.Token(skUser.Token), id)
		},
		allure.NewParameter("token", skUser.Token),
		allure.NewParameter("instanceId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Assert().ErrorAs(err, &nferr, "Error is NotFound")
}

type DeliveryCompanyServiceTestSuite struct {
	suite.Suite
}

type CompanyServiceBuilder struct {
	authenticator *authenticator.MockIAuthenticator
	company       *rdelivery.MockICompanyRepository
}

func NewCompanyServiceBuilder(ctrl *gomock.Controller) *CompanyServiceBuilder {
	return &CompanyServiceBuilder{
		authenticator.NewMockIAuthenticator(ctrl),
		rdelivery.NewMockICompanyRepository(ctrl),
	}
}

func (self *CompanyServiceBuilder) WithAuthenticator(f func(authenticator *authenticator.MockIAuthenticator)) *CompanyServiceBuilder {
	f(self.authenticator)
	return self
}

func (self *CompanyServiceBuilder) WithDeliveryComanyRepository(f func(repo *rdelivery.MockICompanyRepository)) *CompanyServiceBuilder {
	f(self.company)
	return self
}

func (self *CompanyServiceBuilder) GetService() delivery.ICompanyService {
	return service.NewCompany(
		self.authenticator,
		delivery_pmock.NewCompany(self.company),
	)
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

func (self *DeliveryCompanyServiceTestSuite) BeforeEach(t provider.T) {
	testcommon.SetBase(t,
		"DefServices",
		"Default services implementation",
		"Delivery Company service",
	)
}

var describeGetDeliveryCompanyById = testcommon.MethodDescriptor(
	"GetDeliveryCompanyById",
	"Get delivery company by id",
)

var describeListDeliveryCompanies = testcommon.MethodDescriptor(
	"ListDeliveryCompanies",
	"List all delivery companies",
)

func (self *DeliveryCompanyServiceTestSuite) TestListDeliveryCompaniesPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.ICompanyService

	var (
		user              models.User
		deliveryCompanies []models.DeliveryCompany
		reference         []delivery.DeliveryCompany
	)

	describeListDeliveryCompanies(t,
		"List all delivery companies",
		"All compnanies must be returned withou error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create delivery companies", func(sCtx provider.StepCtx) {
			deliveryCompanies = testcommon.AssignParameter(sCtx, "deliveryCompanies",
				collection.Collect(collection.MapIterator(
					func(i *int) models.DeliveryCompany {
						return models_om.DeliveryCompanyExample(fmt.Sprint(*i)).
							Build()
					},
					collection.RangeIterator(collection.RangeEnd(5)),
				)),
			)

			reference = collection.Collect(collection.MapIterator(
				MapDeliveryCompany, collection.SliceIterator(deliveryCompanies),
			))
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewCompanyServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(user.Token)).
						Return(user, nil).
						MinTimes(1)
				}).
				WithDeliveryComanyRepository(func(repo *rdelivery.MockICompanyRepository) {
					repo.EXPECT().GetAll().
						Return(collection.SliceCollection(deliveryCompanies), nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var result collection.Collection[delivery.DeliveryCompany]
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = service.ListDeliveryCompanies(token.Token(user.Token))
		},
		allure.NewParameter("token", user.Token),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().ElementsMatch(reference, collection.Collect(result.Iter()),
		"All companies returned",
	)
}

func (self *DeliveryCompanyServiceTestSuite) TestListDeliveryCompaniesInternalError(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.ICompanyService
	var user models.User

	describeListDeliveryCompanies(t,
		"Internal error during list all deliveries",
		"Error must be mapped",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewCompanyServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(user.Token)).
						Return(user, nil).
						MinTimes(1)
				}).
				WithDeliveryComanyRepository(func(repo *rdelivery.MockICompanyRepository) {
					repo.EXPECT().GetAll().
						Return(nil, errors.New("Some internal error")).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = service.ListDeliveryCompanies(token.Token(user.Token))
		},
		allure.NewParameter("token", user.Token),
	)

	// Assert
	var ierr cmnerrors.ErrorInternal
	var derr cmnerrors.ErrorDataAccess

	t.Require().Error(err, "Error must be returned")
	t.Assert().ErrorAs(err, &ierr, "Error is internal")
	t.Assert().ErrorAs(ierr, &derr, "Error is data access")
}

func (self *DeliveryCompanyServiceTestSuite) TestGetDeliveryCompanyByIdPositive(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.ICompanyService

	var (
		user            models.User
		deliveryCompany models.DeliveryCompany
		reference       delivery.DeliveryCompany
	)

	describeGetDeliveryCompanyById(t,
		"Get delivery company by id",
		"Value must be returned without error",
	)

	// Arrange
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create delivery company", func(sCtx provider.StepCtx) {
			deliveryCompany = testcommon.AssignParameter(sCtx, "deliveryCompany",
				models_om.DeliveryCompanyExample("").Build(),
			)

			reference = MapDeliveryCompany(&deliveryCompany)
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewCompanyServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(user.Token)).
						Return(user, nil).
						MinTimes(1)
				}).
				WithDeliveryComanyRepository(func(repo *rdelivery.MockICompanyRepository) {
					repo.EXPECT().GetById(deliveryCompany.Id).
						Return(deliveryCompany, nil).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var result delivery.DeliveryCompany
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			result, err = service.GetDeliveryCompanyById(token.Token(user.Token), deliveryCompany.Id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("companyId", deliveryCompany.Id),
	)

	// Assert
	t.Require().Nil(err, "No error must be returned")
	t.Require().Equal(reference, result, "Company returned")
}

func (self *DeliveryCompanyServiceTestSuite) TestGetDeliveryCompanyByIdNotFound(t provider.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var service delivery.ICompanyService

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
		t.WithNewStep("Create user", func(sCtx provider.StepCtx) {
			user = testcommon.AssignParameter(sCtx, "user",
				models_om.UserDefault(nullable.None[string]()).Build(),
			)
		})

		t.WithNewStep("Create random id", func(sCtx provider.StepCtx) {
			id = testcommon.AssignParameter(sCtx, "id", uuidgen.Generate())
		})

		t.WithNewStep("Create service", func(sCtx provider.StepCtx) {
			service = NewCompanyServiceBuilder(ctrl).
				WithAuthenticator(func(auth *authenticator.MockIAuthenticator) {
					auth.EXPECT().LoginWithToken(token.Token(user.Token)).
						Return(user, nil).
						MinTimes(1)
				}).
				WithDeliveryComanyRepository(func(repo *rdelivery.MockICompanyRepository) {
					repo.EXPECT().GetById(id).
						Return(
							models.DeliveryCompany{},
							repo_errors.NotFound("delivery_company_id"),
						).
						MinTimes(1)
				}).
				GetService()
		})
	})

	// Act
	var err error

	t.WithNewStep("Accept delivery",
		func(sCtx provider.StepCtx) {
			_, err = service.GetDeliveryCompanyById(token.Token(user.Token), id)
		},
		allure.NewParameter("token", user.Token),
		allure.NewParameter("companyId", id),
	)

	// Assert
	var nferr cmnerrors.ErrorNotFound

	t.Require().NotNil(err, "Error must be returned")
	t.Require().ErrorAs(err, &nferr, "Error is NotFound")
}

func TestDeliveryServiceTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryServiceTestSuite))
}

func TestDeliveryCompanyServiceTestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(DeliveryCompanyServiceTestSuite))
}

