package models

import "time"

type Meta struct {
	MDate   time.Time `gorm:"column:modification_date"`
	MSource string    `gorm:"column:modification_source"`
}

