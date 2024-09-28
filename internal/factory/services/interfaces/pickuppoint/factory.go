package pickuppoint

import "rent_service/internal/logic/services/interfaces/pickuppoint"

type IFactory interface {
	CreatePickUpPointService() pickuppoint.IService
}

type IPhotoFactory interface {
	CreatePickUpPointPhotoService() pickuppoint.IPhotoService
}

type IWorkingHoursFactory interface {
	CreatePickUpPointWorkingHoursService() pickuppoint.IWorkingHoursService
}

