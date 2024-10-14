package models

import "github.com/google/uuid"

type Renter struct {
	Id     uuid.UUID
	UserId uuid.UUID
}

