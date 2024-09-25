package models

import (
	"github.com/google/uuid"
	"time"
)

type Photo struct {
	Id          uuid.UUID
	Path        string
	Mime        string
	Placeholder string
	Description string
	Date        time.Time
}

