package requests

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/requests"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type DeliveryBuilder struct {
	id                 uuid.UUID
	companyId          uuid.UUID
	instanceId         uuid.UUID
	fromId             uuid.UUID
	toId               uuid.UUID
	deliveryId         string
	scheduledBeginDate time.Time
	actualBeginDate    *time.Time
	scheduledEndDate   time.Time
	actualEndDate      *time.Time
	verificationCode   string
	createDate         time.Time
}

func NewDelivery() *DeliveryBuilder {
	return &DeliveryBuilder{}
}

func (self *DeliveryBuilder) WithId(id uuid.UUID) *DeliveryBuilder {
	self.id = id
	return self
}

func (self *DeliveryBuilder) WithCompanyId(companyId uuid.UUID) *DeliveryBuilder {
	self.companyId = companyId
	return self
}

func (self *DeliveryBuilder) WithInstanceId(instanceId uuid.UUID) *DeliveryBuilder {
	self.instanceId = instanceId
	return self
}

func (self *DeliveryBuilder) WithFromId(fromId uuid.UUID) *DeliveryBuilder {
	self.fromId = fromId
	return self
}

func (self *DeliveryBuilder) WithToId(toId uuid.UUID) *DeliveryBuilder {
	self.toId = toId
	return self
}

func (self *DeliveryBuilder) WithDeliveryId(deliveryId string) *DeliveryBuilder {
	self.deliveryId = deliveryId
	return self
}

func (self *DeliveryBuilder) WithScheduledBeginDate(scheduledBeginDate time.Time) *DeliveryBuilder {
	self.scheduledBeginDate = scheduledBeginDate
	return self
}

func (self *DeliveryBuilder) WithActualBeginDate(actualBeginDate *nullable.Nullable[time.Time]) *DeliveryBuilder {
	self.actualBeginDate = nullcommon.CopyPtrIfSome(actualBeginDate)
	return self
}

func (self *DeliveryBuilder) WithScheduledEndDate(scheduledEndDate time.Time) *DeliveryBuilder {
	self.scheduledEndDate = scheduledEndDate
	return self
}

func (self *DeliveryBuilder) WithActualEndDate(actualEndDate *nullable.Nullable[time.Time]) *DeliveryBuilder {
	self.actualEndDate = nullcommon.CopyPtrIfSome(actualEndDate)
	return self
}

func (self *DeliveryBuilder) WithVerificationCode(verificationCode string) *DeliveryBuilder {
	self.verificationCode = verificationCode
	return self
}

func (self *DeliveryBuilder) WithCreateDate(createDate time.Time) *DeliveryBuilder {
	self.createDate = createDate
	return self
}

func (self *DeliveryBuilder) Build() requests.Delivery {
	return requests.Delivery{
		Id:                 self.id,
		CompanyId:          self.companyId,
		InstanceId:         self.instanceId,
		FromId:             self.fromId,
		ToId:               self.toId,
		DeliveryId:         self.deliveryId,
		ScheduledBeginDate: self.scheduledBeginDate,
		ActualBeginDate:    self.actualBeginDate,
		ScheduledEndDate:   self.scheduledEndDate,
		ActualEndDate:      self.actualEndDate,
		VerificationCode:   self.verificationCode,
		CreateDate:         self.createDate,
	}
}

