package models

import (
	"github.com/google/uuid"
	"time"
)

type PickUpPoint struct {
	Id       uuid.UUID
	Address  Address
	Capacity uint64
}

type WorkingHours struct {
	Id    uuid.UUID
	Day   time.Weekday
	Begin time.Duration // Both from 00:00
	End   time.Duration
}

type PickUpPointWorkingHours struct {
	PickUpPointId uuid.UUID
	Map           map[time.Weekday]WorkingHours
}

