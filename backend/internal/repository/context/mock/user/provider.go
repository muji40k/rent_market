package user

import "rent_service/internal/repository/interfaces/user"

type MockProvider struct {
	repository user.IRepository
}

func New(repository user.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetUserRepository() user.IRepository {
	return self.repository
}

type MockProfileProvider struct {
	repository user.IProfileRepository
}

func NewProfile(repository user.IProfileRepository) *MockProfileProvider {
	return &MockProfileProvider{repository}
}

func (self *MockProfileProvider) GetUserProfileRepository() user.IProfileRepository {
	return self.repository
}

type MockFavoriteProvider struct {
	repository user.IFavoriteRepository
}

func NewFavorite(repository user.IFavoriteRepository) *MockFavoriteProvider {
	return &MockFavoriteProvider{repository}
}

func (self *MockFavoriteProvider) GetUserFavoriteRepository() user.IFavoriteRepository {
	return self.repository
}

type MockPayMethodsProvider struct {
	repository user.IPayMethodsRepository
}

func NewPayMethods(repository user.IPayMethodsRepository) *MockPayMethodsProvider {
	return &MockPayMethodsProvider{repository}
}

func (self *MockPayMethodsProvider) GetUserPayMethodsRepository() user.IPayMethodsRepository {
	return self.repository
}

