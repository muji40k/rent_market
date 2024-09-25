package records

import (
	"github.com/google/uuid"
	"time"
)

type Storage struct {
	Id            uuid.UUID
	PickUpPointId uuid.UUID
	InstanceId    uuid.UUID
	InDate        time.Time
	OutDate       *time.Time
}

