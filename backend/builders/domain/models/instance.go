package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type InstanceBuilder struct {
	id          uuid.UUID
	productId   uuid.UUID
	name        string
	description string
	condition   string
}

func NewInstance() *InstanceBuilder {
	return &InstanceBuilder{}
}

func (self *InstanceBuilder) WithId(id uuid.UUID) *InstanceBuilder {
	self.id = id
	return self
}

func (self *InstanceBuilder) WithProductId(productId uuid.UUID) *InstanceBuilder {
	self.productId = productId
	return self
}

func (self *InstanceBuilder) WithName(name string) *InstanceBuilder {
	self.name = name
	return self
}

func (self *InstanceBuilder) WithDescription(description string) *InstanceBuilder {
	self.description = description
	return self
}

func (self *InstanceBuilder) WithCondition(condition string) *InstanceBuilder {
	self.condition = condition
	return self
}

func (self *InstanceBuilder) Build() models.Instance {
	return models.Instance{
		Id:          self.id,
		ProductId:   self.productId,
		Name:        self.name,
		Description: self.description,
		Condition:   self.condition,
	}
}

type InstancePayPlansBuilder struct {
	instanceId uuid.UUID
	mmap       map[uuid.UUID]models.PayPlan
}

func NewInstancePayPlans() *InstancePayPlansBuilder {
	return &InstancePayPlansBuilder{}
}

func (self *InstancePayPlansBuilder) WithInstanceId(instanceId uuid.UUID) *InstancePayPlansBuilder {
	self.instanceId = instanceId
	return self
}

func (self *InstancePayPlansBuilder) WithPayPlans(plans ...models.PayPlan) *InstancePayPlansBuilder {
	self.mmap = make(map[uuid.UUID]models.PayPlan, len(plans))

	for _, p := range plans {
		self.mmap[p.PeriodId] = p
	}

	return self
}

func (self *InstancePayPlansBuilder) Build() models.InstancePayPlans {
	return models.InstancePayPlans{
		InstanceId: self.instanceId,
		Map:        self.mmap,
	}
}

