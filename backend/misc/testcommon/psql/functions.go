package psql

import (
	"rent_service/internal/logic/services/types/date"
	"rent_service/misc/nullable"
	"time"
)

func UnwrapDate(d *date.Date) *time.Time {
	return nullable.GetOr(
		nullable.Map(
			nullable.FromPtr(d),
			func(d *date.Date) *time.Time {
				return &d.Time
			},
		),
		nil,
	)
}

func CompareTimeMicro(t1 time.Time, t2 time.Time) bool {
	return t1.UnixMicro() == t2.UnixMicro()
}

func CompareTimePtrMicro(t1 *time.Time, t2 *time.Time) bool {
	if nil != t1 && nil != t2 {
		return t1.UnixMicro() == t2.UnixMicro()
	} else if nil == t1 && nil == t2 {
		return true
	} else {
		return false
	}
}

