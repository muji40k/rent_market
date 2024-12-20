package tracer

import (
	fv1 "rent_service/internal/factory/repositories/v1"
	cv1 "rent_service/internal/repository/context/v1"
	"rent_service/internal/repository/implementation/tracer/category"
	"rent_service/internal/repository/implementation/tracer/delivery"
	"rent_service/internal/repository/implementation/tracer/instance"
	"rent_service/internal/repository/implementation/tracer/payment"
	"rent_service/internal/repository/implementation/tracer/paymethod"
	"rent_service/internal/repository/implementation/tracer/period"
	"rent_service/internal/repository/implementation/tracer/photo"
	"rent_service/internal/repository/implementation/tracer/pickuppoint"
	"rent_service/internal/repository/implementation/tracer/product"
	"rent_service/internal/repository/implementation/tracer/provision"
	"rent_service/internal/repository/implementation/tracer/rent"
	"rent_service/internal/repository/implementation/tracer/review"
	"rent_service/internal/repository/implementation/tracer/role"
	"rent_service/internal/repository/implementation/tracer/storage"
	"rent_service/internal/repository/implementation/tracer/user"
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
		Payment:                 self,
		PayMethod:               self,
		Period:                  self,
		Photo:                   self,
		PhotoTemp:               self,
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
		Review:                  self,
		RoleAdministrator:       self,
		RoleRenter:              self,
		RoleStorekeeper:         self,
		Storage:                 self,
		User:                    self,
		UserProfile:             self,
		UserFavorite:            self,
		UserPayMethods:          self,
	}
}

func (self *Factory) Clear() {
	self.wrapped.Clear()
}

func (self *Factory) CreateCategoryRepository() repository_category.IRepository {
	return category.New(
		self.factories.Category.CreateCategoryRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateDeliveryRepository() repository_delivery.IRepository {
	return delivery.New(
		self.factories.Delivery.CreateDeliveryRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateDeliveryCompanyRepository() repository_delivery.ICompanyRepository {
	return delivery.NewCompany(
		self.factories.DeliveryCompany.CreateDeliveryCompanyRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstanceRepository() repository_instance.IRepository {
	return instance.New(
		self.factories.Instance.CreateInstanceRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstancePayPlansRepository() repository_instance.IPayPlansRepository {
	return instance.NewPayPlans(
		self.factories.InstancePayPlans.CreateInstancePayPlansRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateInstancePhotoRepository() repository_instance.IPhotoRepository {
	return instance.NewPhoto(
		self.factories.InstancePhoto.CreateInstancePhotoRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePaymentRepository() repository_payment.IRepository {
	return payment.New(
		self.factories.Payment.CreatePaymentRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePayMethodRepository() repository_paymethod.IRepository {
	return paymethod.New(
		self.factories.PayMethod.CreatePayMethodRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePeriodRepository() repository_period.IRepository {
	return period.New(
		self.factories.Period.CreatePeriodRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePhotoRepository() repository_photo.IRepository {
	return photo.New(
		self.factories.Photo.CreatePhotoRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePhotoTempRepository() repository_photo.ITempRepository {
	return photo.NewTemp(
		self.factories.PhotoTemp.CreatePhotoTempRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointRepository() repository_pickuppoint.IRepository {
	return pickuppoint.New(
		self.factories.PickUpPoint.CreatePickUpPointRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointPhotoRepository() repository_pickuppoint.IPhotoRepository {
	return pickuppoint.NewPhoto(
		self.factories.PickUpPointPhoto.CreatePickUpPointPhotoRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreatePickUpPointWorkingHoursRepository() repository_pickuppoint.IWorkingHoursRepository {
	return pickuppoint.NewWorkingHours(
		self.factories.PickUpPointWorkingHours.CreatePickUpPointWorkingHoursRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductRepository() repository_product.IRepository {
	return product.New(
		self.factories.Product.CreateProductRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductCharacteristicsRepository() repository_product.ICharacteristicsRepository {
	return product.NewCharacteristics(
		self.factories.ProductCharacteristics.CreateProductCharacteristicsRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProductPhotoRepository() repository_product.IPhotoRepository {
	return product.NewPhoto(
		self.factories.ProductPhoto.CreateProductPhotoRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionRepository() repository_provision.IRepository {
	return provision.New(
		self.factories.Provision.CreateProvisionRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionRequestRepository() repository_provision.IRequestRepository {
	return provision.NewRequest(
		self.factories.ProvisionRequest.CreateProvisionRequestRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateProvisionRevokeRepository() repository_provision.IRevokeRepository {
	return provision.NewRevoke(
		self.factories.ProvisionRevoke.CreateProvisionRevokeRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentRepository() repository_rent.IRepository {
	return rent.New(
		self.factories.Rent.CreateRentRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentRequestRepository() repository_rent.IRequestRepository {
	return rent.NewRequest(
		self.factories.RentRequest.CreateRentRequestRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRentReturnRepository() repository_rent.IReturnRepository {
	return rent.NewReturn(
		self.factories.RentReturn.CreateRentReturnRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateReviewRepository() repository_review.IRepository {
	return review.New(
		self.factories.Review.CreateReviewRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRoleAdministratorRepository() repository_role.IAdministratorRepository {
	return role.NewAdministrator(
		self.factories.RoleAdministrator.CreateRoleAdministratorRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRoleRenterRepository() repository_role.IRenterRepository {
	return role.NewRenter(
		self.factories.RoleRenter.CreateRoleRenterRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateRoleStorekeeperRepository() repository_role.IStorekeeperRepository {
	return role.NewStorekeeper(
		self.factories.RoleStorekeeper.CreateRoleStorekeeperRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateStorageRepository() repository_storage.IRepository {
	return storage.New(
		self.factories.Storage.CreateStorageRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserRepository() repository_user.IRepository {
	return user.New(
		self.factories.User.CreateUserRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserProfileRepository() repository_user.IProfileRepository {
	return user.NewProfile(
		self.factories.UserProfile.CreateUserProfileRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserFavoriteRepository() repository_user.IFavoriteRepository {
	return user.NewFavorite(
		self.factories.UserFavorite.CreateUserFavoriteRepository(),
		self.hl,
		self.tracer,
	)
}

func (self *Factory) CreateUserPayMethodsRepository() repository_user.IPayMethodsRepository {
	return user.NewPayMethods(
		self.factories.UserPayMethods.CreateUserPayMethodsRepository(),
		self.hl,
		self.tracer,
	)
}

