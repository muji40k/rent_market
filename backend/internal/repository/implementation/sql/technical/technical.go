package technical

import "time"

type Info struct {
	MDate   time.Time `db:"modification_date"`
	MSource string    `db:"modification_source"`
}

type ISetter interface {
	Update(info *Info)
}

