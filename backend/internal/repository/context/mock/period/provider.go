package period

import "rent_service/internal/repository/interfaces/period"

type MockProvider struct {
	repository period.IRepository
}

func New(repository period.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetPeriodRepository() period.IRepository {
	return self.repository
}

