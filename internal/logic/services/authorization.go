package services

import (
	"fmt"

	"github.com/google/uuid"
)

type Role uint

const (
	USER Role = iota
	ADMINISTRATOR
	STOREKEEPER
	RENTER
	role_max
)

type IAuthorizationService interface {
	Authorize(token Token, role Role) error
}

func (role *Role) String() string {
	switch *role {
	case USER:
		return "user"
	case ADMINISTRATOR:
		return "administrator"
	case STOREKEEPER:
		return "storekeeper"
	case RENTER:
		return "renter"
	default:
		return "unknown"
	}
}

func (role *Role) Valid() error {
	if *role >= role_max {
		return ErrorUnknown{[]string{
			fmt.Sprintf("Role value '%d' >= '%d' [max]", role, role_max),
		}}
	}

	return nil
}

type ErrorUnauthorized struct {
	id   uuid.UUID
	Role Role
}
type ErrorUnsupportedRole struct{ Role Role }

func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf(
		"User '%v' attempt to authorize to insufficient role '%v'",
		e.id, e.Role,
	)
}

func (e ErrorUnsupportedRole) Error() string {
	return fmt.Sprintf("Unsupported role value '%d'", e.Role)
}

