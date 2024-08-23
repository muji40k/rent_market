
package models

import (
    "github.com/google/uuid"
)

type PayMethod struct {
    Id uuid.UUID
    Name string
    Description string
}

