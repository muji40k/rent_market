package storage

import (
	"rent_service/internal/domain/records"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create                   func(storage records.Storage) (records.Storage, error)
	update                   func(storage records.Storage) error
	getActiveByPickUpPointId func(pickUpPointId uuid.UUID) (Collection[records.Storage], error)
	getActiveByInstanceId    func(instanceId uuid.UUID) (records.Storage, error)
}

func New() *MockRepository {
	return &MockRepository{
		func(storage records.Storage) (records.Storage, error) {
			return records.Storage{}, cmnerrors.ErrorNotSet
		},
		func(storage records.Storage) error {
			return cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (Collection[records.Storage], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (records.Storage, error) {
			return records.Storage{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(storage records.Storage) (records.Storage, error) {
	return self.create(storage)
}

func (self *MockRepository) WithCreate(f func(storage records.Storage) (records.Storage, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(storage records.Storage) error {
	return self.update(storage)
}

func (self *MockRepository) WithUpdate(f func(storage records.Storage) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetActiveByPickUpPointId(pickUpPointId uuid.UUID) (Collection[records.Storage], error) {
	return self.getActiveByPickUpPointId(pickUpPointId)
}

func (self *MockRepository) WithGetActiveByPickUpPointId(f func(pickUpPointId uuid.UUID) (Collection[records.Storage], error)) *MockRepository {
	self.getActiveByPickUpPointId = f
	return self
}

func (self *MockRepository) GetActiveByInstanceId(instanceId uuid.UUID) (records.Storage, error) {
	return self.getActiveByInstanceId(instanceId)
}

func (self *MockRepository) WithGetActiveByInstanceId(f func(instanceId uuid.UUID) (records.Storage, error)) *MockRepository {
	self.getActiveByInstanceId = f
	return self
}

