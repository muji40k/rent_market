package daytime

import (
	"errors"
	"time"
)

type Time struct {
	Time time.Time
}

func New(time time.Time) Time {
	return Time{time}
}

func NewDuration(duration time.Duration) Time {
	var base time.Time

	return Time{base.Add(duration)}
}

func (self *Time) ToDuration() time.Duration {
	return self.Time.Sub(self.Time.Truncate(24 * time.Hour))
}

const FORMAT = "15:04"

func (self *Time) MarshalJSON() ([]byte, error) {
	return []byte("\"" + self.Time.Format(FORMAT) + "\""), nil
}

func (self *Time) UnmarshalJSON(bytes []byte) error {
	if 2 > len(bytes) || '"' != bytes[0] || '"' != bytes[len(bytes)-1] {
		return errors.New("Not a json string")
	}

	bytes = bytes[1 : len(bytes)-1]

	if time, err := time.Parse(FORMAT, string(bytes)); nil != err {
		return err
	} else {
		self.Time = time
		return nil
	}
}

