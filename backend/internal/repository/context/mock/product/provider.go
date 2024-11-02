package product

import "rent_service/internal/repository/interfaces/product"

type MockProvider struct {
	repository product.IRepository
}

func New(repository product.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetProductRepository() product.IRepository {
	return self.repository
}

type MockCharacteristicsProvider struct {
	repository product.ICharacteristicsRepository
}

func NewCharacteristics(repository product.ICharacteristicsRepository) *MockCharacteristicsProvider {
	return &MockCharacteristicsProvider{repository}
}

func (self *MockCharacteristicsProvider) GetProductCharacteristicsRepository() product.ICharacteristicsRepository {
	return self.repository
}

type MockPhotoProvider struct {
	repository product.IPhotoRepository
}

func NewPhoto(repository product.IPhotoRepository) *MockPhotoProvider {
	return &MockPhotoProvider{repository}
}

func (self *MockPhotoProvider) GetProductPhotoRepository() product.IPhotoRepository {
	return self.repository
}

