package models

import (
    "github.com/google/uuid"
)

type DeliveryCompany struct {
    Id uuid.UUID
    Name string
    Site string
    PhoneNumber string
    Description string
}

