package product

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CategoryId  uuid.UUID `json:"category"`
	Description string    `json:"description"`
}

type Charachteristic struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`
}

