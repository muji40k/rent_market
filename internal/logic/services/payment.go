package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
)

type PayMethod struct {
	Id   uuid.UUID
	Name string
}

type IPayMethodService interface {
	GetPayMethods() (Collection[PayMethod], error)
}

type UserPayMethod struct {
	Id     uuid.UUID
	Method PayMethod
}

type PayMethodRegistrationForm struct {
	MethodId uuid.UUID
	Priority uint
	PayerId  string
}

type IUserPaymentService interface {
	GetPayMethods(token Token) (Collection[UserPayMethod], error)
	RegisterPayMethod(token Token, method PayMethodRegistrationForm) error
	UpdatePayMethodsPriority(
		token Token,
		methodsOrder Collection[uuid.UUID],
	) error
	RemovePayMethod(token Token, methodId uuid.UUID) (bool, error)
}

var ErrorIncompletePayMethodsList = errors.New(
	"Priority list doesn't use all registered methods",
)

type ErrorWrongPayerId struct{ id string }

func (e ErrorWrongPayerId) Error() string {
	return fmt.Sprintf("Can't verify payer id '%v'", e.id)
}

type IRentPaymentService interface {
	GetPaymentsByInstance(
		token Token,
		instanceId uuid.UUID,
	) (Collection[models.Payment], error)
	GetPaymentsByRentId(
		token Token,
		rentId uuid.UUID,
	) (Collection[models.Payment], error)
}

