package requests

import (
	"rent_service/internal/domain/requests"
	"time"

	"github.com/google/uuid"
)

type ReturnBuilder struct {
	id               uuid.UUID
	instanceId       uuid.UUID
	userId           uuid.UUID
	pickUpPointId    uuid.UUID
	rentEndDate      time.Time
	verificationCode string
	createDate       time.Time
}

func NewReturn() *ReturnBuilder {
	return &ReturnBuilder{}
}

func (self *ReturnBuilder) WithId(id uuid.UUID) *ReturnBuilder {
	self.id = id
	return self
}

func (self *ReturnBuilder) WithInstanceId(instanceId uuid.UUID) *ReturnBuilder {
	self.instanceId = instanceId
	return self
}

func (self *ReturnBuilder) WithUserId(userId uuid.UUID) *ReturnBuilder {
	self.userId = userId
	return self
}

func (self *ReturnBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *ReturnBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *ReturnBuilder) WithRentEndDate(rentEndDate time.Time) *ReturnBuilder {
	self.rentEndDate = rentEndDate
	return self
}

func (self *ReturnBuilder) WithVerificationCode(verificationCode string) *ReturnBuilder {
	self.verificationCode = verificationCode
	return self
}

func (self *ReturnBuilder) WithCreateDate(createDate time.Time) *ReturnBuilder {
	self.createDate = createDate
	return self
}

func (self *ReturnBuilder) Build() requests.Return {
	return requests.Return{
		Id:               self.id,
		InstanceId:       self.instanceId,
		UserId:           self.userId,
		PickUpPointId:    self.pickUpPointId,
		RentEndDate:      self.rentEndDate,
		VerificationCode: self.verificationCode,
		CreateDate:       self.createDate,
	}
}

