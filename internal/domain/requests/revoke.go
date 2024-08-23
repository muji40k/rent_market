
package requests

import (
    "time"
    "github.com/google/uuid"
)

type RevokeRequest struct {
    Id uuid.UUID
    InstanceId uuid.UUID
    RenterId uuid.UUID
    PickUpPointId uuid.UUID
    VerificationCode string
    CreateDate time.Time
}

