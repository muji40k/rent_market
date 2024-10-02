package payment

import (
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type PayMethod struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type UserPayMethod struct {
	Id       uuid.UUID `json:"id"`
	MethodId uuid.UUID `json:"pay_method"`
	Name     string    `json:"name"`
}

type PayMethodRegistrationForm struct {
	MethodId uuid.UUID `json:"pay_method"`
	PayerId  string    `json:"payer_id"`
	Name     string    `json:"name"`
}

type Payment struct {
	Id          uuid.UUID         `json:"id"`
	RentId      uuid.UUID         `json:"rent"`
	PayMethodId *uuid.UUID        `json:"pay_method"`
	PeriodStart date.Date         `json:"period_start"`
	PeriodEnd   date.Date         `json:"period_end"`
	Value       currency.Currency `json:"price"`
	Status      string            `json:"status"`
	CreateDate  date.Date         `json:"create_date"`
	PaymentDate *date.Date        `json:"payment_date"`
}

