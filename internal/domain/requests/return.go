
package requests

import (
    "time"
    "github.com/google/uuid"
)

type ReturnRequest struct {
    Id uuid.UUID
    InstanceId uuid.UUID
    UserId uuid.UUID
    PickUpPointId uuid.UUID
    RentEndDate time.Time
    VerificationCode string
    CreateDate time.Time
}

