package period

import "rent_service/internal/repository/interfaces/period"

type IFactory interface {
	CreatePeriodRepository() period.IRepository
}

