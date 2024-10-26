package dummy

import (
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/delivery"
	"time"

	"github.com/google/uuid"
)

var Id = uuid.MustParse("44072c0a-e312-452a-aa20-69429d7b950d")

type creator struct {
	counter uint
}

func New() delivery.ICreator {
	return &creator{}
}

func (self *creator) CreateDelivery(
	_ models.Address,
	_ models.Address,
	_ string,
) (delivery.Delivery, error) {
	self.counter++
	now := time.Now()

	return delivery.Delivery{
		CompanyId:          Id,
		DeliveryId:         fmt.Sprint(self.counter),
		ScheduledBeginDate: now.Add(24 * time.Hour),
		ScheduledEndDate:   now.Add(3 * 24 * time.Hour),
	}, nil
}

