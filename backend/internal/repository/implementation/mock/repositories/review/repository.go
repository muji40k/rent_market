package review

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"
	"rent_service/internal/repository/interfaces/review"
)

type MockRepository struct {
	create        func(review models.Review) (models.Review, error)
	getWithFilter func(filter review.Filter, sort review.Sort) (Collection[models.Review], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(review models.Review) (models.Review, error) {
			return models.Review{}, cmnerrors.ErrorNotSet
		},
		func(filter review.Filter, sort review.Sort) (Collection[models.Review], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(review models.Review) (models.Review, error) {
	return self.create(review)
}

func (self *MockRepository) WithCreate(f func(review models.Review) (models.Review, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) GetWithFilter(filter review.Filter, sort review.Sort) (Collection[models.Review], error) {
	return self.getWithFilter(filter, sort)
}

func (self *MockRepository) WithGetWithFilter(f func(filter review.Filter, sort review.Sort) (Collection[models.Review], error)) *MockRepository {
	self.getWithFilter = f
	return self
}

