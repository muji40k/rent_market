
package models

import (
    "time"
    "github.com/google/uuid"
)

type PickUpPoint struct {
    Id uuid.UUID
    Address string
    Capacity uint64
}

type WorkingHours struct {
    Id uuid.UUID
    Day time.Weekday
    Begin time.Duration // Both from 00:00
    End time.Duration
}

type PickUpPointWorkingHours struct {
    PickUpPointId uuid.UUID
    Map map[time.Weekday]WorkingHours
}

type PickUpPointPhoto struct {
    Id uuid.UUID
    PickUpPointId uuid.UUID
    Photo Photo
}

func NewPickUpPointWorkingHours() PickUpPointWorkingHours {
    out := PickUpPointWorkingHours{}
    out.Map = make(map[time.Weekday]WorkingHours)

    return out
}

