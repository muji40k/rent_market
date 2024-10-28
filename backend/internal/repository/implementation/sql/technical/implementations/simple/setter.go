package simple

import (
	"rent_service/internal/repository/implementation/sql/technical"
	"time"
)

type setter struct {
	source string
}

func New(source string) technical.ISetter {
	return &setter{source}
}

func (self *setter) Update(info *technical.Info) {
	info.MDate = time.Now()
	info.MSource = self.source
}

