package requests

import (
	"github.com/google/uuid"
	"time"
)

type Revoke struct {
	Id               uuid.UUID
	InstanceId       uuid.UUID
	RenterId         uuid.UUID
	PickUpPointId    uuid.UUID
	VerificationCode string
	CreateDate       time.Time
}

