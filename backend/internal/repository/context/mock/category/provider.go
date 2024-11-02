package category

import "rent_service/internal/repository/interfaces/category"

type MockProvider struct {
	repository category.IRepository
}

func New(repository category.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetCategoryRepository() category.IRepository {
	return self.repository
}

