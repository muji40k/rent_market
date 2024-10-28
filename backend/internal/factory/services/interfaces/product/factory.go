package product

import "rent_service/internal/logic/services/interfaces/product"

type IFactory interface {
	CreateProductService() product.IService
}

type ICharacteristicsFactory interface {
	CreateProductCharacteristicsService() product.ICharacteristicsService
}

type IPhotoFactory interface {
	CreateProductPhotoService() product.IPhotoService
}

