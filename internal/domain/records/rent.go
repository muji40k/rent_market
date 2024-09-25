package records

import (
	"github.com/google/uuid"
	"time"
)

type Rent struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	InstanceId      uuid.UUID
	StartDate       time.Time
	EndDate         *time.Time
	PaymentPeriodId uuid.UUID
}

