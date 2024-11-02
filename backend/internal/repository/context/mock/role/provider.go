package role

import "rent_service/internal/repository/interfaces/role"

type MockAdministratorProvider struct {
	repository role.IAdministratorRepository
}

func NewAdministrator(repository role.IAdministratorRepository) *MockAdministratorProvider {
	return &MockAdministratorProvider{repository}
}

func (self *MockAdministratorProvider) GetAdministratorRepository() role.IAdministratorRepository {
	return self.repository
}

type MockRenterProvider struct {
	repository role.IRenterRepository
}

func NewRenter(repository role.IRenterRepository) *MockRenterProvider {
	return &MockRenterProvider{repository}
}

func (self *MockRenterProvider) GetRenterRepository() role.IRenterRepository {
	return self.repository
}

type MockStorekeeperProvider struct {
	repository role.IStorekeeperRepository
}

func NewStorekeeper(repository role.IStorekeeperRepository) *MockStorekeeperProvider {
	return &MockStorekeeperProvider{repository}
}

func (self *MockStorekeeperProvider) GetStorekeeperRepository() role.IStorekeeperRepository {
	return self.repository
}

