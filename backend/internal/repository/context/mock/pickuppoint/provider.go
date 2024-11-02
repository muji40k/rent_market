package pickuppoint

import "rent_service/internal/repository/interfaces/pickuppoint"

type MockProvider struct {
	repository pickuppoint.IRepository
}

func New(repository pickuppoint.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetPickUpPointRepository() pickuppoint.IRepository {
	return self.repository
}

type MockPhotoProvider struct {
	repository pickuppoint.IPhotoRepository
}

func NewPhoto(repository pickuppoint.IPhotoRepository) *MockPhotoProvider {
	return &MockPhotoProvider{repository}
}

func (self *MockPhotoProvider) GetPickUpPointPhotoRepository() pickuppoint.IPhotoRepository {
	return self.repository
}

type MockWorkingHoursProvider struct {
	repository pickuppoint.IWorkingHoursRepository
}

func NewWorkingHours(repository pickuppoint.IWorkingHoursRepository) *MockWorkingHoursProvider {
	return &MockWorkingHoursProvider{repository}
}

func (self *MockWorkingHoursProvider) GetPickUpPointWorkingHoursRepository() pickuppoint.IWorkingHoursRepository {
	return self.repository
}

