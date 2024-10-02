package pickuppoint

import (
	"github.com/google/uuid"
	"rent_service/internal/logic/services/types/day"
	"rent_service/internal/logic/services/types/daytime"
)

type Address struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Street  string  `json:"street"`
	House   string  `json:"house"`
	Flat    *string `json:"flat"`
}

type PickUpPoint struct {
	Id       uuid.UUID `json:"id"`
	Address  Address   `json:"address"`
	Capacity uint64    `json:"capacity"`
}

type WorkingHours struct {
	Id        uuid.UUID    `json:"id"`
	Day       day.Day      `json:"day"`
	StartHour daytime.Time `json:"start_hour"`
	EndHour   daytime.Time `json:"end_hour"`
}

