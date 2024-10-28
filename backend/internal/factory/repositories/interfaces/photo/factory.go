package photo

import "rent_service/internal/repository/interfaces/photo"

type IFactory interface {
	CreatePhotoRepository() photo.IRepository
}

type ITempFactory interface {
	CreatePhotoTempRepository() photo.ITempRepository
}

