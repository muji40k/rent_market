package utctime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type UTCTime struct {
	Time time.Time
}

func FromTime(t time.Time) UTCTime {
	return UTCTime{t}
}

func (self UTCTime) Value() (driver.Value, error) {
	return self.Time.In(time.UTC), nil
}

func (self *UTCTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		self.Time = t.In(time.UTC)
		return nil
	}

	return fmt.Errorf("Unsupproted type for UTCTime: %T", src)
}

