package psql

import "time"

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

