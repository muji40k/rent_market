package records

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/records"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type ProvisionBuilder struct {
	id         uuid.UUID
	renterId   uuid.UUID
	instanceId uuid.UUID
	startDate  time.Time
	endDate    *time.Time
}

func NewProvision() *ProvisionBuilder {
	return &ProvisionBuilder{}
}

func (self *ProvisionBuilder) WithId(id uuid.UUID) *ProvisionBuilder {
	self.id = id
	return self
}

func (self *ProvisionBuilder) WithRenterId(renterId uuid.UUID) *ProvisionBuilder {
	self.renterId = renterId
	return self
}

func (self *ProvisionBuilder) WithInstanceId(instanceId uuid.UUID) *ProvisionBuilder {
	self.instanceId = instanceId
	return self
}

func (self *ProvisionBuilder) WithStartDate(startDate time.Time) *ProvisionBuilder {
	self.startDate = startDate
	return self
}

func (self *ProvisionBuilder) WithEndDate(endDate *nullable.Nullable[time.Time]) *ProvisionBuilder {
	self.endDate = nullcommon.CopyPtrIfSome(endDate)
	return self
}

func (self *ProvisionBuilder) Build() records.Provision {
	return records.Provision{
		Id:         self.id,
		RenterId:   self.renterId,
		InstanceId: self.instanceId,
		StartDate:  self.startDate,
		EndDate:    self.endDate,
	}
}

