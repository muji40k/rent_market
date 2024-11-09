package access

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=instance.go -destination=implementations/mock/instance.go

type IInstance interface {
	Access(userId uuid.UUID, instanceId uuid.UUID) error
}

