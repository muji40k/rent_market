package storage

import "rent_service/internal/logic/services/interfaces/storage"

type IProvider interface {
	GetStorageService() storage.IService
}

