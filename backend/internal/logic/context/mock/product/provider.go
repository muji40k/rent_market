package product

import "rent_service/internal/logic/services/interfaces/product"

type MockProvider struct {
	service product.IService
}

func New(service product.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetProductService() product.IService {
	return self.service
}

type MockCharacteristicsProvider struct {
	service product.ICharacteristicsService
}

func NewCharacteristics(service product.ICharacteristicsService) *MockCharacteristicsProvider {
	return &MockCharacteristicsProvider{service}
}

func (self *MockCharacteristicsProvider) GetProductCharacteristicsService() product.ICharacteristicsService {
	return self.service
}

type MockPhotoProvider struct {
	service product.IPhotoService
}

func NewPhoto(service product.IPhotoService) *MockPhotoProvider {
	return &MockPhotoProvider{service}
}

func (self *MockPhotoProvider) GetProductPhotoService() product.IPhotoService {
	return self.service
}

