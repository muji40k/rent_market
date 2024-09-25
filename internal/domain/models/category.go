package models

import (
	"github.com/google/uuid"
)

type Category struct {
	Id       uuid.UUID
	ParentId uuid.UUID
	Name     string
}

