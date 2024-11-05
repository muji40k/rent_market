package models

import (
	"rent_service/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type PeriodBuilder struct {
	id       uuid.UUID
	name     string
	duration time.Duration
}

func NewPeriod() *PeriodBuilder {
	return &PeriodBuilder{}
}

func (self *PeriodBuilder) WithId(id uuid.UUID) *PeriodBuilder {
	self.id = id
	return self
}

func (self *PeriodBuilder) WithName(name string) *PeriodBuilder {
	self.name = name
	return self
}

func (self *PeriodBuilder) WithDuration(duration time.Duration) *PeriodBuilder {
	self.duration = duration
	return self
}

func (self *PeriodBuilder) Build() models.Period {
	return models.Period{
		Id:       self.id,
		Name:     self.name,
		Duration: self.duration,
	}
}

