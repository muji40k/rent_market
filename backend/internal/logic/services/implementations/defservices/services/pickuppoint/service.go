package pickuppoint

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/pickuppoint"
	"rent_service/internal/logic/services/types/day"
	"rent_service/internal/logic/services/types/daytime"
	. "rent_service/internal/misc/types/collection"
	pickuppoint_provider "rent_service/internal/repository/context/providers/pickuppoint"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	pickUpPoint pickuppoint_provider.IProvider
}

type service struct {
	repos repoproviders
}

func mapPoint(value *models.PickUpPoint) pickuppoint.PickUpPoint {
	return pickuppoint.PickUpPoint{
		Id: value.Id,
		Address: pickuppoint.Address{
			Country: value.Address.Country,
			City:    value.Address.City,
			Street:  value.Address.Street,
			House:   value.Address.House,
			Flat:    value.Address.Flat,
		},
		Capacity: value.Capacity,
	}
}

func New(pickUpPoint pickuppoint_provider.IProvider) pickuppoint.IService {
	return &service{repoproviders{pickUpPoint}}
}

func (self *service) ListPickUpPoints() (Collection[pickuppoint.PickUpPoint], error) {
	var points Collection[pickuppoint.PickUpPoint]
	repo := self.repos.pickUpPoint.GetPickUpPointRepository()
	col, err := repo.GetAll()

	if nil == err {
		points = MapCollection(mapPoint, col)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return points, err
}

func (self *service) GetPickUpPointById(
	pickUpPointId uuid.UUID,
) (pickuppoint.PickUpPoint, error) {
	repo := self.repos.pickUpPoint.GetPickUpPointRepository()
	pickUpPoint, err := repo.GetById(pickUpPointId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return mapPoint(&pickUpPoint), err
}

type photoRepoproviders struct {
	photo pickuppoint_provider.IPhotoProvider
}

type photoService struct {
	repos photoRepoproviders
}

func NewPhoto(photo pickuppoint_provider.IPhotoProvider) pickuppoint.IPhotoService {
	return &photoService{photoRepoproviders{photo}}
}

func (self *photoService) ListPickUpPointPhotos(
	pickUpPointId uuid.UUID,
) (Collection[uuid.UUID], error) {
	repo := self.repos.photo.GetPickUpPointPhotoRepository()
	photos, err := repo.GetById(pickUpPointId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return photos, err
}

type workingHoursRepoproviders struct {
	wh pickuppoint_provider.IWorkingHoursProvider
}

type workingHoursService struct {
	repos workingHoursRepoproviders
}

func NewWorkingHours(
	wh pickuppoint_provider.IWorkingHoursProvider,
) pickuppoint.IWorkingHoursService {
	return &workingHoursService{workingHoursRepoproviders{wh}}
}

func mapWorkingHours(
	wh *models.PickUpPointWorkingHours,
) Collection[pickuppoint.WorkingHours] {
	buf := make([]pickuppoint.WorkingHours, len(wh.Map))
	i := 0

	for k, v := range wh.Map {
		buf[i].Id = v.Id
		buf[i].Day = day.New(k)
		buf[i].StartHour = daytime.NewDuration(v.Begin)
		buf[i].EndHour = daytime.NewDuration(v.End)
		i++
	}

	return SliceCollection(buf)
}

func (self *workingHoursService) ListPickUpPointWorkingHours(
	pickUpPointId uuid.UUID,
) (Collection[pickuppoint.WorkingHours], error) {
	var wh Collection[pickuppoint.WorkingHours]
	repo := self.repos.wh.GetPickUpPointWorkingHoursRepository()
	mwh, err := repo.GetById(pickUpPointId)

	if nil == err {
		wh = mapWorkingHours(&mwh)
	} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return wh, err
}

