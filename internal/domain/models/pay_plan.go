
package models

import (
    "github.com/google/uuid"
    "rent_service/internal/misc/types/currency"
)

type PayPlan struct {
    Id uuid.UUID
    PeriodId uuid.UUID
    Price currency.Currency
}

