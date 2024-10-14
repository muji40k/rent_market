package period

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/period"
	. "rent_service/internal/misc/types/collection"
	period_provider "rent_service/internal/repository/context/providers/period"
)

type repoproviders struct {
	period period_provider.IProvider
}

type service struct {
	repos repoproviders
}

func mapf(value *models.Period) period.Period {
	return period.Period{
		Id:       value.Id,
		Duration: value.Duration,
	}
}

func (self *service) GetPeriods() (Collection[period.Period], error) {
	var periods Collection[period.Period]
	repo := self.repos.period.GetPeriodRepository()
	col, err := repo.GetAll()

	if nil == err {
		periods = MapCollection(mapf, col)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return periods, err
}

