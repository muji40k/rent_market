package instance

import "rent_service/internal/repository/interfaces/instance"

type MockProvider struct {
	repository instance.IRepository
}

func New(repository instance.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetInstanceRepository() instance.IRepository {
	return self.repository
}

type MockPayPlansProvider struct {
	repository instance.IPayPlansRepository
}

func NewPayPlans(repository instance.IPayPlansRepository) *MockPayPlansProvider {
	return &MockPayPlansProvider{repository}
}

func (self *MockPayPlansProvider) GetInstancePayPlansRepository() instance.IPayPlansRepository {
	return self.repository
}

type MockPhotoProvider struct {
	repository instance.IPhotoRepository
}

func NewPhoto(repository instance.IPhotoRepository) *MockPhotoProvider {
	return &MockPhotoProvider{repository}
}

func (self *MockPhotoProvider) GetInstancePhotoRepository() instance.IPhotoRepository {
	return self.repository
}

