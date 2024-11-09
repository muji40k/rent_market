package access

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=rent.go -destination=implementations/mock/rent.go

type IRent interface {
	Access(userId uuid.UUID, rentId uuid.UUID) error
}

type IRentRequest interface {
	Access(userId uuid.UUID, requestId uuid.UUID) error
}

type IRentReturn interface {
	Access(userId uuid.UUID, requestId uuid.UUID) error
}

