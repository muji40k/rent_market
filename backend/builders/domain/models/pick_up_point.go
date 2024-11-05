package models

import (
	"rent_service/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type PickUpPointBuilder struct {
	id       uuid.UUID
	address  models.Address
	capacity uint64
}

func NewPickUpPoint() *PickUpPointBuilder {
	return &PickUpPointBuilder{}
}

func (self *PickUpPointBuilder) WithId(id uuid.UUID) *PickUpPointBuilder {
	self.id = id
	return self
}

func (self *PickUpPointBuilder) WithAddress(address models.Address) *PickUpPointBuilder {
	self.address = address
	return self
}

func (self *PickUpPointBuilder) WithCapacity(capacity uint64) *PickUpPointBuilder {
	self.capacity = capacity
	return self
}

func (self *PickUpPointBuilder) Build() models.PickUpPoint {
	return models.PickUpPoint{
		Id:       self.id,
		Address:  self.address,
		Capacity: self.capacity,
	}
}

type WorkingHoursBuilder struct {
	id    uuid.UUID
	day   time.Weekday
	begin time.Duration
	end   time.Duration
}

func NewWorkingHours() *WorkingHoursBuilder {
	return &WorkingHoursBuilder{}
}

func (self *WorkingHoursBuilder) WithId(id uuid.UUID) *WorkingHoursBuilder {
	self.id = id
	return self
}

func (self *WorkingHoursBuilder) WithDay(day time.Weekday) *WorkingHoursBuilder {
	self.day = day
	return self
}

func (self *WorkingHoursBuilder) WithBegin(begin time.Duration) *WorkingHoursBuilder {
	self.begin = begin
	return self
}

func (self *WorkingHoursBuilder) WithEnd(end time.Duration) *WorkingHoursBuilder {
	self.end = end
	return self
}

func (self *WorkingHoursBuilder) Build() models.WorkingHours {
	return models.WorkingHours{
		Id:    self.id,
		Day:   self.day,
		Begin: self.begin,
		End:   self.end,
	}
}

type PickUpPointWorkingHoursBuilder struct {
	pickUpPointId uuid.UUID
	mmap          map[time.Weekday]models.WorkingHours
}

func NewPickUpPointWorkingHours() *PickUpPointWorkingHoursBuilder {
	return &PickUpPointWorkingHoursBuilder{}
}

func (self *PickUpPointWorkingHoursBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *PickUpPointWorkingHoursBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *PickUpPointWorkingHoursBuilder) WithWorkingHours(wh ...models.WorkingHours) *PickUpPointWorkingHoursBuilder {
	self.mmap = make(map[time.Weekday]models.WorkingHours, len(wh))

	for _, value := range wh {
		self.mmap[value.Day] = value
	}

	return self
}

func (self *PickUpPointWorkingHoursBuilder) Build() models.PickUpPointWorkingHours {
	return models.PickUpPointWorkingHours{
		PickUpPointId: self.pickUpPointId,
		Map:           self.mmap,
	}
}

