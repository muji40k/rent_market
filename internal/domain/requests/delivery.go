package requests

import (
	"github.com/google/uuid"
	"time"
)

type Delivery struct {
	Id                 uuid.UUID
	CompanyId          uuid.UUID
	InstanceId         uuid.UUID
	FromId             uuid.UUID // PickUpPoint
	ToId               uuid.UUID // PickUpPoint
	DeliveryId         string
	ScheduledBeginDate time.Time
	ActualBeginDate    *time.Time
	ScheduledEndDate   time.Time
	ActualEndDate      *time.Time
	VerificationCode   string
	CreateDate         time.Time
}

