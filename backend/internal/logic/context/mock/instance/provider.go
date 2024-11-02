package instance

import "rent_service/internal/logic/services/interfaces/instance"

type MockProvider struct {
	service instance.IService
}

func New(service instance.IService) *MockProvider {
	return &MockProvider{service}
}

func (self *MockProvider) GetInstanceService() instance.IService {
	return self.service
}

type MockPayPlansProvider struct {
	service instance.IPayPlansService
}

func NewPayPlans(service instance.IPayPlansService) *MockPayPlansProvider {
	return &MockPayPlansProvider{service}
}

func (self *MockPayPlansProvider) GetInstancePayPlansService() instance.IPayPlansService {
	return self.service
}

type MockPhotoProvider struct {
	service instance.IPhotoService
}

func NewPhoto(service instance.IPhotoService) *MockPhotoProvider {
	return &MockPhotoProvider{service}
}

func (self *MockPhotoProvider) GetInstancePhotoService() instance.IPhotoService {
	return self.service
}

type MockReviewProvider struct {
	service instance.IReviewService
}

func NewReview(service instance.IReviewService) *MockReviewProvider {
	return &MockReviewProvider{service}
}

func (self *MockReviewProvider) GetInstanceReviewService() instance.IReviewService {
	return self.service
}

