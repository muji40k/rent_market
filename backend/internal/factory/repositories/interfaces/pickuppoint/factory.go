package pickuppoint

import "rent_service/internal/repository/interfaces/pickuppoint"

type IFactory interface {
	CreatePickUpPointRepository() pickuppoint.IRepository
}

type IPhotoFactory interface {
	CreatePickUpPointPhotoRepository() pickuppoint.IPhotoRepository
}

type IWorkingHoursFactory interface {
	CreatePickUpPointWorkingHoursRepository() pickuppoint.IWorkingHoursRepository
}

