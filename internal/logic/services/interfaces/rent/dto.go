package rent

import (
	"rent_service/internal/logic/services/types/date"

	"github.com/google/uuid"
)

type Rent struct {
	Id              uuid.UUID  `json:"id"`
	UserId          uuid.UUID  `json:"user"`
	InstanceId      uuid.UUID  `json:"instance"`
	StartDate       date.Date  `json:"start_date"`
	EndDate         *date.Date `json:"end_date"`
	PaymentPeriodId uuid.UUID  `json:"payment_period"`
}

type RentRequest struct {
	Id               uuid.UUID `json:"id"`
	InstanceId       uuid.UUID `json:"instance"`
	UserId           uuid.UUID `json:"user"`
	PickUpPointId    uuid.UUID `json:"pick_up_point"`
	PaymentPeriodId  uuid.UUID `json:"payment_period"`
	VerificationCode string    `json:"verification_code"`
	CreateDate       date.Date `json:"create_date"`
}

type ReturnRequest struct {
	Id               uuid.UUID `json:"id"`
	InstanceId       uuid.UUID `json:"instance"`
	UserId           uuid.UUID `json:"user"`
	PickUpPointId    uuid.UUID `json:"pick_up_point"`
	RentEndDate      date.Date `json:"rent_end_date"`
	VerificationCode string    `json:"verification_code"`
	CreateDate       date.Date `json:"create_date"`
}

type StartForm struct {
	RequestId        uuid.UUID
	VerificationCode string
	TempPhotos       []uuid.UUID
}

type StopForm struct {
	ReturnId         uuid.UUID
	VerificationCode string
	Comment          *string
	TempPhotos       []uuid.UUID
}

type RequestCreateForm struct {
	InstanceId      uuid.UUID `json:"instance" binding:"required"`
	PickUpPointId   uuid.UUID `json:"pick_up_point" binding:"required"`
	PaymentPeriodId uuid.UUID `json:"payment_period" binding:"required"`
}

type ReturnCreateForm struct {
	RentId        uuid.UUID `json:"rent" binding:"required"`
	PickUpPointId uuid.UUID `json:"pick_up_point" binding:"required"`
	EndDate       date.Date `json:"rent_end_date" binding:"required"`
}

