package instance

import "rent_service/internal/logic/services/interfaces/instance"

type IFactory interface {
	CreateInstanceService() instance.IService
}

type IPayPlansFactory interface {
	CreateInstancePayPlansService() instance.IPayPlansService
}

type IPhotoFactory interface {
	CreateInstancePhotoService() instance.IPhotoService
}

type IReviewFactory interface {
	CreateInstanceReviewService() instance.IReviewService
}

