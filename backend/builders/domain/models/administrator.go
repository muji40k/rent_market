package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type AdministratorBuilder struct {
	id     uuid.UUID
	userId uuid.UUID
}

func NewAdministrator() *AdministratorBuilder {
	return &AdministratorBuilder{}
}

func (self *AdministratorBuilder) WithId(id uuid.UUID) *AdministratorBuilder {
	self.id = id
	return self
}

func (self *AdministratorBuilder) WithUserId(userId uuid.UUID) *AdministratorBuilder {
	self.userId = userId
	return self
}

func (self *AdministratorBuilder) Build() models.Administrator {
	return models.Administrator{
		Id:     self.id,
		UserId: self.userId,
	}
}

