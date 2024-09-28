package storage

import "rent_service/internal/logic/services/interfaces/storage"

type IFactory interface {
	CreateStorageService() storage.IService
}

