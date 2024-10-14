package deffactory

import (
	fv1 "rent_service/internal/factory/services/v1"
	cv1 "rent_service/internal/logic/context/v1"
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

type factory struct {
}

func New() fv1.IFactory {
	return &factory{}
}

func (self *factory) ToFactories() cv1.Factories {
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

func (self *factory) CreateCategoryService() service_category.IService {
	return nil
}

func (self *factory) CreateDeliveryService() service_delivery.IService {
	return nil
}

func (self *factory) CreateDeliveryCompanyService() service_delivery.ICompanyService {
	return nil
}

func (self *factory) CreateInstanceService() service_instance.IService {
	return nil
}

func (self *factory) CreateInstancePayPlansService() service_instance.IPayPlansService {
	return nil
}

func (self *factory) CreateInstancePhotoService() service_instance.IPhotoService {
	return nil
}

func (self *factory) CreateInstanceReviewService() service_instance.IReviewService {
	return nil
}

func (self *factory) CreateLoginService() service_login.IService {
	return nil
}

func (self *factory) CreatePayMethodService() service_payment.IPayMethodService {
	return nil
}

func (self *factory) CreateUserPayMethodService() service_payment.IUserPayMethodService {
	return nil
}

func (self *factory) CreateRentPaymentService() service_payment.IRentPaymentService {
	return nil
}

func (self *factory) CreatePeriodService() service_period.IService {
	return nil
}

func (self *factory) CreatePhotoService() service_photo.IService {
	return nil
}

func (self *factory) CreatePickUpPointService() service_pickuppoint.IService {
	return nil
}

func (self *factory) CreatePickUpPointPhotoService() service_pickuppoint.IPhotoService {
	return nil
}

func (self *factory) CreatePickUpPointWorkingHoursService() service_pickuppoint.IWorkingHoursService {
	return nil
}

func (self *factory) CreateProductService() service_product.IService {
	return nil
}

func (self *factory) CreateProductCharacteristicsService() service_product.ICharacteristicsService {
	return nil
}

func (self *factory) CreateProductPhotoService() service_product.IPhotoService {
	return nil
}

func (self *factory) CreateProvisionService() service_provide.IService {
	return nil
}

func (self *factory) CreateProvisionRequestService() service_provide.IRequestService {
	return nil
}

func (self *factory) CreateProvisionRevokeService() service_provide.IRevokeService {
	return nil
}

func (self *factory) CreateRentService() service_rent.IService {
	return nil
}

func (self *factory) CreateRentRequestService() service_rent.IRequestService {
	return nil
}

func (self *factory) CreateRentReturnService() service_rent.IReturnService {
	return nil
}

func (self *factory) CreateStorageService() service_storage.IService {
	return nil
}

func (self *factory) CreateUserService() service_user.IService {
	return nil
}

func (self *factory) CreateUserProfileService() service_user.IProfileService {
	return nil
}

func (self *factory) CreateUserFavoriteService() service_user.IFavoriteService {
	return nil
}

func (self *factory) CreateRoleService() service_user.IRoleService {
	return nil
}

