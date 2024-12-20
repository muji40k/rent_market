package tracer

import (
	fv1 "rent_service/internal/factory/services/v1"
	cv1 "rent_service/internal/logic/context/v1"
	"rent_service/internal/logic/services/implementations/tracer/category"
	"rent_service/internal/logic/services/implementations/tracer/delivery"
	"rent_service/internal/logic/services/implementations/tracer/instance"
	"rent_service/internal/logic/services/implementations/tracer/login"
	"rent_service/internal/logic/services/implementations/tracer/payment"
	"rent_service/internal/logic/services/implementations/tracer/period"
	"rent_service/internal/logic/services/implementations/tracer/photo"
	"rent_service/internal/logic/services/implementations/tracer/pickuppoint"
	"rent_service/internal/logic/services/implementations/tracer/product"
	"rent_service/internal/logic/services/implementations/tracer/provide"
	"rent_service/internal/logic/services/implementations/tracer/rent"
	"rent_service/internal/logic/services/implementations/tracer/storage"
	"rent_service/internal/logic/services/implementations/tracer/user"
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
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/trace"
)

type Factory struct {
	wrapped   fv1.IFactory
	factories cv1.Factories
	hl        *contextholder.Holder
	tracer    trace.Tracer
}

func New(
	wrapped fv1.IFactory,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) *Factory {
	return &Factory{
		wrapped:   wrapped,
		factories: wrapped.ToFactories(),
		hl:        hl,
		tracer:    tracer,
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

func (self *Factory) Clear() {
	self.wrapped.Clear()
}

// Factory implementation
func (self *Factory) CreateCategoryService() service_category.IService {
	return category.New(
		self.factories.Category.CreateCategoryService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateDeliveryService() service_delivery.IService {
	return delivery.New(
		self.factories.Delivery.CreateDeliveryService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateDeliveryCompanyService() service_delivery.ICompanyService {
	return delivery.NewCompany(
		self.factories.DeliveryCompany.CreateDeliveryCompanyService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstanceService() service_instance.IService {
	return instance.New(
		self.factories.Instance.CreateInstanceService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstancePayPlansService() service_instance.IPayPlansService {
	return instance.NewPayPlans(
		self.factories.InstancePayPlans.CreateInstancePayPlansService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstancePhotoService() service_instance.IPhotoService {
	return instance.NewPhoto(
		self.factories.InstancePhoto.CreateInstancePhotoService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstanceReviewService() service_instance.IReviewService {
	return instance.NewReview(
		self.factories.InstanceReview.CreateInstanceReviewService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateLoginService() service_login.IService {
	return login.New(
		self.factories.Login.CreateLoginService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePayMethodService() service_payment.IPayMethodService {
	return payment.NewPayMethod(
		self.factories.PayMethod.CreatePayMethodService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserPayMethodService() service_payment.IUserPayMethodService {
	return payment.NewUserPayMethod(
		self.factories.UserPayMethod.CreateUserPayMethodService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentPaymentService() service_payment.IRentPaymentService {
	return payment.NewRentPayment(
		self.factories.RentPayment.CreateRentPaymentService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePeriodService() service_period.IService {
	return period.New(
		self.factories.Period.CreatePeriodService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePhotoService() service_photo.IService {
	return photo.New(
		self.factories.Photo.CreatePhotoService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointService() service_pickuppoint.IService {
	return pickuppoint.New(
		self.factories.PickUpPoint.CreatePickUpPointService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointPhotoService() service_pickuppoint.IPhotoService {
	return pickuppoint.NewPhoto(
		self.factories.PickUpPointPhoto.CreatePickUpPointPhotoService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointWorkingHoursService() service_pickuppoint.IWorkingHoursService {
	return pickuppoint.NewWorkingHours(
		self.factories.PickUpPointWorkingHours.CreatePickUpPointWorkingHoursService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductService() service_product.IService {
	return product.New(
		self.factories.Product.CreateProductService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductCharacteristicsService() service_product.ICharacteristicsService {
	return product.NewCharacteristics(
		self.factories.ProductCharacteristics.CreateProductCharacteristicsService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductPhotoService() service_product.IPhotoService {
	return product.NewPhoto(
		self.factories.ProductPhoto.CreateProductPhotoService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionService() service_provide.IService {
	return provide.New(
		self.factories.Provision.CreateProvisionService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionRequestService() service_provide.IRequestService {
	return provide.NewRequest(
		self.factories.ProvisionRequest.CreateProvisionRequestService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionRevokeService() service_provide.IRevokeService {
	return provide.NewRevoke(
		self.factories.ProvisionRevoke.CreateProvisionRevokeService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentService() service_rent.IService {
	return rent.New(
		self.factories.Rent.CreateRentService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentRequestService() service_rent.IRequestService {
	return rent.NewRequest(
		self.factories.RentRequest.CreateRentRequestService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentReturnService() service_rent.IReturnService {
	return rent.NewReturn(
		self.factories.RentReturn.CreateRentReturnService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateStorageService() service_storage.IService {
	return storage.New(
		self.factories.Storage.CreateStorageService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserService() service_user.IService {
	return user.New(
		self.factories.User.CreateUserService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserProfileService() service_user.IProfileService {
	return user.NewProfile(
		self.factories.UserProfile.CreateUserProfileService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserFavoriteService() service_user.IFavoriteService {
	return user.NewFavorite(
		self.factories.UserFavorite.CreateUserFavoriteService(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRoleService() service_user.IRoleService {
	return user.NewRole(
		self.factories.Role.CreateRoleService(),
		self.hl,
		self.tracer,
	)
}

