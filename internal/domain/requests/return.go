package requests

import (
	"github.com/google/uuid"
	"time"
)

type Return struct {
	Id               uuid.UUID
	InstanceId       uuid.UUID
	UserId           uuid.UUID
	PickUpPointId    uuid.UUID
	RentEndDate      time.Time
	VerificationCode string
	CreateDate       time.Time
}

