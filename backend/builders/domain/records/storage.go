package records

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/records"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type StorageBuilder struct {
	id            uuid.UUID
	pickUpPointId uuid.UUID
	instanceId    uuid.UUID
	inDate        time.Time
	outDate       *time.Time
}

func NewStorage() *StorageBuilder {
	return &StorageBuilder{}
}

func (self *StorageBuilder) WithId(id uuid.UUID) *StorageBuilder {
	self.id = id
	return self
}

func (self *StorageBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *StorageBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *StorageBuilder) WithInstanceId(instanceId uuid.UUID) *StorageBuilder {
	self.instanceId = instanceId
	return self
}

func (self *StorageBuilder) WithInDate(inDate time.Time) *StorageBuilder {
	self.inDate = inDate
	return self
}

func (self *StorageBuilder) WithOutDate(outDate *nullable.Nullable[time.Time]) *StorageBuilder {
	self.outDate = nullcommon.CopyPtrIfSome(outDate)
	return self
}

func (self *StorageBuilder) Build() records.Storage {
	return records.Storage{
		Id:            self.id,
		PickUpPointId: self.pickUpPointId,
		InstanceId:    self.instanceId,
		InDate:        self.inDate,
		OutDate:       self.outDate,
	}
}

