package access

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=pick_up_point.go -destination=implementations/mock/pick_up_point.go

type IPickUpPoint interface {
	Access(userId uuid.UUID, pickUpPointId uuid.UUID) error
}

