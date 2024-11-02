package rent

import "rent_service/internal/logic/services/interfaces/rent"

type MockProvider struct {
	service rent.IService
}

func New(service rent.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetRentService() rent.IService {
	return self.service
}

type MockRequestProvider struct {
	service rent.IRequestService
}

func NewRequest(service rent.IRequestService) *MockRequestProvider {
	return &MockRequestProvider{service}
}

func (self *MockRequestProvider) GetRentRequestService() rent.IRequestService {
	return self.service
}

type MockReturnProvider struct {
	service rent.IReturnService
}

func NewReturn(service rent.IReturnService) *MockReturnProvider {
	return &MockReturnProvider{service}
}

func (self *MockReturnProvider) GetRentReturnService() rent.IReturnService {
	return self.service
}

