package provision

import "rent_service/internal/repository/interfaces/provision"

type MockProvider struct {
	repository provision.IRepository
}

func New(repository provision.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetProvisionRepository() provision.IRepository {
	return self.repository
}

type MockRequestProvider struct {
	repository provision.IRequestRepository
}

func NewRequest(repository provision.IRequestRepository) *MockRequestProvider {
	return &MockRequestProvider{repository}
}

func (self *MockRequestProvider) GetProvisionRequestRepository() provision.IRequestRepository {
	return self.repository
}

type MockRevokeProvider struct {
	repository provision.IRevokeRepository
}

func NewRevoke(repository provision.IRevokeRepository) *MockRevokeProvider {
	return &MockRevokeProvider{repository}
}

func (self *MockRevokeProvider) GetRevokeProvisionRepository() provision.IRevokeRepository {
	return self.repository
}

