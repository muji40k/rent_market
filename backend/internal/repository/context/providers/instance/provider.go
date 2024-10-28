package instance

import "rent_service/internal/repository/interfaces/instance"

type IProvider interface {
	GetInstanceRepository() instance.IRepository
}

type IPayPlansProvider interface {
	GetInstancePayPlansRepository() instance.IPayPlansRepository
}

type IPhotoProvider interface {
	GetInstancePhotoRepository() instance.IPhotoRepository
}

