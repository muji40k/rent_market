package requests

import (
	"rent_service/internal/domain/requests"
	"time"

	"github.com/google/uuid"
)

type RevokeBuilder struct {
	id               uuid.UUID
	instanceId       uuid.UUID
	renterId         uuid.UUID
	pickUpPointId    uuid.UUID
	verificationCode string
	createDate       time.Time
}

func NewRevoke() *RevokeBuilder {
	return &RevokeBuilder{}
}

func (self *RevokeBuilder) WithId(id uuid.UUID) *RevokeBuilder {
	self.id = id
	return self
}

func (self *RevokeBuilder) WithInstanceId(instanceId uuid.UUID) *RevokeBuilder {
	self.instanceId = instanceId
	return self
}

func (self *RevokeBuilder) WithRenterId(renterId uuid.UUID) *RevokeBuilder {
	self.renterId = renterId
	return self
}

func (self *RevokeBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *RevokeBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *RevokeBuilder) WithVerificationCode(verificationCode string) *RevokeBuilder {
	self.verificationCode = verificationCode
	return self
}

func (self *RevokeBuilder) WithCreateDate(createDate time.Time) *RevokeBuilder {
	self.createDate = createDate
	return self
}

func (self *RevokeBuilder) Build() requests.Revoke {
	return requests.Revoke{
		Id:               self.id,
		InstanceId:       self.instanceId,
		RenterId:         self.renterId,
		PickUpPointId:    self.pickUpPointId,
		VerificationCode: self.verificationCode,
		CreateDate:       self.createDate,
	}
}

