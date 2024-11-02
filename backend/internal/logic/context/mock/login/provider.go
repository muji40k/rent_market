package login

import "rent_service/internal/logic/services/interfaces/login"

type MockProvider struct {
	service login.IService
}

func New(service login.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetLoginService() login.IService {
	return self.service
}

