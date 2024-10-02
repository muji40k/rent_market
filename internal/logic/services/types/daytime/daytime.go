package daytime

import (
	"errors"
	"time"
)

type Time struct {
	Time time.Time
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

