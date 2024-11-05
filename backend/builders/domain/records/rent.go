package records

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/records"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type RentBuilder struct {
	id              uuid.UUID
	userId          uuid.UUID
	instanceId      uuid.UUID
	startDate       time.Time
	endDate         *time.Time
	paymentPeriodId uuid.UUID
}

func NewRent() *RentBuilder {
	return &RentBuilder{}
}

func (self *RentBuilder) WithId(id uuid.UUID) *RentBuilder {
	self.id = id
	return self
}

func (self *RentBuilder) WithUserId(userId uuid.UUID) *RentBuilder {
	self.userId = userId
	return self
}

func (self *RentBuilder) WithInstanceId(instanceId uuid.UUID) *RentBuilder {
	self.instanceId = instanceId
	return self
}

func (self *RentBuilder) WithStartDate(startDate time.Time) *RentBuilder {
	self.startDate = startDate
	return self
}

func (self *RentBuilder) WithEndDate(endDate *nullable.Nullable[time.Time]) *RentBuilder {
	self.endDate = nullcommon.CopyPtrIfSome(endDate)
	return self
}

func (self *RentBuilder) WithPaymentPeriodId(paymentPeriodId uuid.UUID) *RentBuilder {
	self.paymentPeriodId = paymentPeriodId
	return self
}

func (self *RentBuilder) Build() records.Rent {
	return records.Rent{
		Id:              self.id,
		UserId:          self.userId,
		InstanceId:      self.instanceId,
		StartDate:       self.startDate,
		EndDate:         self.endDate,
		PaymentPeriodId: self.paymentPeriodId,
	}
}

