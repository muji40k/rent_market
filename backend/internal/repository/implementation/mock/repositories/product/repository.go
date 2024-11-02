package product

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"
	"rent_service/internal/repository/interfaces/product"

	"github.com/google/uuid"
)

type MockRepository struct {
	getById       func(productId uuid.UUID) (models.Product, error)
	getWithFilter func(filter product.Filter, sort product.Sort) (collection.Collection[models.Product], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(productId uuid.UUID) (models.Product, error) {
			return models.Product{}, cmnerrors.ErrorNotSet
		},
		func(filter product.Filter, sort product.Sort) (collection.Collection[models.Product], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) GetById(productId uuid.UUID) (models.Product, error) {
	return self.getById(productId)
}

func (self *MockRepository) WithGetById(f func(productId uuid.UUID) (models.Product, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetWithFilter(filter product.Filter, sort product.Sort) (collection.Collection[models.Product], error) {
	return self.getWithFilter(filter, sort)
}

func (self *MockRepository) WithGetWithFilter(f func(filter product.Filter, sort product.Sort) (collection.Collection[models.Product], error)) *MockRepository {
	self.getWithFilter = f
	return self
}

type MockCharacteristicsRepository struct {
	getByProductId func(productId uuid.UUID) (models.ProductCharacteristics, error)
}

func NewCharacteristics() *MockCharacteristicsRepository {
	return &MockCharacteristicsRepository{
		func(productId uuid.UUID) (models.ProductCharacteristics, error) {
			return models.ProductCharacteristics{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockCharacteristicsRepository) GetByProductId(productId uuid.UUID) (models.ProductCharacteristics, error) {
	return self.getByProductId(productId)
}

func (self *MockCharacteristicsRepository) WithGetByProductId(f func(productId uuid.UUID) (models.ProductCharacteristics, error)) *MockCharacteristicsRepository {
	self.getByProductId = f
	return self
}

type MockPhotoRepository struct {
	getByProductId func(productId uuid.UUID) (collection.Collection[uuid.UUID], error)
}

func NewPhoto() *MockPhotoRepository {
	return &MockPhotoRepository{
		func(productId uuid.UUID) (collection.Collection[uuid.UUID], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockPhotoRepository) GetByProductId(productId uuid.UUID) (collection.Collection[uuid.UUID], error) {
	return self.getByProductId(productId)
}

func (self *MockPhotoRepository) WithGetByProductId(f func(productId uuid.UUID) (collection.Collection[uuid.UUID], error)) *MockPhotoRepository {
	self.getByProductId = f
	return self
}

