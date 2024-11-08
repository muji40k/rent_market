package mock

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/delivery"
	derrors "rent_service/internal/logic/delivery/errors"

	"github.com/google/uuid"
)

var Id = uuid.MustParse("d3bf06de-7e21-442a-8296-734d127b8f24")

type MockDelivery struct {
	createDelivery func(
		from models.Address,
		to models.Address,
		verificationCode string,
	) (delivery.Delivery, error)
}

func New() *MockDelivery {
	return &MockDelivery{
		func(
			from models.Address,
			to models.Address,
			verificationCode string,
		) (delivery.Delivery, error) {
			return delivery.Delivery{}, derrors.Internal(errors.New("Method not set"))
		},
	}
}

func (self *MockDelivery) WithCreateDelivery(f func(
	from models.Address,
	to models.Address,
	verificationCode string,
) (delivery.Delivery, error)) *MockDelivery {
	self.createDelivery = f
	return self
}

func (self *MockDelivery) CreateDelivery(
	from models.Address,
	to models.Address,
	verificationCode string,
) (delivery.Delivery, error) {
	return self.createDelivery(from, to, verificationCode)
}

