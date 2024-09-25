package models

import (
	"github.com/google/uuid"
	"time"
)

type Period struct {
	Id       uuid.UUID
	Name     string
	Duration time.Duration
}

