package models

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	Id         uuid.UUID
	InstanceId uuid.UUID
	UserId     uuid.UUID
	Content    string
	Rating     float64
	Date       time.Time
}

