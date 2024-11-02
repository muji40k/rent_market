package photo

import "rent_service/internal/repository/interfaces/photo"

type MockProvider struct {
	repository photo.IRepository
}

func New(repository photo.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetPhotoRepository() photo.IRepository {
	return self.repository
}

type MockTempProvider struct {
	repository photo.ITempRepository
}

func NewTemp(repository photo.ITempRepository) *MockTempProvider {
	return &MockTempProvider{repository}
}

func (self *MockTempProvider) GetTempPhotoRepository() photo.ITempRepository {
	return self.repository
}

