package delivery

import (
	"rent_service/internal/logic/services/types/date"

	"github.com/google/uuid"
)

type Delivery struct {
	Id         uuid.UUID `json:"id"`
	CompanyId  uuid.UUID `json:"company"`
	InstanceId uuid.UUID `json:"instance"`
	FromId     uuid.UUID `json:"from"`
	ToId       uuid.UUID `json:"to"`
	BeginDate  struct {
		Scheduled date.Date  `json:"scheduled"`
		Actual    *date.Date `json:"actual"`
	} `json:"begin"`
	EndDate struct {
		Scheduled date.Date  `json:"scheduled"`
		Actual    *date.Date `json:"actual"`
	} `json:"end"`
	VerificationCode string    `json:"verification_code"`
	CreateDate       date.Date `json:"create_date"`
}

type DeliveryCompany struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Site        string    `json:"site"`
	PhoneNumber string    `json:"phone_number"`
	Description string    `json:"description"`
}

type CreateForm struct {
	InstanceId uuid.UUID `json:"instance"`
	From       uuid.UUID `json:"from"`
	To         uuid.UUID `json:"to"`
}

type SendForm struct {
	DeliveryId       uuid.UUID
	VerificationCode string
	TempPhotos       []uuid.UUID
}

type AcceptForm struct {
	DeliveryId       uuid.UUID
	Comment          *string
	VerificationCode string
	TempPhotos       []uuid.UUID
}

