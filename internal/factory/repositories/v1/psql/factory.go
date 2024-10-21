package psql

import (
	cv1 "rent_service/internal/repository/context/v1"
	"rent_service/internal/repository/implementation/sql/repositories/address"
	"rent_service/internal/repository/implementation/sql/repositories/category"
	"rent_service/internal/repository/implementation/sql/repositories/currency"
	"rent_service/internal/repository/implementation/sql/repositories/delivery"
	"rent_service/internal/repository/implementation/sql/repositories/instance"
	"rent_service/internal/repository/implementation/sql/repositories/payment"
	"rent_service/internal/repository/implementation/sql/repositories/paymethod"
	"rent_service/internal/repository/implementation/sql/repositories/period"
	"rent_service/internal/repository/implementation/sql/repositories/photo"
	"rent_service/internal/repository/implementation/sql/repositories/pickuppoint"
	"rent_service/internal/repository/implementation/sql/repositories/product"
	"rent_service/internal/repository/implementation/sql/repositories/provision"
	"rent_service/internal/repository/implementation/sql/repositories/rent"
	"rent_service/internal/repository/implementation/sql/repositories/review"
	"rent_service/internal/repository/implementation/sql/repositories/role"
	"rent_service/internal/repository/implementation/sql/repositories/storage"
	"rent_service/internal/repository/implementation/sql/repositories/user"
	"rent_service/internal/repository/implementation/sql/technical"
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

	"github.com/jmoiron/sqlx"
)

type innerRepos struct {
	currency *currency.Repository
	address  *address.Repository
}

type Factory struct {
	connection *sqlx.DB
	setter     technical.ISetter
	hasher     user.Hasher
	repos      innerRepos
}

func New(
	connection *sqlx.DB,
	setter technical.ISetter,
	hasher user.Hasher,
) *Factory {
	return &Factory{connection, setter, hasher, innerRepos{}}
}

func (self *Factory) CreateAddressRepository() *address.Repository {
	if nil == self.repos.address {
		self.repos.address = address.New(self.connection)
	}

	return self.repos.address
}

func (self *Factory) CreateCurrencyRepository() *currency.Repository {
	if nil == self.repos.currency {
		self.repos.currency = currency.New(self.connection)
	}

	return self.repos.currency
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

func (self *Factory) CreateCategoryRepository() repository_category.IRepository {
	return category.New(self.connection)
}

func (self *Factory) CreateDeliveryRepository() repository_delivery.IRepository {
	return delivery.New(self.connection, self.setter)
}

func (self *Factory) CreateDeliveryCompanyRepository() repository_delivery.ICompanyRepository {
	return delivery.NewCompany(self.connection)
}

func (self *Factory) CreateInstanceRepository() repository_instance.IRepository {
	return instance.New(self.connection, self.setter)
}

func (self *Factory) CreateInstancePayPlansRepository() repository_instance.IPayPlansRepository {
	return instance.NewPayPlans(
		self.connection,
		self.setter,
		self.createCurrencyRepository(),
	)
}

func (self *Factory) CreateInstancePhotoRepository() repository_instance.IPhotoRepository {
	return instance.NewPhoto(self.connection, self.setter)
}

func (self *Factory) CreatePaymentRepository() repository_payment.IRepository {
	return payment.New(self.connection, self.createCurrencyRepository())
}

func (self *Factory) CreatePayMethodRepository() repository_paymethod.IRepository {
	return paymethod.New(self.connection)
}

func (self *Factory) CreatePeriodRepository() repository_period.IRepository {
	return period.New(self.connection)
}

func (self *Factory) CreatePhotoRepository() repository_photo.IRepository {
	return photo.New(self.connection, self.setter)
}

func (self *Factory) CreatePhotoTempRepository() repository_photo.ITempRepository {
	return photo.NewTemp(self.connection, self.setter)
}

func (self *Factory) CreatePickUpPointRepository() repository_pickuppoint.IRepository {
	return pickuppoint.New(self.connection, self.createAddressRepository())
}

func (self *Factory) CreatePickUpPointPhotoRepository() repository_pickuppoint.IPhotoRepository {
	return pickuppoint.NewPhoto(self.connection)
}

func (self *Factory) CreatePickUpPointWorkingHoursRepository() repository_pickuppoint.IWorkingHoursRepository {
	return pickuppoint.NewWorkingHours(self.connection)
}

func (self *Factory) CreateProductRepository() repository_product.IRepository {
	return product.New(self.connection)
}

func (self *Factory) CreateProductCharacteristicsRepository() repository_product.ICharacteristicsRepository {
	return product.NewCharacteristics(self.connection)
}

func (self *Factory) CreateProductPhotoRepository() repository_product.IPhotoRepository {
	return product.NewPhoto(self.connection)
}

func (self *Factory) CreateProvisionRepository() repository_provision.IRepository {
	return provision.New(self.connection, self.setter)
}

func (self *Factory) CreateProvisionRequestRepository() repository_provision.IRequestRepository {
	return provision.NewRequest(
		self.connection,
		self.setter,
		self.createCurrencyRepository(),
	)
}

func (self *Factory) CreateProvisionRevokeRepository() repository_provision.IRevokeRepository {
	return provision.NewRevoke(self.connection, self.setter)
}

func (self *Factory) CreateRentRepository() repository_rent.IRepository {
	return rent.New(self.connection, self.setter)
}

func (self *Factory) CreateRentRequestRepository() repository_rent.IRequestRepository {
	return rent.NewRequest(self.connection, self.setter)
}

func (self *Factory) CreateRentReturnRepository() repository_rent.IReturnRepository {
	return rent.NewReturn(self.connection, self.setter)
}

func (self *Factory) CreateReviewRepository() repository_review.IRepository {
	return review.New(self.connection, self.setter)
}

func (self *Factory) CreateRoleAdministratorRepository() repository_role.IAdministratorRepository {
	return role.NewAdministrator(self.connection)
}

func (self *Factory) CreateRoleRenterRepository() repository_role.IRenterRepository {
	return role.NewRenter(self.connection, self.setter)
}

func (self *Factory) CreateRoleStorekeeperRepository() repository_role.IStorekeeperRepository {
	return role.NewStorekeeper(self.connection)
}

func (self *Factory) CreateStorageRepository() repository_storage.IRepository {
	return storage.New(self.connection, self.setter)
}

func (self *Factory) CreateUserRepository() repository_user.IRepository {
	return user.New(self.connection, self.setter, self.hasher)
}

func (self *Factory) CreateUserProfileRepository() repository_user.IProfileRepository {
	return user.NewProfile(self.connection, self.setter)
}

func (self *Factory) CreateUserFavoriteRepository() repository_user.IFavoriteRepository {
	return user.NewFavorite(self.connection, self.setter)
}

func (self *Factory) CreateUserPayMethodsRepository() repository_user.IPayMethodsRepository {
	return user.NewPayMethod(self.connection, self.setter)
}

