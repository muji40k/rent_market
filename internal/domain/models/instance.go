package models

import (
	"github.com/google/uuid"
)

type Instance struct {
	Id          uuid.UUID
	ProductId   uuid.UUID
	Name        string
	Description string
	Condition   string
}

type InstancePayPlans struct {
	InstanceId uuid.UUID
	Map        map[uuid.UUID]PayPlan // Indexed by Period uuid
}

func NewInstancePayPlans() InstancePayPlans {
	out := InstancePayPlans{}
	out.Map = make(map[uuid.UUID]PayPlan)

	return out
}

