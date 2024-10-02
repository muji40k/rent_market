package payment

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"
)

type IPayMethodService interface {
	GetPayMethods() (Collection[PayMethod], error)
}

type IUserPayMethodService interface {
	GetPayMethods(token token.Token) (Collection[UserPayMethod], error)
	RegisterPayMethod(
		token token.Token,
		method PayMethodRegistrationForm,
	) (uuid.UUID, error)
	UpdatePayMethodsPriority(
		token token.Token,
		methodsOrder Collection[uuid.UUID],
	) error
	RemovePayMethod(token token.Token, methodId uuid.UUID) (bool, error)
}

var ErrorIncompletePayMethodsList = errors.New(
	"Priority list doesn't use all registered methods",
)

type ErrorWrongPayerId struct{ Id string }

func (e ErrorWrongPayerId) Error() string {
	return fmt.Sprintf("Can't verify payer id '%v'", e.Id)
}

type IRentPaymentService interface {
	GetPaymentsByInstance(
		token token.Token,
		instanceId uuid.UUID,
	) (Collection[Payment], error)
	GetPaymentsByRentId(
		token token.Token,
		rentId uuid.UUID,
	) (Collection[Payment], error)
}

