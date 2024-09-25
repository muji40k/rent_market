package requests

import (
	"github.com/google/uuid"
	"time"
)

type Rent struct {
	Id               uuid.UUID
	InstanceId       uuid.UUID
	UserId           uuid.UUID
	PickUpPointId    uuid.UUID
	PaymentPeriodId  uuid.UUID
	VerificationCode string
	CreateDate       time.Time
}

