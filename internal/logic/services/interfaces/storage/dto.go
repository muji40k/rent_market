package storage

import (
	"rent_service/internal/logic/services/types/date"

	"github.com/google/uuid"
)

type Storage struct {
	Id            uuid.UUID  `json:"id"`
	PickUpPointId uuid.UUID  `json:"pick_up_point"`
	InstanceId    uuid.UUID  `json:"instance"`
	InDate        date.Date  `json:"in"`
	OutDate       *date.Date `json:"out"`
}

