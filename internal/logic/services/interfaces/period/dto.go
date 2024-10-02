package period

import (
	"time"

	"github.com/google/uuid"
)

type Period struct {
	Id       uuid.UUID     `json:"id"`
	Duration time.Duration `json:"duration"`
}

