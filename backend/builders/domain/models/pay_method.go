package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type PayMethodBuilder struct {
	id          uuid.UUID
	name        string
	description string
}

func NewPayMethod() *PayMethodBuilder {
	return &PayMethodBuilder{}
}

func (self *PayMethodBuilder) WithId(id uuid.UUID) *PayMethodBuilder {
	self.id = id
	return self
}

func (self *PayMethodBuilder) WithName(name string) *PayMethodBuilder {
	self.name = name
	return self
}

func (self *PayMethodBuilder) WithDescription(description string) *PayMethodBuilder {
	self.description = description
	return self
}

func (self *PayMethodBuilder) Build() models.PayMethod {
	return models.PayMethod{
		Id:          self.id,
		Name:        self.name,
		Description: self.description,
	}
}

