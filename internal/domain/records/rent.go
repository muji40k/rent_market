package records

import (
    "time"
    "github.com/google/uuid"
)

type Rent struct {
    Id uuid.UUID
    UserId uuid.UUID
    InstanceId uuid.UUID
    StartDate time.Time
    EndDate *time.Time
    PaymentPeriodId uuid.UUID
}

