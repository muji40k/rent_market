package day

import (
	"errors"
	"strconv"
	"time"
)

type Day struct {
	Day time.Weekday
}

func New(day time.Weekday) Day {
	return Day{day}
}

func toMondayStart(incorrect int) int {
	incorrect -= 1

	if 0 > incorrect {
		return 7 - incorrect
	}

	return incorrect
}

func toSundayStart(correct int) int {
	return (correct + 1) % 7
}

func (self Day) MarshalJSON() ([]byte, error) {
	return []byte(string(strconv.Itoa(toMondayStart(int(self.Day))))), nil
}

func (self *Day) UnmarshalJSON(bytes []byte) error {
	if value, err := strconv.Atoi(string(bytes)); nil == err {
		if 0 <= value && 6 >= value {
			self.Day = time.Weekday(toSundayStart(value))

			return nil
		} else {
			return errors.New("Day value exceed [0; 6]")
		}
	} else {
		return err
	}
}

