package requests

import (
	"rent_service/builders/misc/dategen"
	"time"
)

var getStartDate = dategen.CreateGetter(
	dategen.NewDate(2020, 1, 1),
	dategen.FromTime(time.Now()),
)

