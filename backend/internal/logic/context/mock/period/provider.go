package period

import "rent_service/internal/logic/services/interfaces/period"

type MockProvider struct {
	service period.IService
}

func New(service period.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetPeriodService() period.IService {
	return self.service
}

