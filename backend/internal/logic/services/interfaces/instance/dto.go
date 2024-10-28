package instance

import (
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type Instance struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"product"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Condition   string    `json:"condition"`
}

type PayPlan struct {
	InstanceId uuid.UUID         `json:"id"`
	PeriodId   uuid.UUID         `json:"period"`
	Price      currency.Currency `json:"price"`
}

type PayPlanUpdateForm struct {
	PeriodId uuid.UUID         `json:"period"`
	Price    currency.Currency `json:"price"`
}

type PayPlansUpdateForm []PayPlanUpdateForm

type Review struct {
	Id         uuid.UUID `json:"id"`
	InstanceId uuid.UUID `json:"instance"`
	UserId     uuid.UUID `json:"user"`
	Content    string    `json:"content"`
	Rating     float64   `json:"rating"`
	Date       date.Date `json:"date"`
}

type ReviewPostForm struct {
	Content string  `json:"content"`
	Rating  float64 `json:"rating"`
}

