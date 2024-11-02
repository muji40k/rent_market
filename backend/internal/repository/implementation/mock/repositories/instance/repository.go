package instance

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/implementation/mock/cmnerrors"
	"rent_service/internal/repository/interfaces/instance"

	"github.com/google/uuid"
)

type MockRepository struct {
	create        func(instance models.Instance) (models.Instance, error)
	update        func(instance models.Instance) error
	getById       func(instanceId uuid.UUID) (models.Instance, error)
	getWithFilter func(filter instance.Filter, sort instance.Sort) (Collection[models.Instance], error)
}

func New() *MockRepository {
	return &MockRepository{
		func(instance models.Instance) (models.Instance, error) {
			return models.Instance{}, cmnerrors.ErrorNotSet
		},
		func(instance models.Instance) error {
			return cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (models.Instance, error) {
			return models.Instance{}, cmnerrors.ErrorNotSet
		},
		func(filter instance.Filter, sort instance.Sort) (Collection[models.Instance], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockRepository) Create(instance models.Instance) (models.Instance, error) {
	return self.create(instance)
}

func (self *MockRepository) WithCreate(f func(instance models.Instance) (models.Instance, error)) *MockRepository {
	self.create = f
	return self
}

func (self *MockRepository) Update(instance models.Instance) error {
	return self.update(instance)
}

func (self *MockRepository) WithUpdate(f func(instance models.Instance) error) *MockRepository {
	self.update = f
	return self
}

func (self *MockRepository) GetById(instanceId uuid.UUID) (models.Instance, error) {
	return self.getById(instanceId)
}

func (self *MockRepository) WithGetById(f func(instanceId uuid.UUID) (models.Instance, error)) *MockRepository {
	self.getById = f
	return self
}

func (self *MockRepository) GetWithFilter(filter instance.Filter, sort instance.Sort) (Collection[models.Instance], error) {
	return self.getWithFilter(filter, sort)
}

func (self *MockRepository) WithGetWithFilter(f func(filter instance.Filter, sort instance.Sort) (Collection[models.Instance], error)) *MockRepository {
	self.getWithFilter = f
	return self
}

type MockPayPlansRepository struct {
	create          func(payPlans models.InstancePayPlans) (models.InstancePayPlans, error)
	addPayPlan      func(instanceId uuid.UUID, plan models.PayPlan) (models.InstancePayPlans, error)
	update          func(plans models.InstancePayPlans) error
	getByInstanceId func(instanceId uuid.UUID) (models.InstancePayPlans, error)
}

func NewPayPlansRepository() *MockPayPlansRepository {
	return &MockPayPlansRepository{
		func(payPlans models.InstancePayPlans) (models.InstancePayPlans, error) {
			return models.InstancePayPlans{}, cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID, plan models.PayPlan) (models.InstancePayPlans, error) {
			return models.InstancePayPlans{}, cmnerrors.ErrorNotSet
		},
		func(plans models.InstancePayPlans) error {
			return cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (models.InstancePayPlans, error) {
			return models.InstancePayPlans{}, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockPayPlansRepository) Create(payPlans models.InstancePayPlans) (models.InstancePayPlans, error) {
	return self.create(payPlans)
}

func (self *MockPayPlansRepository) WithCreate(f func(payPlans models.InstancePayPlans) (models.InstancePayPlans, error)) *MockPayPlansRepository {
	self.create = f
	return self
}

func (self *MockPayPlansRepository) AddPayPlan(instanceId uuid.UUID, plan models.PayPlan) (models.InstancePayPlans, error) {
	return self.addPayPlan(instanceId, plan)
}

func (self *MockPayPlansRepository) WithAddPayPlan(f func(instanceId uuid.UUID, plan models.PayPlan) (models.InstancePayPlans, error)) *MockPayPlansRepository {
	self.addPayPlan = f
	return self
}

func (self *MockPayPlansRepository) Update(plans models.InstancePayPlans) error {
	return self.update(plans)
}

func (self *MockPayPlansRepository) WithUpdate(f func(plans models.InstancePayPlans) error) *MockPayPlansRepository {
	self.update = f
	return self
}

func (self *MockPayPlansRepository) GetByInstanceId(instanceId uuid.UUID) (models.InstancePayPlans, error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockPayPlansRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (models.InstancePayPlans, error)) *MockPayPlansRepository {
	self.getByInstanceId = f
	return self
}

type MockPhotoRepository struct {
	create          func(instanceId uuid.UUID, photoId uuid.UUID) error
	getByInstanceId func(instanceId uuid.UUID) (Collection[uuid.UUID], error)
}

func NewPhotoRepository() *MockPhotoRepository {
	return &MockPhotoRepository{
		func(instanceId uuid.UUID, photoId uuid.UUID) error {
			return cmnerrors.ErrorNotSet
		},
		func(instanceId uuid.UUID) (Collection[uuid.UUID], error) {
			return nil, cmnerrors.ErrorNotSet
		},
	}
}

func (self *MockPhotoRepository) Create(instanceId uuid.UUID, photoId uuid.UUID) error {
	return self.create(instanceId, photoId)
}

func (self *MockPhotoRepository) WithCreate(f func(instanceId uuid.UUID, photoId uuid.UUID) error) *MockPhotoRepository {
	self.create = f
	return self
}

func (self *MockPhotoRepository) GetByInstanceId(instanceId uuid.UUID) (Collection[uuid.UUID], error) {
	return self.getByInstanceId(instanceId)
}

func (self *MockPhotoRepository) WithGetByInstanceId(f func(instanceId uuid.UUID) (Collection[uuid.UUID], error)) *MockPhotoRepository {
	self.getByInstanceId = f
	return self
}

