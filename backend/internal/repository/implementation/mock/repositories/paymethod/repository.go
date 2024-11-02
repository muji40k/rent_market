package paymethod

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"
)

type MockRepository struct {
	getAll func() (collection.Collection[models.PayMethod], error)
}

func New() *MockRepository {
	return &MockRepository{
		func() (collection.Collection[models.PayMethod], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) GetAll() (collection.Collection[models.PayMethod], error) {
	return self.getAll()
}

func (self *MockRepository) WithGetAll(f func() (collection.Collection[models.PayMethod], error)) *MockRepository {
	self.getAll = f
	return self
}

