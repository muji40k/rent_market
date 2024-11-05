package models

import (
	"rent_service/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type ReviewBuilder struct {
	id         uuid.UUID
	instanceId uuid.UUID
	userId     uuid.UUID
	content    string
	rating     float64
	date       time.Time
}

func NewReview() *ReviewBuilder {
	return &ReviewBuilder{}
}

func (self *ReviewBuilder) WithId(id uuid.UUID) *ReviewBuilder {
	self.id = id
	return self
}

func (self *ReviewBuilder) WithInstanceId(instanceId uuid.UUID) *ReviewBuilder {
	self.instanceId = instanceId
	return self
}

func (self *ReviewBuilder) WithUserId(userId uuid.UUID) *ReviewBuilder {
	self.userId = userId
	return self
}

func (self *ReviewBuilder) WithContent(content string) *ReviewBuilder {
	self.content = content
	return self
}

func (self *ReviewBuilder) WithRating(rating float64) *ReviewBuilder {
	self.rating = rating
	return self
}

func (self *ReviewBuilder) WithDate(date time.Time) *ReviewBuilder {
	self.date = date
	return self
}

func (self *ReviewBuilder) Build() models.Review {
	return models.Review{
		Id:         self.id,
		InstanceId: self.instanceId,
		UserId:     self.userId,
		Content:    self.content,
		Rating:     self.rating,
		Date:       self.date,
	}
}

