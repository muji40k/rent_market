package provide

import (
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type Provision struct {
	Id         uuid.UUID  `json:"id"`
	UserId     uuid.UUID  `json:"user"`
	InstanceId uuid.UUID  `json:"instance"`
	StartDate  date.Date  `json:"start_date"`
	EndDate    *date.Date `json:"end_date"`
}

type PayPlan struct {
	Id       uuid.UUID         `json:"id"`
	PeriodId uuid.UUID         `json:"period"`
	Price    currency.Currency `json:"price"`
}

type ProvideRequest struct {
	Id               uuid.UUID `json:"id"`
	ProductId        uuid.UUID `json:"product"`
	UserId           uuid.UUID `json:"user"`
	PickUpPointId    uuid.UUID `json:"pick_up_point"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Condition        string    `json:"condition"`
	PayPlans         []PayPlan `json:"pay_plans"`
	VerificationCode string    `json:"verification_code"`
	CreateDate       date.Date `json:"create_date"`
}

type RevokeRequest struct {
	Id               uuid.UUID `json:"id"`
	InstanceId       uuid.UUID `json:"instance"`
	UserId           uuid.UUID `json:"user"`
	PickUpPointId    uuid.UUID `json:"pick_up_point"`
	VerificationCode string    `json:"verification_code"`
	CreateDate       date.Date `json:"create_date"`
}

type Overrides struct {
	ProductId   *uuid.UUID `json:"product"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Condition   *string    `json:"condition"`
	PayPlans    []PayPlan  `json:"pay_plans"`
}

type StartForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	Overrides        Overrides
	TempPhotos       []uuid.UUID
}

type StopForm struct {
	RevokeId         uuid.UUID
	VerificationCode string
	TempPhotos       []uuid.UUID
}

type PayPlanCreateForm struct {
	PeriodId uuid.UUID         `json:"period" binding:"required"`
	Price    currency.Currency `json:"price" binding:"required"`
}

type RequestCreateForm struct {
	ProductId     uuid.UUID           `json:"product" binding:"required"`
	PickUpPointId uuid.UUID           `json:"pick_up_point" binding:"required"`
	Name          string              `json:"name" binding:"required"`
	Description   string              `json:"description" binding:"required"`
	Condition     string              `json:"condition" binding:"required"`
	PayPlans      []PayPlanCreateForm `json:"pay_plans" binding:"required"`
}

type RevokeCreateForm struct {
	ProvisionId   uuid.UUID `json:"provision" binding:"required"`
	PickUpPointId uuid.UUID `json:"pick_up_point" binding:"required"`
}

