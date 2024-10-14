package models

import (
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	Name        string
	CategoryId  uuid.UUID
	Description string
}

type Charachteristic struct {
	Id    uuid.UUID
	Name  string
	Value string
}

type ProductCharacteristics struct {
	ProductId uuid.UUID
	Map       map[string]Charachteristic
}

// func NewProductCharacteristics() ProductCharacteristics {
//     out := ProductCharacteristics{}
//     out.Map = make(map[string]Charachteristic)
//
//     return out
// }

