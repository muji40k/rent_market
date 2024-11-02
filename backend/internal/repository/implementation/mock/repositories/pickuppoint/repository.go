package pickuppoint

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	getById func(pickUpPointId uuid.UUID) (models.PickUpPoint, error)
	getAll  func() (collection.Collection[models.PickUpPoint], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(pickUpPointId uuid.UUID) (models.PickUpPoint, error) {
			return models.PickUpPoint{}, cmnerrors.ErrorNotSet
		},
		func() (collection.Collection[models.PickUpPoint], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) GetById(pickUpPointId uuid.UUID) (models.PickUpPoint, error) {
	return self.getById(pickUpPointId)
}

func (self *MockRepository) WithGetById(f func(pickUpPointId uuid.UUID) (models.PickUpPoint, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetAll() (collection.Collection[models.PickUpPoint], error) {
	return self.getAll()
}

func (self *MockRepository) WithGetAll(f func() (collection.Collection[models.PickUpPoint], error)) *MockRepository {
	self.getAll = f
	return self
}

type MockPhotoRepository struct {
	getById func(pickUpPointId uuid.UUID) (collection.Collection[uuid.UUID], error)
}

func NewPhoto() *MockPhotoRepository {
	return &MockPhotoRepository{
		func(pickUpPointId uuid.UUID) (collection.Collection[uuid.UUID], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockPhotoRepository) GetById(pickUpPointId uuid.UUID) (collection.Collection[uuid.UUID], error) {
	return self.getById(pickUpPointId)
}

func (self *MockPhotoRepository) WithGetById(f func(pickUpPointId uuid.UUID) (collection.Collection[uuid.UUID], error)) *MockPhotoRepository {
	self.getById = f
	return self
}

type MockWorkingHoursRepository struct {
	getById func(pickUpPointId uuid.UUID) (models.PickUpPointWorkingHours, error)
}

func NewWorkingHours() *MockWorkingHoursRepository {
	return &MockWorkingHoursRepository{
		func(pickUpPointId uuid.UUID) (models.PickUpPointWorkingHours, error) {
			return models.PickUpPointWorkingHours{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockWorkingHoursRepository) GetById(pickUpPointId uuid.UUID) (models.PickUpPointWorkingHours, error) {
	return self.getById(pickUpPointId)
}

func (self *MockWorkingHoursRepository) WithGetById(f func(pickUpPointId uuid.UUID) (models.PickUpPointWorkingHours, error)) *MockWorkingHoursRepository {
	self.getById = f
	return self
}

