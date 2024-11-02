package period

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	getById func(periodId uuid.UUID) (models.Period, error)
	getAll  func() (collection.Collection[models.Period], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(periodId uuid.UUID) (models.Period, error) {
			return models.Period{}, cmnerrors.ErrorNotSet
		},
		func() (collection.Collection[models.Period], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) GetById(periodId uuid.UUID) (models.Period, error) {
	return self.getById(periodId)
}

func (self *MockRepository) WithGetById(f func(periodId uuid.UUID) (models.Period, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetAll() (collection.Collection[models.Period], error) {
	return self.getAll()
}

func (self *MockRepository) WithGetAll(f func() (collection.Collection[models.Period], error)) *MockRepository {
	self.getAll = f
	return self
}

