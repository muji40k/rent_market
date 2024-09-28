package composite

import (
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/delivery"
	"rent_service/internal/logic/delivery/errors"
)

type DeliveryComposite struct {
	deliveries []delivery.IDelivery
}

func New(deliveries []delivery.IDelivery) DeliveryComposite {
	return DeliveryComposite{deliveries}
}

func (self *DeliveryComposite) Add(delivery delivery.IDelivery) {
	self.deliveries = append(self.deliveries, delivery)
}

func (self *DeliveryComposite) CreateDelivery(
	from models.Address,
	to models.Address,
	code string,
) (delivery.Delivery, error) {
	var errs = make([]error, 0, len(self.deliveries))

	for _, delivery := range self.deliveries {
		if res, err := delivery.CreateDelivery(from, to, code); nil == err {
			return res, nil
		} else {
			errs = append(errs, err)
		}
	}

	return delivery.Delivery{}, errors.Internal(ErrorComposition{errs})
}

type ErrorComposition struct{ Errors []error }

func (self ErrorComposition) Error() string {
	return fmt.Sprintf("All delivery requests failed: %v", self.Errors)
}

