package delivery

import (
	"rent_service/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type Delivery struct {
	CompanyId          uuid.UUID
	DeliveryId         string
	ScheduledBeginDate time.Time
	ScheduledEndDate   time.Time
}

type ICreator interface {
	CreateDelivery(
		from models.Address,
		to models.Address,
		verificationCode string,
	) (Delivery, error)
	CancelDelivery(companyId uuid.UUID, deliveryId string) error
}

