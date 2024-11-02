package payment

import "rent_service/internal/repository/interfaces/payment"

type MockProvider struct {
	repository payment.IRepository
}

func New(repository payment.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetPaymentRepository() payment.IRepository {
	return self.repository
}

