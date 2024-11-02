package review

import "rent_service/internal/repository/interfaces/review"

type MockProvider struct {
	repository review.IRepository
}

func New(repository review.IRepository) *MockProvider {
	return &MockProvider{repository}
}

func (self *MockProvider) GetReviewRepository() review.IRepository {
	return self.repository
}

