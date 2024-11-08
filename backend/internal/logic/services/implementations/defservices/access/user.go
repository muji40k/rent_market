package access

import (
	"github.com/google/uuid"
)

type IUser interface {
	Access(rqUserId uuid.UUID, userId uuid.UUID) error
}

