package instance

import "rent_service/internal/logic/services/interfaces/instance"

type IProvider interface {
	GetInstanceService() instance.IService
}

type IPayPlansProvider interface {
	GetInstancePayPlansService() instance.IPayPlansService
}

type IPhotoProvider interface {
	GetInstancePhotoService() instance.IPhotoService
}

type IReviewProvider interface {
	GetInstanceReviewService() instance.IReviewService
}

