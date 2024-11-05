package models

import (
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type StorekeeperBuilder struct {
	id            uuid.UUID
	userId        uuid.UUID
	pickUpPointId uuid.UUID
}

func NewStorekeeper() *StorekeeperBuilder {
	return &StorekeeperBuilder{}
}

func (self *StorekeeperBuilder) WithId(id uuid.UUID) *StorekeeperBuilder {
	self.id = id
	return self
}

func (self *StorekeeperBuilder) WithUserId(userId uuid.UUID) *StorekeeperBuilder {
	self.userId = userId
	return self
}

func (self *StorekeeperBuilder) WithPickUpPointId(pickUpPointId uuid.UUID) *StorekeeperBuilder {
	self.pickUpPointId = pickUpPointId
	return self
}

func (self *StorekeeperBuilder) Build() models.Storekeeper {
	return models.Storekeeper{
		Id:            self.id,
		UserId:        self.userId,
		PickUpPointId: self.pickUpPointId,
	}
}

