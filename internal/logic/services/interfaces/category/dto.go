package category

import "github.com/google/uuid"

type Category struct {
	Id       uuid.UUID `json:"id"`
	ParentId uuid.UUID `json:"parent"`
	Name     string    `json:"name"`
}

