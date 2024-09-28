package models

import "github.com/google/uuid"

type Address struct {
	Id      uuid.UUID
	Country string
	City    string
	Street  string
	House   string
	Flat    *string
}

