package payment

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	getByInstanceId func(instanceId uuid.UUID) (Collection[models.Payment], error)
	getByRentId     func(rentId uuid.UUID) (Collection[models.Payment], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(instanceId uuid.UUID) (Collection[models.Payment], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(rentId uuid.UUID) (Collection[models.Payment], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) GetByInstanceId(instanceId uuid.UUID) (Collection[models.Payment], error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (Collection[models.Payment], error)) *MockRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockRepository) GetByRentId(rentId uuid.UUID) (Collection[models.Payment], error) {
	return self.getByRentId(rentId)
}

func (self *MockRepository) WithGetByRentId(f func(rentId uuid.UUID) (Collection[models.Payment], error)) *MockRepository {
	self.getByRentId = f
	return self
}

