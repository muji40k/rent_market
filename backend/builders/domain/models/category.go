package models

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"

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

func (self *CategoryBuilder) WithParentId(parentId *nullable.Nullable[uuid.UUID]) *CategoryBuilder {
	self.parentId = nullcommon.CopyPtrIfSome(parentId)
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

