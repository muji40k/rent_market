package user

import "rent_service/internal/logic/services/interfaces/user"

type MockProvider struct {
	service user.IService
}

func New(service user.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetUserService() user.IService {
	return self.service
}

type MockProfileProvider struct {
	service user.IProfileService
}

func NewProfile(service user.IProfileService) *MockProfileProvider {
	return &MockProfileProvider{service}
}

func (self *MockProfileProvider) GetUserProfileService() user.IProfileService {
	return self.service
}

type MockFavoriteProvider struct {
	service user.IFavoriteService
}

func NewFavorite(service user.IFavoriteService) *MockFavoriteProvider {
	return &MockFavoriteProvider{service}
}

func (self *MockFavoriteProvider) GetUserFavoriteService() user.IFavoriteService {
	return self.service
}

type MockRoleProvider struct {
	service user.IRoleService
}

func NewRole(service user.IRoleService) *MockRoleProvider {
	return &MockRoleProvider{service}
}

func (self *MockRoleProvider) GetRoleService() user.IRoleService {
	return self.service
}

