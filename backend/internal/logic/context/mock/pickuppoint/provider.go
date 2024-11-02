package pickuppoint

import "rent_service/internal/logic/services/interfaces/pickuppoint"

type MockProvider struct {
	service pickuppoint.IService
}

func New(service pickuppoint.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetPickUpPointService() pickuppoint.IService {
	return self.service
}

type MockPhotoProvider struct {
	service pickuppoint.IPhotoService
}

func NewPhoto(service pickuppoint.IPhotoService) *MockPhotoProvider {
	return &MockPhotoProvider{service}
}

func (self *MockPhotoProvider) GetPickUpPointPhotoService() pickuppoint.IPhotoService {
	return self.service
}

type MockWorkingHoursProvider struct {
	service pickuppoint.IWorkingHoursService
}

func NewWorkingHours(service pickuppoint.IWorkingHoursService) *MockWorkingHoursProvider {
	return &MockWorkingHoursProvider{service}
}

func (self *MockWorkingHoursProvider) GetPickUpPointWorkingHoursService() pickuppoint.IWorkingHoursService {
	return self.service
}

