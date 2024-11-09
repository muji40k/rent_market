package deffactory

import (
	cv1 "rent_service/internal/logic/context/v1"
	delivery_creator "rent_service/internal/logic/delivery"
	"rent_service/internal/logic/services/implementations/defservices/access"
	"rent_service/internal/logic/services/implementations/defservices/access/implementations/defaccess"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/authenticator/implementations/defauth"
	"rent_service/internal/logic/services/implementations/defservices/authorizer"
	"rent_service/internal/logic/services/implementations/defservices/authorizer/implementations/defauthorizer"
	"rent_service/internal/logic/services/implementations/defservices/codegen"
	"rent_service/internal/logic/services/implementations/defservices/paymentcheckers"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry"
	"rent_service/internal/logic/services/implementations/defservices/services/category"
	"rent_service/internal/logic/services/implementations/defservices/services/delivery"
	"rent_service/internal/logic/services/implementations/defservices/services/instance"
	"rent_service/internal/logic/services/implementations/defservices/services/login"
	"rent_service/internal/logic/services/implementations/defservices/services/payment"
	"rent_service/internal/logic/services/implementations/defservices/services/period"
	"rent_service/internal/logic/services/implementations/defservices/services/photo"
	"rent_service/internal/logic/services/implementations/defservices/services/pickuppoint"
	"rent_service/internal/logic/services/implementations/defservices/services/product"
	"rent_service/internal/logic/services/implementations/defservices/services/provide"
	"rent_service/internal/logic/services/implementations/defservices/services/rent"
	"rent_service/internal/logic/services/implementations/defservices/services/storage"
	"rent_service/internal/logic/services/implementations/defservices/services/user"
	"rent_service/internal/logic/services/implementations/defservices/states"
	"rent_service/internal/logic/services/implementations/defservices/states/implementations/defstates"
	service_category "rent_service/internal/logic/services/interfaces/category"
	service_delivery "rent_service/internal/logic/services/interfaces/delivery"
	service_instance "rent_service/internal/logic/services/interfaces/instance"
	service_login "rent_service/internal/logic/services/interfaces/login"
	service_payment "rent_service/internal/logic/services/interfaces/payment"
	service_period "rent_service/internal/logic/services/interfaces/period"
	service_photo "rent_service/internal/logic/services/interfaces/photo"
	service_pickuppoint "rent_service/internal/logic/services/interfaces/pickuppoint"
	service_product "rent_service/internal/logic/services/interfaces/product"
	service_provide "rent_service/internal/logic/services/interfaces/provide"
	service_rent "rent_service/internal/logic/services/interfaces/rent"
	service_storage "rent_service/internal/logic/services/interfaces/storage"
	service_user "rent_service/internal/logic/services/interfaces/user"
	rv1 "rent_service/internal/repository/context/v1"

	"github.com/google/uuid"
)

type accessors struct {
	isntance         access.IInstance
	pickUpPoint      access.IPickUpPoint
	provision        access.IProvision
	provisionRequest access.IProvisionRequest
	provisionRevoke  access.IProvisionRevoke
	rent             access.IRent
	rentRequest      access.IRentRequest
	rentReturn       access.IRentReturn
	renter           access.IRenter
	user             access.IUser
}

type static struct {
	accessors     accessors
	authenticator authenticator.IAuthenticator
	authorizer    authorizer.IAuthorizer
	stateMachine  states.IInstanceStateMachine
	registry      photoregistry.IRegistry
}

type Factory struct {
	repositories      *rv1.Context
	generator         codegen.IGenerator
	registryStorage   defregistry.IStorage
	deliveryCreator   delivery_creator.ICreator
	payMethodCheckers map[uuid.UUID]paymentcheckers.IRegistrationChecker
	static            static
}

func New(
	repositories *rv1.Context,
	generator codegen.IGenerator,
	registryStorage defregistry.IStorage,
	deliveryCreator delivery_creator.ICreator,
	payMethodCheckers map[uuid.UUID]paymentcheckers.IRegistrationChecker,
) *Factory {
	return &Factory{
		repositories,
		generator,
		registryStorage,
		deliveryCreator,
		payMethodCheckers,
		static{},
	}
}

func (self *Factory) ToFactories() cv1.Factories {
	return cv1.Factories{
		Category:                self,
		Delivery:                self,
		DeliveryCompany:         self,
		Instance:                self,
		InstancePayPlans:        self,
		InstancePhoto:           self,
		InstanceReview:          self,
		Login:                   self,
		PayMethod:               self,
		UserPayMethod:           self,
		RentPayment:             self,
		Period:                  self,
		Photo:                   self,
		PickUpPoint:             self,
		PickUpPointPhoto:        self,
		PickUpPointWorkingHours: self,
		Product:                 self,
		ProductCharacteristics:  self,
		ProductPhoto:            self,
		Provision:               self,
		ProvisionRequest:        self,
		ProvisionRevoke:         self,
		Rent:                    self,
		RentRequest:             self,
		RentReturn:              self,
		Storage:                 self,
		User:                    self,
		UserProfile:             self,
		UserFavorite:            self,
		Role:                    self,
	}
}

func (self *Factory) Clear() {}

func (self *Factory) CreateAuthenticator() authenticator.IAuthenticator {
	if nil == self.static.authenticator {
		self.static.authenticator = defauth.New(self.repositories)
	}

	return self.static.authenticator
}

func (self *Factory) CreateAuthorizer() authorizer.IAuthorizer {
	if nil == self.static.authorizer {
		self.static.authorizer = defauthorizer.New(
			self.repositories,
			self.repositories,
			self.repositories,
		)
	}

	return self.static.authorizer
}

func (self *Factory) CreateStateMachine() states.IInstanceStateMachine {
	if nil == self.static.stateMachine {
		self.static.stateMachine = defstates.New(
			self.deliveryCreator,
			self.generator,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
			self.repositories,
		)
	}

	return self.static.stateMachine
}

func (self *Factory) CreatePhotoRegistry() photoregistry.IRegistry {
	if nil == self.static.registry {
		self.static.registry = defregistry.New(
			self.repositories,
			self.repositories,
			self.registryStorage,
		)
	}

	return self.static.registry
}

func (self *Factory) CreateInstanceAccessor() access.IInstance {
	if nil == self.static.accessors.isntance {
		self.static.accessors.isntance = defaccess.NewInstance(
			self.CreateAuthorizer(),
			self.repositories,
			self.repositories,
			self.repositories,
		)
	}

	return self.static.accessors.isntance
}

func (self *Factory) CreatePickUpPointAccessor() access.IPickUpPoint {
	if nil == self.static.accessors.pickUpPoint {
		self.static.accessors.pickUpPoint = defaccess.NewPickUpPoint(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.pickUpPoint
}

func (self *Factory) CreateProvisionAccessor() access.IProvision {
	if nil == self.static.accessors.provision {
		self.static.accessors.provision = defaccess.NewProvision(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.provision
}

func (self *Factory) CreateProvisionRequestAccessor() access.IProvisionRequest {
	if nil == self.static.accessors.provisionRequest {
		self.static.accessors.provisionRequest = defaccess.NewProvisionRequest(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.provisionRequest
}

func (self *Factory) CreateProvisionRevokeAccessor() access.IProvisionRevoke {
	if nil == self.static.accessors.provisionRevoke {
		self.static.accessors.provisionRevoke = defaccess.NewProvisionRevoke(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.provisionRevoke
}

func (self *Factory) CreateRentAccessor() access.IRent {
	if nil == self.static.accessors.rent {
		self.static.accessors.rent = defaccess.NewRent(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.rent
}

func (self *Factory) CreateRentRequestAccessor() access.IRentRequest {
	if nil == self.static.accessors.rentRequest {
		self.static.accessors.rentRequest = defaccess.NewRentRequest(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.rentRequest
}

func (self *Factory) CreateRentReturnAccessor() access.IRentReturn {
	if nil == self.static.accessors.rentReturn {
		self.static.accessors.rentReturn = defaccess.NewRentReturn(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.rentReturn
}

func (self *Factory) CreateRenterAccessor() access.IRenter {
	if nil == self.static.accessors.renter {
		self.static.accessors.renter = defaccess.NewRenter(
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.renter
}

func (self *Factory) CreateUserAccessor() access.IUser {
	if nil == self.static.accessors.user {
		self.static.accessors.user = defaccess.NewUser(
			self.repositories,
			self.CreateAuthorizer(),
		)
	}

	return self.static.accessors.user
}

// Factory implementation
func (self *Factory) CreateCategoryService() service_category.IService {
	return category.New(self.repositories)
}

func (self *Factory) CreateDeliveryService() service_delivery.IService {
	return delivery.New(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
		self.repositories,
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreatePickUpPointAccessor(),
	)
}

func (self *Factory) CreateDeliveryCompanyService() service_delivery.ICompanyService {
	return delivery.NewCompany(self.CreateAuthenticator(), self.repositories)
}

func (self *Factory) CreateInstanceService() service_instance.IService {
	return instance.New(
		self.CreateAuthenticator(),
		self.repositories,
		self.CreateInstanceAccessor(),
	)
}

func (self *Factory) CreateInstancePayPlansService() service_instance.IPayPlansService {
	return instance.NewPayPlans(
		self.CreateAuthenticator(),
		self.repositories,
		self.CreateInstanceAccessor(),
	)
}

func (self *Factory) CreateInstancePhotoService() service_instance.IPhotoService {
	return instance.NewPhoto(
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
		self.repositories,
		self.CreateInstanceAccessor(),
	)
}

func (self *Factory) CreateInstanceReviewService() service_instance.IReviewService {
	return instance.NewReview(
		self.CreateAuthenticator(),
		self.repositories,
		self.repositories,
	)
}

func (self *Factory) CreateLoginService() service_login.IService {
	return login.New(self.repositories)
}

func (self *Factory) CreatePayMethodService() service_payment.IPayMethodService {
	return payment.NewPayMethod(self.repositories)
}

func (self *Factory) CreateUserPayMethodService() service_payment.IUserPayMethodService {
	return payment.NewUserPayMethod(
		self.CreateAuthenticator(),
		self.CreateAuthorizer(),
		self.payMethodCheckers,
		self.repositories,
	)
}

func (self *Factory) CreateRentPaymentService() service_payment.IRentPaymentService {
	return payment.NewRentPayment(
		self.CreateAuthenticator(),
		self.CreateAuthorizer(),
		self.CreateInstanceAccessor(),
		self.CreateRentAccessor(),
		self.repositories,
	)
}

func (self *Factory) CreatePeriodService() service_period.IService {
	return period.New(self.repositories)
}

func (self *Factory) CreatePhotoService() service_photo.IService {
	return photo.New(
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
		self.repositories,
		self.repositories,
	)
}

func (self *Factory) CreatePickUpPointService() service_pickuppoint.IService {
	return pickuppoint.New(self.repositories)
}

func (self *Factory) CreatePickUpPointPhotoService() service_pickuppoint.IPhotoService {
	return pickuppoint.NewPhoto(self.repositories)
}

func (self *Factory) CreatePickUpPointWorkingHoursService() service_pickuppoint.IWorkingHoursService {
	return pickuppoint.NewWorkingHours(self.repositories)
}

func (self *Factory) CreateProductService() service_product.IService {
	return product.New(self.repositories)
}

func (self *Factory) CreateProductCharacteristicsService() service_product.ICharacteristicsService {
	return product.NewCharacteristics(self.repositories)
}

func (self *Factory) CreateProductPhotoService() service_product.IPhotoService {
	return product.NewPhoto(self.repositories)
}

func (self *Factory) CreateProvisionService() service_provide.IService {
	return provide.New(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
		self.repositories,
		self.repositories,
		self.repositories,
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreateProvisionRequestAccessor(),
		self.CreateProvisionRevokeAccessor(),
	)
}

func (self *Factory) CreateProvisionRequestService() service_provide.IRequestService {
	return provide.NewRequest(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.CreateAuthorizer(),
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreatePickUpPointAccessor(),
	)
}

func (self *Factory) CreateProvisionRevokeService() service_provide.IRevokeService {
	return provide.NewRevoke(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.CreateAuthorizer(),
		self.repositories,
		self.repositories,
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreatePickUpPointAccessor(),
		self.CreateProvisionAccessor(),
		self.CreateProvisionRevokeAccessor(),
	)
}

func (self *Factory) CreateRentService() service_rent.IService {
	return rent.New(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
		self.repositories,
		self.repositories,
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreateRentRequestAccessor(),
		self.CreateRentReturnAccessor(),
	)
}

func (self *Factory) CreateRentRequestService() service_rent.IRequestService {
	return rent.NewRequest(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreatePickUpPointAccessor(),
	)
}

func (self *Factory) CreateRentReturnService() service_rent.IReturnService {
	return rent.NewReturn(
		self.CreateStateMachine(),
		self.CreateAuthenticator(),
		self.repositories,
		self.repositories,
		self.CreateInstanceAccessor(),
		self.CreateUserAccessor(),
		self.CreatePickUpPointAccessor(),
		self.CreateRentAccessor(),
		self.CreateRentReturnAccessor(),
	)
}

func (self *Factory) CreateStorageService() service_storage.IService {
	return storage.New(
		self.CreateAuthenticator(),
		self.repositories,
		self.CreatePickUpPointAccessor(),
	)
}

func (self *Factory) CreateUserService() service_user.IService {
	return user.New(self.repositories, self.CreateAuthenticator())
}

func (self *Factory) CreateUserProfileService() service_user.IProfileService {
	return user.NewProfile(
		self.repositories,
		self.CreateAuthenticator(),
		self.CreatePhotoRegistry(),
	)
}

func (self *Factory) CreateUserFavoriteService() service_user.IFavoriteService {
	return user.NewFavorite(self.repositories, self.CreateAuthenticator())
}

func (self *Factory) CreateRoleService() service_user.IRoleService {
	return user.NewRole(
		self.CreateAuthenticator(),
		self.CreateAuthorizer(),
		self.repositories,
	)
}

