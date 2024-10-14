package models

import "github.com/google/uuid"

type Storekeeper struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	PickUpPointId uuid.UUID
}

