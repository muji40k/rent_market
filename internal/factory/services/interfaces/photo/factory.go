package photo

import "rent_service/internal/logic/services/interfaces/photo"

type IFactory interface {
	CreatePhotoService() photo.IService
}

