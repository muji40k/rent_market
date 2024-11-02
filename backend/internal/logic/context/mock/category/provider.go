package category

import "rent_service/internal/logic/services/interfaces/category"

type MockProvider struct {
	service category.IService
}

func New(service category.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetCategoryService() category.IService {
	return self.service
}

