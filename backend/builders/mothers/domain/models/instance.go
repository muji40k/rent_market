package models

import (
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/collect"
	"rent_service/builders/misc/uuidgen"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

func InstanceRandomId() *modelsb.InstanceBuilder {
	return modelsb.NewInstance().
		WithId(uuidgen.Generate())
}

func InstanceExample(prefix string, productId uuid.UUID) *modelsb.InstanceBuilder {
	return InstanceRandomId().
		WithProductId(productId).
		WithName("Example " + prefix).
		WithDescription("Example Instance for tests").
		WithCondition("As new")
}

func InstancePayPlansExample(instanceId uuid.UUID, periods ...models.Period) *modelsb.InstancePayPlansBuilder {
	return modelsb.NewInstancePayPlans().
		WithInstanceId(instanceId).
		WithPayPlans(collect.Do(PayPlansWithPeriods(periods...)...)...)
}

