package v1

import (
	factory_category "rent_service/internal/factory/repositories/interfaces/category"
	factory_delivery "rent_service/internal/factory/repositories/interfaces/delivery"
	factory_instance "rent_service/internal/factory/repositories/interfaces/instance"
	factory_payment "rent_service/internal/factory/repositories/interfaces/payment"
	factory_paymethod "rent_service/internal/factory/repositories/interfaces/paymethod"
	factory_period "rent_service/internal/factory/repositories/interfaces/period"
	factory_photo "rent_service/internal/factory/repositories/interfaces/photo"
	factory_pickuppoint "rent_service/internal/factory/repositories/interfaces/pickuppoint"
	factory_product "rent_service/internal/factory/repositories/interfaces/product"
	factory_provision "rent_service/internal/factory/repositories/interfaces/provision"
	factory_rent "rent_service/internal/factory/repositories/interfaces/rent"
	factory_review "rent_service/internal/factory/repositories/interfaces/review"
	factory_role "rent_service/internal/factory/repositories/interfaces/role"
	factory_storage "rent_service/internal/factory/repositories/interfaces/storage"
	factory_user "rent_service/internal/factory/repositories/interfaces/user"
	repository_category "rent_service/internal/repository/interfaces/category"
	repository_delivery "rent_service/internal/repository/interfaces/delivery"
	repository_instance "rent_service/internal/repository/interfaces/instance"
	repository_payment "rent_service/internal/repository/interfaces/payment"
	repository_paymethod "rent_service/internal/repository/interfaces/paymethod"
	repository_period "rent_service/internal/repository/interfaces/period"
	repository_photo "rent_service/internal/repository/interfaces/photo"
	repository_pickuppoint "rent_service/internal/repository/interfaces/pickuppoint"
	repository_product "rent_service/internal/repository/interfaces/product"
	repository_provision "rent_service/internal/repository/interfaces/provision"
	repository_rent "rent_service/internal/repository/interfaces/rent"
	repository_review "rent_service/internal/repository/interfaces/review"
	repository_role "rent_service/internal/repository/interfaces/role"
	repository_storage "rent_service/internal/repository/interfaces/storage"
	repository_user "rent_service/internal/repository/interfaces/user"
)

type Factories struct {
	Category                factory_category.IFactory
	Delivery                factory_delivery.IFactory
	DeliveryCompany         factory_delivery.ICompanyFactory
	Instance                factory_instance.IFactory
	InstancePayPlans        factory_instance.IPayPlansFactory
	InstancePhoto           factory_instance.IPhotoFactory
	Payment                 factory_payment.IFactory
	PayMethod               factory_paymethod.IFactory
	Period                  factory_period.IFactory
	Photo                   factory_photo.IFactory
	PhotoTemp               factory_photo.ITempFactory
	PickUpPoint             factory_pickuppoint.IFactory
	PickUpPointPhoto        factory_pickuppoint.IPhotoFactory
	PickUpPointWorkingHours factory_pickuppoint.IWorkingHoursFactory
	Product                 factory_product.IFactory
	ProductCharacteristics  factory_product.ICharacteristicsFactory
	ProductPhoto            factory_product.IPhotoFactory
	Provision               factory_provision.IFactory
	ProvisionRequest        factory_provision.IRequestFactory
	ProvisionRevoke         factory_provision.IRevokeFactory
	Rent                    factory_rent.IFactory
	RentRequest             factory_rent.IRequestFactory
	RentReturn              factory_rent.IReturnFactory
	Review                  factory_review.IFactory
	RoleAdministrator       factory_role.IAdministratorFactory
	RoleRenter              factory_role.IRenterFactory
	RoleStorekeeper         factory_role.IStorekeeperFactory
	Storage                 factory_storage.IFactory
	User                    factory_user.IFactory
	UserProfile             factory_user.IProfileFactory
	UserFavorite            factory_user.IFavoriteFactory
	UserPayMethods          factory_user.IPayMethodsFactory
}

type static struct {
	Category                repository_category.IRepository
	Delivery                repository_delivery.IRepository
	DeliveryCompany         repository_delivery.ICompanyRepository
	Instance                repository_instance.IRepository
	InstancePayPlans        repository_instance.IPayPlansRepository
	InstancePhoto           repository_instance.IPhotoRepository
	Payment                 repository_payment.IRepository
	PayMethod               repository_paymethod.IRepository
	Period                  repository_period.IRepository
	Photo                   repository_photo.IRepository
	PhotoTemp               repository_photo.ITempRepository
	PickUpPoint             repository_pickuppoint.IRepository
	PickUpPointPhoto        repository_pickuppoint.IPhotoRepository
	PickUpPointWorkingHours repository_pickuppoint.IWorkingHoursRepository
	Product                 repository_product.IRepository
	ProductCharacteristics  repository_product.ICharacteristicsRepository
	ProductPhoto            repository_product.IPhotoRepository
	Provision               repository_provision.IRepository
	ProvisionRequest        repository_provision.IRequestRepository
	ProvisionRevoke         repository_provision.IRevokeRepository
	Rent                    repository_rent.IRepository
	RentRequest             repository_rent.IRequestRepository
	RentReturn              repository_rent.IReturnRepository
	Review                  repository_review.IRepository
	RoleAdministrator       repository_role.IAdministratorRepository
	RoleRenter              repository_role.IRenterRepository
	RoleStorekeeper         repository_role.IStorekeeperRepository
	Storage                 repository_storage.IRepository
	User                    repository_user.IRepository
	UserProfile             repository_user.IProfileRepository
	UserFavorite            repository_user.IFavoriteRepository
	UserPayMethods          repository_user.IPayMethodsRepository
}

type Context struct {
	factories    Factories
	repositories static
}

func New(factories Factories) *Context {
	return &Context{
		factories,
		static{},
	}
}

// Provider implementations
func (self *Context) GetCategoryRepository() repository_category.IRepository {
	if nil == self.repositories.Category {
		self.repositories.Category = self.factories.Category.CreateCategoryRepository()
	}

	return self.repositories.Category
}

func (self *Context) GetDeliveryRepository() repository_delivery.IRepository {
	if nil == self.repositories.Delivery {
		self.repositories.Delivery = self.factories.Delivery.CreateDeliveryRepository()
	}

	return self.repositories.Delivery
}

func (self *Context) GetDeliveryCompanyRepository() repository_delivery.ICompanyRepository {
	if nil == self.repositories.DeliveryCompany {
		self.repositories.DeliveryCompany = self.factories.DeliveryCompany.CreateDeliveryCompanyRepository()
	}

	return self.repositories.DeliveryCompany
}

func (self *Context) GetInstanceRepository() repository_instance.IRepository {
	if nil == self.repositories.Instance {
		self.repositories.Instance = self.factories.Instance.CreateInstanceRepository()
	}

	return self.repositories.Instance
}

func (self *Context) GetInstancePayPlansRepository() repository_instance.IPayPlansRepository {
	if nil == self.repositories.InstancePayPlans {
		self.repositories.InstancePayPlans = self.factories.InstancePayPlans.CreateInstancePayPlansRepository()
	}

	return self.repositories.InstancePayPlans
}

func (self *Context) GetInstancePhotoRepository() repository_instance.IPhotoRepository {
	if nil == self.repositories.InstancePhoto {
		self.repositories.InstancePhoto = self.factories.InstancePhoto.CreateInstancePhotoRepository()
	}

	return self.repositories.InstancePhoto
}

func (self *Context) GetPaymentRepository() repository_payment.IRepository {
	if nil == self.repositories.Payment {
		self.repositories.Payment = self.factories.Payment.CreatePaymentRepository()
	}

	return self.repositories.Payment
}

func (self *Context) GetPayMethodRepository() repository_paymethod.IRepository {
	if nil == self.repositories.PayMethod {
		self.repositories.PayMethod = self.factories.PayMethod.CreatePayMethodRepository()
	}

	return self.repositories.PayMethod
}

func (self *Context) GetPeriodRepository() repository_period.IRepository {
	if nil == self.repositories.Period {
		self.repositories.Period = self.factories.Period.CreatePeriodRepository()
	}

	return self.repositories.Period
}

func (self *Context) GetPhotoRepository() repository_photo.IRepository {
	if nil == self.repositories.Photo {
		self.repositories.Photo = self.factories.Photo.CreatePhotoRepository()
	}

	return self.repositories.Photo
}

func (self *Context) GetTempPhotoRepository() repository_photo.ITempRepository {
	if nil == self.repositories.PhotoTemp {
		self.repositories.PhotoTemp = self.factories.PhotoTemp.CreatePhotoTempRepository()
	}

	return self.repositories.PhotoTemp
}

func (self *Context) GetPickUpPointRepository() repository_pickuppoint.IRepository {
	if nil == self.repositories.PickUpPoint {
		self.repositories.PickUpPoint = self.factories.PickUpPoint.CreatePickUpPointRepository()
	}

	return self.repositories.PickUpPoint
}

func (self *Context) GetPickUpPointPhotoRepository() repository_pickuppoint.IPhotoRepository {
	if nil == self.repositories.PickUpPointPhoto {
		self.repositories.PickUpPointPhoto = self.factories.PickUpPointPhoto.CreatePickUpPointPhotoRepository()
	}

	return self.repositories.PickUpPointPhoto
}

func (self *Context) GetPickUpPointWorkingHoursRepository() repository_pickuppoint.IWorkingHoursRepository {
	if nil == self.repositories.PickUpPointWorkingHours {
		self.repositories.PickUpPointWorkingHours = self.factories.PickUpPointWorkingHours.CreatePickUpPointWorkingHoursRepository()
	}

	return self.repositories.PickUpPointWorkingHours
}

func (self *Context) GetProductRepository() repository_product.IRepository {
	if nil == self.repositories.Product {
		self.repositories.Product = self.factories.Product.CreateProductRepository()
	}

	return self.repositories.Product
}

func (self *Context) GetProductCharacteristicsRepository() repository_product.ICharacteristicsRepository {
	if nil == self.repositories.ProductCharacteristics {
		self.repositories.ProductCharacteristics = self.factories.ProductCharacteristics.CreateProductCharacteristicsRepository()
	}

	return self.repositories.ProductCharacteristics
}

func (self *Context) GetProductPhotoRepository() repository_product.IPhotoRepository {
	if nil == self.repositories.ProductPhoto {
		self.repositories.ProductPhoto = self.factories.ProductPhoto.CreateProductPhotoRepository()
	}

	return self.repositories.ProductPhoto
}

func (self *Context) GetProvisionRepository() repository_provision.IRepository {
	if nil == self.repositories.Provision {
		self.repositories.Provision = self.factories.Provision.CreateProvisionRepository()
	}

	return self.repositories.Provision
}

func (self *Context) GetProvisionRequestRepository() repository_provision.IRequestRepository {
	if nil == self.repositories.ProvisionRequest {
		self.repositories.ProvisionRequest = self.factories.ProvisionRequest.CreateProvisionRequestRepository()
	}

	return self.repositories.ProvisionRequest
}

func (self *Context) GetRevokeProvisionRepository() repository_provision.IRevokeRepository {
	if nil == self.repositories.ProvisionRevoke {
		self.repositories.ProvisionRevoke = self.factories.ProvisionRevoke.CreateProvisionRevokeRepository()
	}

	return self.repositories.ProvisionRevoke
}

func (self *Context) GetRentRepository() repository_rent.IRepository {
	if nil == self.repositories.Rent {
		self.repositories.Rent = self.factories.Rent.CreateRentRepository()
	}

	return self.repositories.Rent
}

func (self *Context) GetRentRequestRepository() repository_rent.IRequestRepository {
	if nil == self.repositories.RentRequest {
		self.repositories.RentRequest = self.factories.RentRequest.CreateRentRequestRepository()
	}

	return self.repositories.RentRequest
}

func (self *Context) GetRentReturnRepository() repository_rent.IReturnRepository {
	if nil == self.repositories.RentReturn {
		self.repositories.RentReturn = self.factories.RentReturn.CreateRentReturnRepository()
	}

	return self.repositories.RentReturn
}

func (self *Context) GetReviewRepository() repository_review.IRepository {
	if nil == self.repositories.Review {
		self.repositories.Review = self.factories.Review.CreateReviewRepository()
	}

	return self.repositories.Review
}

func (self *Context) GetAdministratorRepository() repository_role.IAdministratorRepository {
	if nil == self.repositories.RoleAdministrator {
		self.repositories.RoleAdministrator = self.factories.RoleAdministrator.CreateRoleAdministratorRepository()
	}

	return self.repositories.RoleAdministrator
}

func (self *Context) GetRenterRepository() repository_role.IRenterRepository {
	if nil == self.repositories.RoleRenter {
		self.repositories.RoleRenter = self.factories.RoleRenter.CreateRoleRenterRepository()
	}

	return self.repositories.RoleRenter
}

func (self *Context) GetStorekeeperRepository() repository_role.IStorekeeperRepository {
	if nil == self.repositories.RoleStorekeeper {
		self.repositories.RoleStorekeeper = self.factories.RoleStorekeeper.CreateRoleStorekeeperRepository()
	}

	return self.repositories.RoleStorekeeper
}

func (self *Context) GetStorageRepository() repository_storage.IRepository {
	if nil == self.repositories.Storage {
		self.repositories.Storage = self.factories.Storage.CreateStorageRepository()
	}

	return self.repositories.Storage
}

func (self *Context) GetUserRepository() repository_user.IRepository {
	if nil == self.repositories.User {
		self.repositories.User = self.factories.User.CreateUserRepository()
	}

	return self.repositories.User
}

func (self *Context) GetUserProfileRepository() repository_user.IProfileRepository {
	if nil == self.repositories.UserProfile {
		self.repositories.UserProfile = self.factories.UserProfile.CreateUserProfileRepository()
	}

	return self.repositories.UserProfile
}

func (self *Context) GetUserFavoriteRepository() repository_user.IFavoriteRepository {
	if nil == self.repositories.UserFavorite {
		self.repositories.UserFavorite = self.factories.UserFavorite.CreateUserFavoriteRepository()
	}

	return self.repositories.UserFavorite
}

func (self *Context) GetUserPayMethodsRepository() repository_user.IPayMethodsRepository {
	if nil == self.repositories.UserPayMethods {
		self.repositories.UserPayMethods = self.factories.UserPayMethods.CreateUserPayMethodsRepository()
	}

	return self.repositories.UserPayMethods
}

