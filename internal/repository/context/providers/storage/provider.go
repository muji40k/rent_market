package storage

import "rent_service/internal/repository/interfaces/storage"

type IProvider interface {
	GetStorageRepository() storage.IRepository
}

