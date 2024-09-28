package pickuppoint

import "rent_service/internal/logic/services/interfaces/pickuppoint"

type IProvider interface {
	GetPickUpPointService() pickuppoint.IService
}

type IPhotoProvider interface {
	GetPickUpPointPhotoService() pickuppoint.IPhotoService
}

type IWorkingHoursProvider interface {
	GetPickUpPointWorkingHoursService() pickuppoint.IWorkingHoursService
}

