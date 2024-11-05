package requests

import (
	"rent_service/internal/domain/requests"
	"time"

	"github.com/google/uuid"
)

type RentBuilder struct {
	id               uuid.UUID
	instanceId       uuid.UUID
	userId           uuid.UUID
	pickUpPointId    uuid.UUID
	paymentPeriodId  uuid.UUID
	verificationCode string
	createDate       time.Time
}

func NewRent() *RentBuilder {
	return &RentBuilder{}
}

func (self *RentBuilder) WithId(id uuid.UUID) *RentBuilder {
	self.id = id
	return self
}

func (self *RentBuilder) WithInstanceId(instanceId uuid.UUID) *RentBuilder {
	self.instanceId = instanceId
	return self
}

func (self *RentBuilder) WithUserId(userId uuid.UUID) *RentBuilder {
	self.userId = userId
	return self
}

func (self *RentBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *RentBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *RentBuilder) WithPaymentPeriodId(paymentPeriodId uuid.UUID) *RentBuilder {
	self.paymentPeriodId = paymentPeriodId
	return self
}

func (self *RentBuilder) WithVerificationCode(verificationCode string) *RentBuilder {
	self.verificationCode = verificationCode
	return self
}

func (self *RentBuilder) WithCreateDate(createDate time.Time) *RentBuilder {
	self.createDate = createDate
	return self
}

func (self *RentBuilder) Build() requests.Rent {
	return requests.Rent{
		Id:               self.id,
		InstanceId:       self.instanceId,
		UserId:           self.userId,
		PickUpPointId:    self.pickUpPointId,
		PaymentPeriodId:  self.paymentPeriodId,
		VerificationCode: self.verificationCode,
		CreateDate:       self.createDate,
	}
}

