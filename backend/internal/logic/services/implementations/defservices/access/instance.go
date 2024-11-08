package access

import (
	"github.com/google/uuid"
)

type IInstance interface {
	Access(userId uuid.UUID, instanceId uuid.UUID) error
}

