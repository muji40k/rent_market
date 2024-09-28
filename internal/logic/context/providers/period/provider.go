package period

import "rent_service/internal/logic/services/interfaces/period"

type IProvider interface {
	GetPeriodService() period.IService
}

