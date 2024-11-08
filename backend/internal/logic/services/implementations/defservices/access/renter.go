package access

import (
	"github.com/google/uuid"
)

type IRenter interface {
	Access(userId uuid.UUID, renterUserId uuid.UUID) error
}

