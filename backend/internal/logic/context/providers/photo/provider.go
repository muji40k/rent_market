package photo

import "rent_service/internal/logic/services/interfaces/photo"

type IProvider interface {
	GetPhotoService() photo.IService
}

