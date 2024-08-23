
package requests

import (
    "time"
    "github.com/google/uuid"
)

type RentRequest struct {
    Id uuid.UUID
    InstanceId uuid.UUID
    UserId uuid.UUID
    PickUpPointId uuid.UUID
    PaymentPeriodId uuid.UUID
    VerificationCode string
    CreateDate time.Time
}

