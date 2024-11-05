package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type RenterBuilder struct {
	id     uuid.UUID
	userId uuid.UUID
}

func NewRenter() *RenterBuilder {
	return &RenterBuilder{}
}

func (self *RenterBuilder) WithId(id uuid.UUID) *RenterBuilder {
	self.id = id
	return self
}

func (self *RenterBuilder) WithUserId(userId uuid.UUID) *RenterBuilder {
	self.userId = userId
	return self
}

func (self *RenterBuilder) Build() models.Renter {
	return models.Renter{
		Id:     self.id,
		UserId: self.userId,
	}
}

