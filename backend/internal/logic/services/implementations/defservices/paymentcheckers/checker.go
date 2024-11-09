package paymentcheckers

import "github.com/google/uuid"

//go:generate mockgen -source=checker.go -destination=mock/mock.go

type IRegistrationChecker interface {
	MethodId() uuid.UUID
	CheckPayerId(payerId string) error
}

