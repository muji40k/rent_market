package dategen

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const format = "2006-01-02"

type formattable struct {
	year  uint
	month uint
	day   uint
}

type Date struct {
	time *time.Time
	form *formattable
}

func FromTime(t time.Time) Date {
	out := Date{new(time.Time), nil}
	*out.time = t

	return out
}

func NewDate(year uint, month uint, day uint) Date {
	out := Date{nil, new(formattable)}
	*out.form = formattable{year, month, day}

	return out
}

func (self *formattable) String() string {
	return fmt.Sprintf("%v-%v-%v", self.year, self.month, self.day)
}

func (self *Date) getTime() time.Time {
	if nil == self.form && nil != self.time {
		return *self.time
	} else if nil != self.form && nil == self.time {
		t, _ := time.Parse(format, self.form.String())
		return t
	} else {
		panic("Unreachable state!")
	}
}

func CreateGetter(from Date, to Date) func() time.Time {
	begin := from.getTime()
	end := to.getTime()
	diff := end.Sub(begin)

	return func() time.Time {
		return begin.Add(time.Duration(float64(diff) * rand.Float64()))
	}

}

func GetDate(from Date, to Date) time.Time {
	begin := from.getTime()
	end := to.getTime()
	diff := end.Sub(begin)

	return begin.Add(time.Duration(float64(diff) * rand.Float64()))
}

