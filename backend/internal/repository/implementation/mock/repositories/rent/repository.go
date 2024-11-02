package rent

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"

	"github.com/google/uuid"
)

type MockRepository struct {
	create                func(rent records.Rent) (records.Rent, error)
	update                func(rent records.Rent) error
	getById               func(rentId uuid.UUID) (records.Rent, error)
	getByUserId           func(userId uuid.UUID) (Collection[records.Rent], error)
	getActiveByInstanceId func(instanceId uuid.UUID) (records.Rent, error)
	getPastByUserId       func(userId uuid.UUID) (Collection[records.Rent], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(rent records.Rent) (records.Rent, error) {
			return records.Rent{}, cmnerrors.ErrorNotSet
		},
		func(rent records.Rent) error {
			return cmnerrors.ErrorNotSet
		},
		func(rentId uuid.UUID) (records.Rent, error) {
			return records.Rent{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[records.Rent], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (records.Rent, error) {
			return records.Rent{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[records.Rent], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(rent records.Rent) (records.Rent, error) {
	return self.create(rent)
}

func (self *MockRepository) WithCreate(f func(rent records.Rent) (records.Rent, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(rent records.Rent) error {
	return self.update(rent)
}

func (self *MockRepository) WithUpdate(f func(rent records.Rent) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetById(rentId uuid.UUID) (records.Rent, error) {
	return self.getById(rentId)
}

func (self *MockRepository) WithGetById(f func(rentId uuid.UUID) (records.Rent, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetByUserId(userId uuid.UUID) (Collection[records.Rent], error) {
	return self.getByUserId(userId)
}

func (self *MockRepository) WithGetByUserId(f func(userId uuid.UUID) (Collection[records.Rent], error)) *MockRepository {
	self.getByUserId = f
	return self
}

func (self *MockRepository) GetActiveByInstanceId(instanceId uuid.UUID) (records.Rent, error) {
	return self.getActiveByInstanceId(instanceId)
}

func (self *MockRepository) WithGetActiveByInstanceId(f func(instanceId uuid.UUID) (records.Rent, error)) *MockRepository {
	self.getActiveByInstanceId = f
	return self
}

func (self *MockRepository) GetPastByUserId(userId uuid.UUID) (Collection[records.Rent], error) {
	return self.getPastByUserId(userId)
}

func (self *MockRepository) WithGetPastByUserId(f func(userId uuid.UUID) (Collection[records.Rent], error)) *MockRepository {
	self.getPastByUserId = f
	return self
}

type MockRequestRepository struct {
	create             func(request requests.Rent) (requests.Rent, error)
	getById            func(requestId uuid.UUID) (requests.Rent, error)
	getByUserId        func(userId uuid.UUID) (Collection[requests.Rent], error)
	getByInstanceId    func(instanceId uuid.UUID) (requests.Rent, error)
	getByPickUpPointId func(pickUpPointId uuid.UUID) (Collection[requests.Rent], error)
	remove             func(requestId uuid.UUID) error
}

func NewRequest() *MockRequestRepository {
	return &MockRequestRepository{
		func(request requests.Rent) (requests.Rent, error) {
			return requests.Rent{}, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) (requests.Rent, error) {
			return requests.Rent{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[requests.Rent], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (requests.Rent, error) {
			return requests.Rent{}, cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (Collection[requests.Rent], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRequestRepository) Create(request requests.Rent) (requests.Rent, error) {
	return self.create(request)
}

func (self *MockRequestRepository) WithCreate(f func(request requests.Rent) (requests.Rent, error)) *MockRequestRepository {
	self.create = f
	return self
}

func (self *MockRequestRepository) GetById(requestId uuid.UUID) (requests.Rent, error) {
	return self.getById(requestId)
}

func (self *MockRequestRepository) WithGetById(f func(requestId uuid.UUID) (requests.Rent, error)) *MockRequestRepository {
	self.getById = f
	return self
}

func (self *MockRequestRepository) GetByUserId(userId uuid.UUID) (Collection[requests.Rent], error) {
	return self.getByUserId(userId)
}

func (self *MockRequestRepository) WithGetByUserId(f func(userId uuid.UUID) (Collection[requests.Rent], error)) *MockRequestRepository {
	self.getByUserId = f
	return self
}

func (self *MockRequestRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Rent, error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockRequestRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (requests.Rent, error)) *MockRequestRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockRequestRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (Collection[requests.Rent], error) {
	return self.getByPickUpPointId(pickUpPointId)
}

func (self *MockRequestRepository) WithGetByPickUpPointId(f func(pickUpPointId uuid.UUID) (Collection[requests.Rent], error)) *MockRequestRepository {
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

type MockReturnRepository struct {
	create             func(request requests.Return) (requests.Return, error)
	getById            func(requestId uuid.UUID) (requests.Return, error)
	getByUserId        func(userId uuid.UUID) (Collection[requests.Return], error)
	getByInstanceId    func(instanceId uuid.UUID) (requests.Return, error)
	getByPickUpPointId func(pickUpPointId uuid.UUID) (Collection[requests.Return], error)
	remove             func(requestId uuid.UUID) error
}

func NewReturn() *MockReturnRepository {
	return &MockReturnRepository{
		func(request requests.Return) (requests.Return, error) {
			return requests.Return{}, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) (requests.Return, error) {
			return requests.Return{}, cmnerrors.ErrorNotSet
		},
		func(userId uuid.UUID) (Collection[requests.Return], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (requests.Return, error) {
			return requests.Return{}, cmnerrors.ErrorNotSet
		},
		func(pickUpPointId uuid.UUID) (Collection[requests.Return], error) {
			return nil, cmnerrors.ErrorNotSet
		},
		func(requestId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockReturnRepository) Create(request requests.Return) (requests.Return, error) {
	return self.create(request)
}

func (self *MockReturnRepository) WithCreate(f func(request requests.Return) (requests.Return, error)) *MockReturnRepository {
	self.create = f
	return self
}

func (self *MockReturnRepository) GetById(requestId uuid.UUID) (requests.Return, error) {
	return self.getById(requestId)
}

func (self *MockReturnRepository) WithGetById(f func(requestId uuid.UUID) (requests.Return, error)) *MockReturnRepository {
	self.getById = f
	return self
}

func (self *MockReturnRepository) GetByUserId(userId uuid.UUID) (Collection[requests.Return], error) {
	return self.getByUserId(userId)
}

func (self *MockReturnRepository) WithGetByUserId(f func(userId uuid.UUID) (Collection[requests.Return], error)) *MockReturnRepository {
	self.getByUserId = f
	return self
}

func (self *MockReturnRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Return, error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockReturnRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (requests.Return, error)) *MockReturnRepository {
	self.getByInstanceId = f
	return self
}

func (self *MockReturnRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (Collection[requests.Return], error) {
	return self.getByPickUpPointId(pickUpPointId)
}

func (self *MockReturnRepository) WithGetByPickUpPointId(f func(pickUpPointId uuid.UUID) (Collection[requests.Return], error)) *MockReturnRepository {
	self.getByPickUpPointId = f
	return self
}

func (self *MockReturnRepository) Remove(requestId uuid.UUID) error {
	return self.remove(requestId)
}

func (self *MockReturnRepository) WithRemove(f func(requestId uuid.UUID) error) *MockReturnRepository {
	self.remove = f
	return self
}

