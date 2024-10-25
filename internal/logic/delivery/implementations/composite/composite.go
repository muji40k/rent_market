package composite

import (
	"fmt"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/delivery"
	"rent_service/internal/logic/delivery/errors"

	"github.com/google/uuid"
)

type DeliveryComposite struct {
	deliveries map[uuid.UUID]delivery.ICreator
}

type Dpair struct {
	id      uuid.UUID
	creator delivery.ICreator
}

func Pair(id uuid.UUID, creator delivery.ICreator) Dpair {
	return Dpair{id, creator}
}

func New(deliveries ...Dpair) DeliveryComposite {
	m := make(map[uuid.UUID]delivery.ICreator, len(deliveries))

	for _, v := range deliveries {
		m[v.id] = v.creator
	}

	return DeliveryComposite{m}
}

func (self *DeliveryComposite) Add(id uuid.UUID, creator delivery.ICreator) {
	self.deliveries[id] = creator
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

