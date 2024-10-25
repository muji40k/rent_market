package dummy

import (
	"rent_service/internal/logic/services/implementations/defservices/payment"

	"github.com/google/uuid"
)

var Id = uuid.MustParse("a2e7d44f-0af0-4b9b-bd70-65d3397d5ad9")

type checker struct{}

func New() payment.IRegistrationChecker {
	return &checker{}
}

func (self *checker) CheckPayerId(payerId string) error {
	return nil
}

func (self *checker) MethodId() uuid.UUID {
	return Id
}

