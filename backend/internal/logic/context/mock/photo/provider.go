package photo

import "rent_service/internal/logic/services/interfaces/photo"

type MockProvider struct {
	service photo.IService
}

func New(service photo.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetPhotoService() photo.IService {
	return self.service
}

