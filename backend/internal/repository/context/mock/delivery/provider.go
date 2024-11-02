package delivery

import "rent_service/internal/repository/interfaces/delivery"

type MockProvider struct {
	repository delivery.IRepository
}

func New(repository delivery.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetDeliveryRepository() delivery.IRepository {
	return self.repository
}

type MockCompanyProvider struct {
	repository delivery.ICompanyRepository
}

func NewCompany(repository delivery.ICompanyRepository) *MockCompanyProvider {
	return &MockCompanyProvider{repository}
}

func (self *MockCompanyProvider) GetDeliveryCompanyRepository() delivery.ICompanyRepository {
	return self.repository
}

