package access

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=user.go -destination=implementations/mock/user.go

type IUser interface {
	Access(rqUserId uuid.UUID, userId uuid.UUID) error
}

