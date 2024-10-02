package date

import (
	"errors"
	"time"
)

type Date struct {
	Time time.Time
}

const FORMAT = "2006-01-02"

func (self *Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + self.Time.Format(FORMAT) + "\""), nil
}

func (self *Date) UnmarshalJSON(bytes []byte) error {
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

