package models

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/currency"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type PaymentBuilder struct {
	id          uuid.UUID
	rentId      uuid.UUID
	payMethodId *uuid.UUID
	paymentId   *string
	periodStart time.Time
	periodEnd   time.Time
	value       currency.Currency
	status      string
	createDate  time.Time
	paymentDate *time.Time
}

func NewPayment() *PaymentBuilder {
	return &PaymentBuilder{}
}

func (self *PaymentBuilder) WithId(id uuid.UUID) *PaymentBuilder {
	self.id = id
	return self
}

func (self *PaymentBuilder) WithRentId(rentId uuid.UUID) *PaymentBuilder {
	self.rentId = rentId
	return self
}

func (self *PaymentBuilder) WithPayMethodId(payMethodId *nullable.Nullable[uuid.UUID]) *PaymentBuilder {
	self.payMethodId = nullcommon.CopyPtrIfSome(payMethodId)
	return self
}

func (self *PaymentBuilder) WithPaymentId(paymentId *nullable.Nullable[string]) *PaymentBuilder {
	self.paymentId = nullcommon.CopyPtrIfSome(paymentId)
	return self
}

func (self *PaymentBuilder) WithPeriodStart(periodStart time.Time) *PaymentBuilder {
	self.periodStart = periodStart
	return self
}

func (self *PaymentBuilder) WithPeriodEnd(periodEnd time.Time) *PaymentBuilder {
	self.periodEnd = periodEnd
	return self
}

func (self *PaymentBuilder) WithValue(value currency.Currency) *PaymentBuilder {
	self.value = value
	return self
}

func (self *PaymentBuilder) WithStatus(status string) *PaymentBuilder {
	self.status = status
	return self
}

func (self *PaymentBuilder) WithCreateDate(createDate time.Time) *PaymentBuilder {
	self.createDate = createDate
	return self
}

func (self *PaymentBuilder) WithPaymentDate(paymentDate *nullable.Nullable[time.Time]) *PaymentBuilder {
	self.paymentDate = nullcommon.CopyPtrIfSome(paymentDate)
	return self
}

func (self *PaymentBuilder) Build() models.Payment {
	return models.Payment{
		Id:          self.id,
		RentId:      self.rentId,
		PayMethodId: self.payMethodId,
		PaymentId:   self.paymentId,
		PeriodStart: self.periodStart,
		PeriodEnd:   self.periodEnd,
		Value:       self.value,
		Status:      self.status,
		CreateDate:  self.createDate,
		PaymentDate: self.paymentDate,
	}
}

