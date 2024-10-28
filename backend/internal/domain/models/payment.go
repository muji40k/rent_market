package models

import (
	"github.com/google/uuid"
	"time"

	"rent_service/internal/misc/types/currency"
)

type Payment struct {
	Id          uuid.UUID
	RentId      uuid.UUID
	PayMethodId *uuid.UUID
	PaymentId   *string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Value       currency.Currency
	Status      string
	CreateDate  time.Time
	PaymentDate *time.Time
}

