package records

import (
    "time"
    "github.com/google/uuid"
)

type Storage struct {
    Id uuid.UUID
    PickUpPointId uuid.UUID
    InstanceId uuid.UUID
    InDate time.Time
    OutDate *time.Time
}

