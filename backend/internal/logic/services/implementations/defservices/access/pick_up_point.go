package access

import (
	"github.com/google/uuid"
)

type IPickUpPoint interface {
	Access(userId uuid.UUID, pickUpPointId uuid.UUID) error
}

