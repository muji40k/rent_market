package period

import "rent_service/internal/logic/services/interfaces/period"

type IFactory interface {
	CreatePeriodService() period.IService
}

