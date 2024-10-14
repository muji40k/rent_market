package pickuppoint

import "rent_service/internal/repository/interfaces/pickuppoint"

type IProvider interface {
	GetPickUpPointRepository() pickuppoint.IRepository
}

type IPhotoProvider interface {
	GetPickUpPointPhotoRepository() pickuppoint.IPhotoRepository
}

type IWorkingHoursProvider interface {
	GetPickUpPointWorkingHoursRepository() pickuppoint.IWorkingHoursRepository
}

