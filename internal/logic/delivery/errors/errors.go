package errors

import (
	"fmt"
	"rent_service/internal/domain/models"

	"github.com/google/uuid"
)

type ErrorAddressCantBeReached struct {
	company uuid.UUID
	Address models.Address
}
type ErrorOverloaded struct{ company uuid.UUID }
type ErrorRejected struct {
	company uuid.UUID
	Reason  string
}
type ErrorInternal struct{ Err error }

// Creators
func AddressCantBeReached(
	company uuid.UUID,
	address models.Address,
) ErrorAddressCantBeReached {
	return ErrorAddressCantBeReached{company, address}
}

func Overloaded(company uuid.UUID) ErrorOverloaded {
	return ErrorOverloaded{company}
}

func Rejected(company uuid.UUID, reason string) ErrorRejected {
	return ErrorRejected{company, reason}
}

func Internal(err error) ErrorInternal {
	return ErrorInternal{err}
}

// Error implementation
func (self ErrorAddressCantBeReached) Error() string {
	return fmt.Sprintf(
		"Address '%v' can't be reached by company '%d'",
		self.Address, self.company,
	)
}

func (self ErrorOverloaded) Error() string {
	return fmt.Sprintf("Company '%v' is overloaded", self.company)
}

func (self ErrorRejected) Error() string {
	return fmt.Sprintf(
		"Company '%v' rejected delivery: %v", self.company, self.Reason,
	)
}

func (self ErrorInternal) Error() string {
	return fmt.Sprintf("Internal error during delivery request: %v", self.Err)
}

func (self ErrorInternal) Unwrap() error {
	return self.Err
}

