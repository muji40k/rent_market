package product

import "rent_service/internal/logic/services/interfaces/product"

type IProvider interface {
	GetProductService() product.IService
}

type ICharacteristicsProvider interface {
	GetProductCharacteristicsService() product.ICharacteristicsService
}

type IPhotoProvider interface {
	GetProductPhotoService() product.IPhotoService
}

