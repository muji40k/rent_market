package role

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockAdministratorRepository struct {
	getByUserId func(userId uuid.UUID) (models.Administrator, error)
}

func NewAdministrator() *MockAdministratorRepository {
	return &MockAdministratorRepository{
		func(userId uuid.UUID) (models.Administrator, error) {
			return models.Administrator{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockAdministratorRepository) GetByUserId(userId uuid.UUID) (models.Administrator, error) {
	return self.getByUserId(userId)
}

func (self *MockAdministratorRepository) WithGetByUserId(f func(userId uuid.UUID) (models.Administrator, error)) *MockAdministratorRepository {
	self.getByUserId = f
	return self
}

type MockRenterRepository struct {
	create      func(userId uuid.UUID) (models.Renter, error)
	getById     func(renterId uuid.UUID) (models.Renter, error)
	getByUserId func(userId uuid.UUID) (models.Renter, error)
}

func NewRenter() *MockRenterRepository {
	return &MockRenterRepository{
		func(userId uuid.UUID) (models.Renter, error) {
			return models.Renter{}, cmnerrors.ErrorNotSet
		},
		func(renterId uuid.UUID) (models.Renter, error) {
			return models.Renter{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (models.Renter, error) {
			return models.Renter{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRenterRepository) Create(userId uuid.UUID) (models.Renter, error) {
	return self.create(userId)
}

func (self *MockRenterRepository) WithCreate(f func(userId uuid.UUID) (models.Renter, error)) *MockRenterRepository {
	self.create = f
	return self
}

func (self *MockRenterRepository) GetById(renterId uuid.UUID) (models.Renter, error) {
	return self.getById(renterId)
}

func (self *MockRenterRepository) WithGetById(f func(renterId uuid.UUID) (models.Renter, error)) *MockRenterRepository {
	self.getById = f
	return self
}

func (self *MockRenterRepository) GetByUserId(userId uuid.UUID) (models.Renter, error) {
	return self.getByUserId(userId)
}

func (self *MockRenterRepository) WithGetByUserId(f func(userId uuid.UUID) (models.Renter, error)) *MockRenterRepository {
	self.getByUserId = f
	return self
}

type MockStorekeeperRepository struct {
	getByUserId func(userId uuid.UUID) (models.Storekeeper, error)
}

func NewStorekeeper() *MockStorekeeperRepository {
	return &MockStorekeeperRepository{
		func(userId uuid.UUID) (models.Storekeeper, error) {
			return models.Storekeeper{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockStorekeeperRepository) GetByUserId(userId uuid.UUID) (models.Storekeeper, error) {
	return self.getByUserId(userId)
}

func (self *MockStorekeeperRepository) WithGetByUserId(f func(userId uuid.UUID) (models.Storekeeper, error)) *MockStorekeeperRepository {
	self.getByUserId = f
	return self
}

