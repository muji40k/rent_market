package delivery

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create                   func(delivery requests.Delivery) (requests.Delivery, error)
	update                   func(delivery requests.Delivery) error
	getById                  func(deliveryId uuid.UUID) (requests.Delivery, error)
	getActiveByPickUpPointId func(pickUpPointId uuid.UUID) (collection.Collection[requests.Delivery], error)
	getActiveByInstanceId    func(instanceId uuid.UUID) (requests.Delivery, error)
}

func New() *MockRepository {
	return &MockRepository{
		func(delivery requests.Delivery) (requests.Delivery, error) {
			return requests.Delivery{}, cmnerrors.ErrorNotSet
		},
		func(delivery requests.Delivery) error {
			return cmnerrors.ErrorNotSet
		},
		func(deliveryId uuid.UUID) (requests.Delivery, error) {
			return requests.Delivery{}, cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (collection.Collection[requests.Delivery], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (requests.Delivery, error) {
			return requests.Delivery{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(delivery requests.Delivery) (requests.Delivery, error) {
	return self.create(delivery)
}

func (self *MockRepository) WithCreate(f func(delivery requests.Delivery) (requests.Delivery, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(delivery requests.Delivery) error {
	return self.update(delivery)
}

func (self *MockRepository) WithUpdate(f func(delivery requests.Delivery) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetById(deliveryId uuid.UUID) (requests.Delivery, error) {
	return self.getById(deliveryId)
}

func (self *MockRepository) WithGetById(f func(deliveryId uuid.UUID) (requests.Delivery, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetActiveByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Delivery], error) {
	return self.getActiveByPickUpPointId(pickUpPointId)
}

func (self *MockRepository) WithGetActiveByPickUpPointId(f func(pickUpPointId uuid.UUID) (collection.Collection[requests.Delivery], error)) *MockRepository {
	self.getActiveByPickUpPointId = f
	return self
}

func (self *MockRepository) GetActiveByInstanceId(instanceId uuid.UUID) (requests.Delivery, error) {
	return self.getActiveByInstanceId(instanceId)
}

func (self *MockRepository) WithGetActiveByInstanceId(f func(instanceId uuid.UUID) (requests.Delivery, error)) *MockRepository {
	self.getActiveByInstanceId = f
	return self
}

type MockCompanyRepository struct {
	getById func(companyId uuid.UUID) (models.DeliveryCompany, error)
	getAll  func() (collection.Collection[models.DeliveryCompany], error)
}

func NewCompany() *MockCompanyRepository {
	return &MockCompanyRepository{
		func(companyId uuid.UUID) (models.DeliveryCompany, error) {
			return models.DeliveryCompany{}, cmnerrors.ErrorNotSet
		},
		func() (collection.Collection[models.DeliveryCompany], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockCompanyRepository) GetById(companyId uuid.UUID) (models.DeliveryCompany, error) {
	return self.getById(companyId)
}

func (self *MockCompanyRepository) WithGetById(f func(companyId uuid.UUID) (models.DeliveryCompany, error)) *MockCompanyRepository {
	self.getById = f
	return self
}

func (self *MockCompanyRepository) GetAll() (collection.Collection[models.DeliveryCompany], error) {
	return self.getAll()
}

func (self *MockCompanyRepository) WithGetAll(f func() (collection.Collection[models.DeliveryCompany], error)) *MockCompanyRepository {
	self.getAll = f
	return self
}

