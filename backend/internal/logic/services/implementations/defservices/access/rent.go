package access

import (
	"github.com/google/uuid"
)

type IRent interface {
	Access(userId uuid.UUID, rentId uuid.UUID) error
}

type IRentRequest interface {
	Access(userId uuid.UUID, requestId uuid.UUID) error
}

type IRentReturn interface {
	Access(userId uuid.UUID, requestId uuid.UUID) error
}

