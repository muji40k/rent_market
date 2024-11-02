package storage

import "rent_service/internal/repository/interfaces/storage"

type MockProvider struct {
	repository storage.IRepository
}

func New(repository storage.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetStorageRepository() storage.IRepository {
	return self.repository
}

