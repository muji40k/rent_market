package access

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=renter.go -destination=implementations/mock/renter.go

type IRenter interface {
	Access(userId uuid.UUID, renterUserId uuid.UUID) error
}

