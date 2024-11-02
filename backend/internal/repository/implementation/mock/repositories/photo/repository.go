package photo

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create  func(photo models.Photo) (models.Photo, error)
	getById func(photoId uuid.UUID) (models.Photo, error)
}

func New() *MockRepository {
	return &MockRepository{
		func(photo models.Photo) (models.Photo, error) {
			return models.Photo{}, cmnerrors.ErrorNotSet
		},
		func(photoId uuid.UUID) (models.Photo, error) {
			return models.Photo{}, cmnerrors.ErrorNotSet

		},
	}
}

func (self *MockRepository) Create(photo models.Photo) (models.Photo, error) {
	return self.create(photo)
}

func (self *MockRepository) WithCreate(f func(photo models.Photo) (models.Photo, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) GetById(photoId uuid.UUID) (models.Photo, error) {
	return self.getById(photoId)
}

func (self *MockRepository) WithGetById(f func(photoId uuid.UUID) (models.Photo, error)) *MockRepository {
	self.getById = f
	return self
}

type MockTempRepository struct {
	create  func(photo models.TempPhoto) (models.TempPhoto, error)
	update  func(photo models.TempPhoto) error
	getById func(photoId uuid.UUID) (models.TempPhoto, error)
	remove  func(photoId uuid.UUID) error
}

func NewTemp() *MockTempRepository {
	return &MockTempRepository{
		func(photo models.TempPhoto) (models.TempPhoto, error) {
			return models.TempPhoto{}, cmnerrors.ErrorNotSet
		},
		func(photo models.TempPhoto) error {
			return cmnerrors.ErrorNotSet
		},
		func(photoId uuid.UUID) (models.TempPhoto, error) {
			return models.TempPhoto{}, cmnerrors.ErrorNotSet
		},
		func(photoId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockTempRepository) Create(photo models.TempPhoto) (models.TempPhoto, error) {
	return self.create(photo)
}

func (self *MockTempRepository) WithCreate(f func(photo models.TempPhoto) (models.TempPhoto, error)) *MockTempRepository {
	self.create = f
	return self
}

func (self *MockTempRepository) Update(photo models.TempPhoto) error {
	return self.update(photo)
}

func (self *MockTempRepository) WithUpdate(f func(photo models.TempPhoto) error) *MockTempRepository {
	self.update = f
	return self
}

func (self *MockTempRepository) GetById(photoId uuid.UUID) (models.TempPhoto, error) {
	return self.getById(photoId)
}

func (self *MockTempRepository) WithGetById(f func(photoId uuid.UUID) (models.TempPhoto, error)) *MockTempRepository {
	self.getById = f
	return self
}

func (self *MockTempRepository) Remove(photoId uuid.UUID) error {
	return self.remove(photoId)
}

func (self *MockTempRepository) WithRemove(f func(photoId uuid.UUID) error) *MockTempRepository {
	self.remove = f
	return self
}

