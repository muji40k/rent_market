package storage

import "rent_service/internal/logic/services/interfaces/storage"

type MockProvider struct {
	service storage.IService
}

func New(service storage.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetStorageService() storage.IService {
	return self.service
}

