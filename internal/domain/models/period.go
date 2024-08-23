
package models

import (
    "time"
    "github.com/google/uuid"
)

type Period struct {
    Id uuid.UUID
    Name string
    Duration time.Duration
}

