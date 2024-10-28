package records

import (
	"github.com/google/uuid"
	"time"
)

type Provision struct {
	Id         uuid.UUID
	RenterId   uuid.UUID
	InstanceId uuid.UUID
	StartDate  time.Time
	EndDate    *time.Time
}

