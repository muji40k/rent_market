package provision

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create                func(provision records.Provision) (records.Provision, error)
	update                func(provision records.Provision) error
	getById               func(provisionId uuid.UUID) (records.Provision, error)
	getByRenterUserId     func(userId uuid.UUID) (Collection[records.Provision], error)
	getByInstanceId       func(instanceId uuid.UUID) (Collection[records.Provision], error)
	getActiveByInstanceId func(instanceId uuid.UUID) (records.Provision, error)
}

func New() *MockRepository {
	return &MockRepository{
		func(provision records.Provision) (records.Provision, error) {
			return records.Provision{}, cmnerrors.ErrorNotSet
		},
		func(provision records.Provision) error {
			return cmnerrors.ErrorNotSet
		},
		func(provisionId uuid.UUID) (records.Provision, error) {
			return records.Provision{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[records.Provision], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (Collection[records.Provision], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (records.Provision, error) {
			return records.Provision{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(provision records.Provision) (records.Provision, error) {
	return self.create(provision)
}

func (self *MockRepository) WithCreate(f func(provision records.Provision) (records.Provision, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(provision records.Provision) error {
	return self.update(provision)
}

func (self *MockRepository) WithUpdate(f func(provision records.Provision) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetById(provisionId uuid.UUID) (records.Provision, error) {
	return self.getById(provisionId)
}

func (self *MockRepository) WithGetById(f func(provisionId uuid.UUID) (records.Provision, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetByRenterUserId(userId uuid.UUID) (Collection[records.Provision], error) {
	return self.getByRenterUserId(userId)
}

func (self *MockRepository) WithGetByRenterUserId(f func(userId uuid.UUID) (Collection[records.Provision], error)) *MockRepository {
	self.getByRenterUserId = f
	return self
}

func (self *MockRepository) GetByInstanceId(instanceId uuid.UUID) (Collection[records.Provision], error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (Collection[records.Provision], error)) *MockRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockRepository) GetActiveByInstanceId(instanceId uuid.UUID) (records.Provision, error) {
	return self.getActiveByInstanceId(instanceId)
}

func (self *MockRepository) WithGetActiveByInstanceId(f func(instanceId uuid.UUID) (records.Provision, error)) *MockRepository {
	self.getActiveByInstanceId = f
	return self
}

type MockRequestRepository struct {
	create             func(request requests.Provide) (requests.Provide, error)
	getById            func(requestId uuid.UUID) (requests.Provide, error)
	getByUserId        func(userId uuid.UUID) (Collection[requests.Provide], error)
	getByInstanceId    func(instanceId uuid.UUID) (requests.Provide, error)
	getByPickUpPointId func(pickUpPointId uuid.UUID) (Collection[requests.Provide], error)
	remove             func(requestId uuid.UUID) error
}

func NewRequest() *MockRequestRepository {
	return &MockRequestRepository{
		func(request requests.Provide) (requests.Provide, error) {
			return requests.Provide{}, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) (requests.Provide, error) {
			return requests.Provide{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[requests.Provide], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (requests.Provide, error) {
			return requests.Provide{}, cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (Collection[requests.Provide], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRequestRepository) Create(request requests.Provide) (requests.Provide, error) {
	return self.create(request)
}

func (self *MockRequestRepository) WithCreate(f func(request requests.Provide) (requests.Provide, error)) *MockRequestRepository {
	self.create = f
	return self
}

func (self *MockRequestRepository) GetById(requestId uuid.UUID) (requests.Provide, error) {
	return self.getById(requestId)
}

func (self *MockRequestRepository) WithGetById(f func(requestId uuid.UUID) (requests.Provide, error)) *MockRequestRepository {
	self.getById = f
	return self
}

func (self *MockRequestRepository) GetByUserId(userId uuid.UUID) (Collection[requests.Provide], error) {
	return self.getByUserId(userId)
}

func (self *MockRequestRepository) WithGetByUserId(f func(userId uuid.UUID) (Collection[requests.Provide], error)) *MockRequestRepository {
	self.getByUserId = f
	return self
}

func (self *MockRequestRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Provide, error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockRequestRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (requests.Provide, error)) *MockRequestRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockRequestRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (Collection[requests.Provide], error) {
	return self.getByPickUpPointId(pickUpPointId)
}

func (self *MockRequestRepository) WithGetByPickUpPointId(f func(pickUpPointId uuid.UUID) (Collection[requests.Provide], error)) *MockRequestRepository {
	self.getByPickUpPointId = f
	return self
}

func (self *MockRequestRepository) Remove(requestId uuid.UUID) error {
	return self.remove(requestId)
}

func (self *MockRequestRepository) WithRemove(f func(requestId uuid.UUID) error) *MockRequestRepository {
	self.remove = f
	return self
}

type MockRevokeRepository struct {
	create             func(request requests.Revoke) (requests.Revoke, error)
	getById            func(requestId uuid.UUID) (requests.Revoke, error)
	getByUserId        func(userId uuid.UUID) (Collection[requests.Revoke], error)
	getByInstanceId    func(instanceId uuid.UUID) (requests.Revoke, error)
	getByPickUpPointId func(pickUpPointId uuid.UUID) (Collection[requests.Revoke], error)
	remove             func(requestId uuid.UUID) error
}

func NewRevoke() *MockRevokeRepository {
	return &MockRevokeRepository{
		func(request requests.Revoke) (requests.Revoke, error) {
			return requests.Revoke{}, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) (requests.Revoke, error) {
			return requests.Revoke{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[requests.Revoke], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (requests.Revoke, error) {
			return requests.Revoke{}, cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (Collection[requests.Revoke], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRevokeRepository) Create(request requests.Revoke) (requests.Revoke, error) {
	return self.create(request)
}

func (self *MockRevokeRepository) WithCreate(f func(request requests.Revoke) (requests.Revoke, error)) *MockRevokeRepository {
	self.create = f
	return self
}

func (self *MockRevokeRepository) GetById(requestId uuid.UUID) (requests.Revoke, error) {
	return self.getById(requestId)
}

func (self *MockRevokeRepository) WithGetById(f func(requestId uuid.UUID) (requests.Revoke, error)) *MockRevokeRepository {
	self.getById = f
	return self
}

func (self *MockRevokeRepository) GetByUserId(userId uuid.UUID) (Collection[requests.Revoke], error) {
	return self.getByUserId(userId)
}

func (self *MockRevokeRepository) WithGetByUserId(f func(userId uuid.UUID) (Collection[requests.Revoke], error)) *MockRevokeRepository {
	self.getByUserId = f
	return self
}

func (self *MockRevokeRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Revoke, error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockRevokeRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (requests.Revoke, error)) *MockRevokeRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockRevokeRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (Collection[requests.Revoke], error) {
	return self.getByPickUpPointId(pickUpPointId)
}

func (self *MockRevokeRepository) WithGetByPickUpPointId(f func(pickUpPointId uuid.UUID) (Collection[requests.Revoke], error)) *MockRevokeRepository {
	self.getByPickUpPointId = f
	return self
}

func (self *MockRevokeRepository) Remove(requestId uuid.UUID) error {
	return self.remove(requestId)
}

func (self *MockRevokeRepository) WithRemove(f func(requestId uuid.UUID) error) *MockRevokeRepository {
	self.remove = f
	return self
}

