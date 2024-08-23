package records

import (
    "time"
    "github.com/google/uuid"
)

type Provision struct {
    Id uuid.UUID
    RenterId uuid.UUID
    InstanceId uuid.UUID
    StartDate time.Time
    EndDate *time.Time
}

