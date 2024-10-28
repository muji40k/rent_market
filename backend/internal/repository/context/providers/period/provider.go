package period

import "rent_service/internal/repository/interfaces/period"

type IProvider interface {
	GetPeriodRepository() period.IRepository
}

