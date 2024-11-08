package authorizer

import (
	"fmt"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type IAuthorizer interface {
	IsAdministrator(userId uuid.UUID) (models.Administrator, error)
	IsRenter(userId uuid.UUID) (models.Renter, error)
	IsStorekeeper(userId uuid.UUID) (models.Storekeeper, error)
}

func Unauthorized(id uuid.UUID, role string) ErrorUnauthorized {
	return ErrorUnauthorized{id, role}
}

type ErrorUnauthorized struct {
	id   uuid.UUID
	Role string
}

func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf(
		"User '%v' attempt to authorize to insufficient role '%v'",
		e.id, e.Role,
	)
}

