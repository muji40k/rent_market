package models

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type PayPlanBuilder struct {
	id       uuid.UUID
	periodId uuid.UUID
	price    currency.Currency
}

func NewPayPlan() *PayPlanBuilder {
	return &PayPlanBuilder{}
}

func (self *PayPlanBuilder) WithId(id uuid.UUID) *PayPlanBuilder {
	self.id = id
	return self
}

func (self *PayPlanBuilder) WithPeriodId(periodId uuid.UUID) *PayPlanBuilder {
	self.periodId = periodId
	return self
}

func (self *PayPlanBuilder) WithPrice(price currency.Currency) *PayPlanBuilder {
	self.price = price
	return self
}

func (self *PayPlanBuilder) Build() models.PayPlan {
	return models.PayPlan{
		Id:       self.id,
		PeriodId: self.periodId,
		Price:    self.price,
	}
}

