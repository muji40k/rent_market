package v1

import (
	factory_category "rent_service/internal/factory/services/interfaces/category"
	factory_delivery "rent_service/internal/factory/services/interfaces/delivery"
	factory_instance "rent_service/internal/factory/services/interfaces/instance"
	factory_login "rent_service/internal/factory/services/interfaces/login"
	factory_payment "rent_service/internal/factory/services/interfaces/payment"
	factory_period "rent_service/internal/factory/services/interfaces/period"
	factory_photo "rent_service/internal/factory/services/interfaces/photo"
	factory_pickuppoint "rent_service/internal/factory/services/interfaces/pickuppoint"
	factory_product "rent_service/internal/factory/services/interfaces/product"
	factory_provide "rent_service/internal/factory/services/interfaces/provide"
	factory_rent "rent_service/internal/factory/services/interfaces/rent"
	factory_storage "rent_service/internal/factory/services/interfaces/storage"
	factory_user "rent_service/internal/factory/services/interfaces/user"
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
)

type Factories struct {
	Category                factory_category.IFactory
	Delivery                factory_delivery.IFactory
	DeliveryCompany         factory_delivery.ICompanyFactory
	Instance                factory_instance.IFactory
	InstancePayPlans        factory_instance.IPayPlansFactory
	InstancePhoto           factory_instance.IPhotoFactory
	InstanceReview          factory_instance.IReviewFactory
	Login                   factory_login.IFactory
	PayMethod               factory_payment.IPayMethodFactory
	UserPayMethod           factory_payment.IUserPayMethodFactory
	RentPayment             factory_payment.IRentPaymentFactory
	Period                  factory_period.IFactory
	Photo                   factory_photo.IFactory
	PickUpPoint             factory_pickuppoint.IFactory
	PickUpPointPhoto        factory_pickuppoint.IPhotoFactory
	PickUpPointWorkingHours factory_pickuppoint.IWorkingHoursFactory
	Product                 factory_product.IFactory
	ProductCharacteristics  factory_product.ICharacteristicsFactory
	ProductPhoto            factory_product.IPhotoFactory
	Provision               factory_provide.IFactory
	ProvisionRequest        factory_provide.IRequestFactory
	ProvisionRevoke         factory_provide.IRevokeFactory
	Rent                    factory_rent.IFactory
	RentRequest             factory_rent.IRequestFactory
	RentReturn              factory_rent.IReturnFactory
	Storage                 factory_storage.IFactory
	User                    factory_user.IFactory
	UserProfile             factory_user.IProfileFactory
	UserFavorite            factory_user.IFavoriteFactory
	Role                    factory_user.IRoleFactory
}

type static struct {
	Category                service_category.IService
	Delivery                service_delivery.IService
	DeliveryCompany         service_delivery.ICompanyService
	Instance                service_instance.IService
	InstancePayPlans        service_instance.IPayPlansService
	InstancePhoto           service_instance.IPhotoService
	InstanceReview          service_instance.IReviewService
	Login                   service_login.IService
	PayMethod               service_payment.IPayMethodService
	UserPayMethod           service_payment.IUserPayMethodService
	RentPayment             service_payment.IRentPaymentService
	Period                  service_period.IService
	Photo                   service_photo.IService
	PickUpPoint             service_pickuppoint.IService
	PickUpPointPhoto        service_pickuppoint.IPhotoService
	PickUpPointWorkingHours service_pickuppoint.IWorkingHoursService
	Product                 service_product.IService
	ProductCharacteristics  service_product.ICharacteristicsService
	ProductPhoto            service_product.IPhotoService
	Provision               service_provide.IService
	ProvisionRequest        service_provide.IRequestService
	ProvisionRevoke         service_provide.IRevokeService
	Rent                    service_rent.IService
	RentRequest             service_rent.IRequestService
	RentReturn              service_rent.IReturnService
	Storage                 service_storage.IService
	User                    service_user.IService
	UserProfile             service_user.IProfileService
	UserFavorite            service_user.IFavoriteService
	Role                    service_user.IRoleService
}

type Context struct {
	factories Factories
	services  static
}

func New(factories Factories) Context {
	var out Context

	out.factories = factories

	return out
}

// Provider implementations
func (self *Context) GetCategoryService() service_category.IService {
	if nil == self.services.Category {
		self.services.Category = self.factories.Category.CreateCategoryService()
	}

	return self.services.Category
}

func (self *Context) GetDeliveryService() service_delivery.IService {
	if nil == self.services.Delivery {
		self.services.Delivery = self.factories.Delivery.CreateDeliveryService()
	}

	return self.services.Delivery
}

func (self *Context) GetDeliveryCompanyService() service_delivery.ICompanyService {
	if nil == self.services.DeliveryCompany {
		self.services.DeliveryCompany = self.factories.DeliveryCompany.CreateDeliveryCompanyService()
	}

	return self.services.DeliveryCompany
}

func (self *Context) GetInstanceService() service_instance.IService {
	if nil == self.services.Instance {
		self.services.Instance = self.factories.Instance.CreateInstanceService()
	}

	return self.services.Instance
}

func (self *Context) GetInstancePayPlansService() service_instance.IPayPlansService {
	if nil == self.services.InstancePayPlans {
		self.services.InstancePayPlans = self.factories.InstancePayPlans.CreateInstancePayPlansService()
	}

	return self.services.InstancePayPlans
}

func (self *Context) GetInstancePhotoService() service_instance.IPhotoService {
	if nil == self.services.InstancePhoto {
		self.services.InstancePhoto = self.factories.InstancePhoto.CreateInstancePhotoService()
	}

	return self.services.InstancePhoto
}

func (self *Context) GetInstanceReviewService() service_instance.IReviewService {
	if nil == self.services.InstanceReview {
		self.services.InstanceReview = self.factories.InstanceReview.CreateInstanceReviewService()
	}

	return self.services.InstanceReview
}

func (self *Context) GetLoginService() service_login.IService {
	if nil == self.services.Login {
		self.services.Login = self.factories.Login.CreateLoginService()
	}

	return self.services.Login
}

func (self *Context) GetPayMethodService() service_payment.IPayMethodService {
	if nil == self.services.PayMethod {
		self.services.PayMethod = self.factories.PayMethod.CreatePayMethodService()
	}

	return self.services.PayMethod
}

func (self *Context) GetUserPayMethodService() service_payment.IUserPayMethodService {
	if nil == self.services.UserPayMethod {
		self.services.UserPayMethod = self.factories.UserPayMethod.CreateUserPayMethodService()
	}

	return self.services.UserPayMethod
}

func (self *Context) GetRentPaymentService() service_payment.IRentPaymentService {
	if nil == self.services.RentPayment {
		self.services.RentPayment = self.factories.RentPayment.CreateRentPaymentService()
	}

	return self.services.RentPayment
}

func (self *Context) GetPeriodService() service_period.IService {
	if nil == self.services.Period {
		self.services.Period = self.factories.Period.CreatePeriodService()
	}

	return self.services.Period
}

func (self *Context) GetPhotoService() service_photo.IService {
	if nil == self.services.Photo {
		self.services.Photo = self.factories.Photo.CreatePhotoService()
	}

	return self.services.Photo
}

func (self *Context) GetPickUpPointService() service_pickuppoint.IService {
	if nil == self.services.PickUpPoint {
		self.services.PickUpPoint = self.factories.PickUpPoint.CreatePickUpPointService()
	}

	return self.services.PickUpPoint
}

func (self *Context) GetPickUpPointPhotoService() service_pickuppoint.IPhotoService {
	if nil == self.services.PickUpPointPhoto {
		self.services.PickUpPointPhoto = self.factories.PickUpPointPhoto.CreatePickUpPointPhotoService()
	}

	return self.services.PickUpPointPhoto
}

func (self *Context) GetPickUpPointWorkingHoursService() service_pickuppoint.IWorkingHoursService {
	if nil == self.services.PickUpPointWorkingHours {
		self.services.PickUpPointWorkingHours = self.factories.PickUpPointWorkingHours.CreatePickUpPointWorkingHoursService()
	}

	return self.services.PickUpPointWorkingHours
}

func (self *Context) GetProductService() service_product.IService {
	if nil == self.services.Product {
		self.services.Product = self.factories.Product.CreateProductService()
	}

	return self.services.Product
}

func (self *Context) GetProductCharacteristicsService() service_product.ICharacteristicsService {
	if nil == self.services.ProductCharacteristics {
		self.services.ProductCharacteristics = self.factories.ProductCharacteristics.CreateProductCharacteristicsService()
	}

	return self.services.ProductCharacteristics
}

func (self *Context) GetProductPhotoService() service_product.IPhotoService {
	if nil == self.services.ProductPhoto {
		self.services.ProductPhoto = self.factories.ProductPhoto.CreateProductPhotoService()
	}

	return self.services.ProductPhoto
}

func (self *Context) GetProvisionService() service_provide.IService {
	if nil == self.services.Provision {
		self.services.Provision = self.factories.Provision.CreateProvisionService()
	}

	return self.services.Provision
}

func (self *Context) GetProvisionRequestService() service_provide.IRequestService {
	if nil == self.services.ProvisionRequest {
		self.services.ProvisionRequest = self.factories.ProvisionRequest.CreateProvisionRequestService()
	}

	return self.services.ProvisionRequest
}

func (self *Context) GetProvisionRevokeService() service_provide.IRevokeService {
	if nil == self.services.ProvisionRevoke {
		self.services.ProvisionRevoke = self.factories.ProvisionRevoke.CreateProvisionRevokeService()
	}

	return self.services.ProvisionRevoke
}

func (self *Context) GetRentService() service_rent.IService {
	if nil == self.services.Rent {
		self.services.Rent = self.factories.Rent.CreateRentService()
	}

	return self.services.Rent
}

func (self *Context) GetRentRequestService() service_rent.IRequestService {
	if nil == self.services.RentRequest {
		self.services.RentRequest = self.factories.RentRequest.CreateRentRequestService()
	}

	return self.services.RentRequest
}

func (self *Context) GetRentReturnService() service_rent.IReturnService {
	if nil == self.services.RentReturn {
		self.services.RentReturn = self.factories.RentReturn.CreateRentReturnService()
	}

	return self.services.RentReturn
}

func (self *Context) GetStorageService() service_storage.IService {
	if nil == self.services.Storage {
		self.services.Storage = self.factories.Storage.CreateStorageService()
	}

	return self.services.Storage
}

func (self *Context) GetUserService() service_user.IService {
	if nil == self.services.User {
		self.services.User = self.factories.User.CreateUserService()
	}

	return self.services.User
}

func (self *Context) GetUserProfileService() service_user.IProfileService {
	if nil == self.services.UserProfile {
		self.services.UserProfile = self.factories.UserProfile.CreateUserProfileService()
	}

	return self.services.UserProfile
}

func (self *Context) GetUserFavoriteService() service_user.IFavoriteService {
	if nil == self.services.UserFavorite {
		self.services.UserFavorite = self.factories.UserFavorite.CreateUserFavoriteService()
	}

	return self.services.UserFavorite
}

func (self *Context) GetRoleService() service_user.IRoleService {
	if nil == self.services.Role {
		self.services.Role = self.factories.Role.CreateRoleService()
	}

	return self.services.Role
}

