package provide

import "rent_service/internal/logic/services/interfaces/provide"

type MockProvider struct {
	service provide.IService
}

func New(service provide.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetProvisionService() provide.IService {
	return self.service
}

type MockRequestProvider struct {
	service provide.IRequestService
}

func NewRequest(service provide.IRequestService) *MockRequestProvider {
	return &MockRequestProvider{service}
}

func (self *MockRequestProvider) GetProvisionRequestService() provide.IRequestService {
	return self.service
}

type MockRevokeProvider struct {
	service provide.IRevokeService
}

func NewRevoke(service provide.IRevokeService) *MockRevokeProvider {
	return &MockRevokeProvider{service}
}

func (self *MockRevokeProvider) GetProvisionRevokeService() provide.IRevokeService {
	return self.service
}

