package rent

import "rent_service/internal/repository/interfaces/rent"

type MockProvider struct {
	repository rent.IRepository
}

func New(repository rent.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetRentRepository() rent.IRepository {
	return self.repository
}

type MockRequestProvider struct {
	repository rent.IRequestRepository
}

func NewRequest(repository rent.IRequestRepository) *MockRequestProvider {
	return &MockRequestProvider{repository}
}

func (self *MockRequestProvider) GetRentRequestRepository() rent.IRequestRepository {
	return self.repository
}

type MockReturnProvider struct {
	repository rent.IReturnRepository
}

func NewReturn(repository rent.IReturnRepository) *MockReturnProvider {
	return &MockReturnProvider{repository}
}

func (self *MockReturnProvider) GetRentReturnRepository() rent.IReturnRepository {
	return self.repository
}

