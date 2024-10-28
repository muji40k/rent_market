package instance

import "rent_service/internal/repository/interfaces/instance"

type IFactory interface {
	CreateInstanceRepository() instance.IRepository
}

type IPayPlansFactory interface {
	CreateInstancePayPlansRepository() instance.IPayPlansRepository
}

type IPhotoFactory interface {
	CreateInstancePhotoRepository() instance.IPhotoRepository
}

