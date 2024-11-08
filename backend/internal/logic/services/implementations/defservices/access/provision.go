package access

import (
	"github.com/google/uuid"
)

type IProvision interface {
	Access(userId uuid.UUID, provisionId uuid.UUID) error
}

type IProvisionRequest interface {
	Access(userId uuid.UUID, requestId uuid.UUID) error
}

type IProvisionRevoke interface {
	Access(userId uuid.UUID, revokeId uuid.UUID) error
}

