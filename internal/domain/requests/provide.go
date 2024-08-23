
package requests

import (
    "time"
    "github.com/google/uuid"

    "rent_service/internal/domain/models"
)

type ProvideRequest struct {
    Id uuid.UUID
    ProductId uuid.UUID
    RenterId uuid.UUID
    PickUpPointId uuid.UUID
    Name string
    Description string
    VerificationCode string
    CreateDate time.Time
}

type ProvideRequestPayPlans struct {
    ProvideRequestId uuid.UUID
    Map map[uuid.UUID]models.PayPlan // Indexed by Period uuid
}

func NewProvideRequestPayPlans() ProvideRequestPayPlans {
    out := ProvideRequestPayPlans{}
    out.Map = make(map[uuid.UUID]models.PayPlan)

    return out
}

