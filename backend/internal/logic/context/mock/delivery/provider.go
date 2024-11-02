package delivery

import "rent_service/internal/logic/services/interfaces/delivery"

type MockProvider struct {
	service delivery.IService
}

func New(service delivery.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetDeliveryService() delivery.IService {
	return self.service
}

type MockCompanyProvider struct {
	service delivery.ICompanyService
}

func NewCompany(service delivery.ICompanyService) *MockCompanyProvider {
	return &MockCompanyProvider{service}
}

func (self *MockCompanyProvider) GetDeliveryCompanyService() delivery.ICompanyService {
	return self.service
}

