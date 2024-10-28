package product

import "rent_service/internal/repository/interfaces/product"

type IProvider interface {
	GetProductRepository() product.IRepository
}

type ICharacteristicsProvider interface {
	GetProductCharacteristicsRepository() product.ICharacteristicsRepository
}

type IPhotoProvider interface {
	GetProductPhotoRepository() product.IPhotoRepository
}

