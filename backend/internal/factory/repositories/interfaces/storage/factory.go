package storage

import "rent_service/internal/repository/interfaces/storage"

type IFactory interface {
	CreateStorageRepository() storage.IRepository
}

