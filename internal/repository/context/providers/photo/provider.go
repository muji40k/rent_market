package photo

import "rent_service/internal/repository/interfaces/photo"

type IProvider interface {
	GetPhotoRepository() photo.IRepository
}

type ITempProvider interface {
	GetTempPhotoRepository() photo.ITempRepository
}

