package models

import "github.com/google/uuid"

type Administrator struct {
	Id     uuid.UUID
	UserId uuid.UUID
}

