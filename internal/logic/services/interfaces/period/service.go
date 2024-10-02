package period

import (
	. "rent_service/internal/misc/types/collection"
)

type IService interface {
	GetPeriods() (Collection[Period], error)
}

