package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type CategoryBuilder struct {
	id       uuid.UUID
	parentId *uuid.UUID
	name     string
}

func NewCategory() *CategoryBuilder {
	return &CategoryBuilder{}
}

func (self *CategoryBuilder) WithId(id uuid.UUID) *CategoryBuilder {
	self.id = id
	return self
}

func (self *CategoryBuilder) WithParentId(parentId *uuid.UUID) *CategoryBuilder {
	if nil == parentId {
		self.parentId = nil
	} else {
		self.parentId = new(uuid.UUID)
		*self.parentId = *parentId
	}
	return self
}

func (self *CategoryBuilder) WithName(name string) *CategoryBuilder {
	self.name = name
	return self
}

func (self *CategoryBuilder) Build() models.Category {
	return models.Category{
		Id:       self.id,
		ParentId: self.parentId,
		Name:     self.name,
	}
}

