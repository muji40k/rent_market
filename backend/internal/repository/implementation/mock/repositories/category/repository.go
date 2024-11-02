package category

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	getAll  func() (collection.Collection[models.Category], error)
	getPath func(leaf uuid.UUID) (collection.Collection[models.Category], error)
}

func New() *MockRepository {
	return &MockRepository{
		func() (collection.Collection[models.Category], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(leaf uuid.UUID) (collection.Collection[models.Category], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) WithGetAll(f func() (collection.Collection[models.Category], error)) *MockRepository {
	self.getAll = f
	return self
}

func (self *MockRepository) WithGetPath(f func(leaf uuid.UUID) (collection.Collection[models.Category], error)) *MockRepository {
	self.getPath = f
	return self
}

func (self *MockRepository) GetAll() (collection.Collection[models.Category], error) {
	return self.getAll()
}

func (self *MockRepository) GetPath(leaf uuid.UUID) (collection.Collection[models.Category], error) {
	return self.getPath(leaf)
}

