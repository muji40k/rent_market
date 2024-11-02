package paymethod

import "rent_service/internal/repository/interfaces/paymethod"

type MockProvider struct {
	repository paymethod.IRepository
}

func New(repository paymethod.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetPayMethodRepository() paymethod.IRepository {
	return self.repository
}

