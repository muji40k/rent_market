package photo

import (
	"github.com/google/uuid"
	"rent_service/internal/logic/services/types/date"
)

type Description struct {
	Mime        string `json:"mime" binding:"required"`
	Placeholder string `json:"placeholder" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type Photo struct {
	Id          uuid.UUID `json:"id"`
	Mime        string    `json:"mime"`
	Placeholder string    `json:"placeholder"`
	Description string    `json:"description"`
	Href        string    `json:"href"`
	Date        date.Date `json:"date"`
}

type TempPhoto struct {
	Id          uuid.UUID
	Mime        string
	Placeholder string
	Description string
	Date        date.Date
}

