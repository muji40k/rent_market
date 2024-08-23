
package models

import (
    "time"
    "github.com/google/uuid"
)

type Photo struct {
    Id uuid.UUID
    Path string
    Placeholder string
    Description string
    Date time.Time
}

