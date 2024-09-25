package requests

import (
	"github.com/google/uuid"
	"time"

	"rent_service/internal/domain/models"
)

type Provide struct {
	Id               uuid.UUID
	ProductId        uuid.UUID
	RenterId         uuid.UUID
	PickUpPointId    uuid.UUID
	PayPlans         map[uuid.UUID]models.PayPlan // Indexed by Period uuid
	Name             string
	Description      string
	Condition        string
	VerificationCode string
	CreateDate       time.Time
}

func NewProvide() Provide {
	out := Provide{}
	out.PayPlans = make(map[uuid.UUID]models.PayPlan)

	return out
}

