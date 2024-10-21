package product

import "rent_service/internal/repository/interfaces/product"

type IFactory interface {
	CreateProductRepository() product.IRepository
}

type ICharacteristicsFactory interface {
	CreateProductCharacteristicsRepository() product.ICharacteristicsRepository
}

type IPhotoFactory interface {
	CreateProductPhotoRepository() product.IPhotoRepository
}

