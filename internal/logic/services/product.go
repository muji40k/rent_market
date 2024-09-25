package services

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type ProductSort uint

const (
	PRODUCT_SORT_NONE ProductSort = iota
	PRODUCT_SORT_OFFERS_ASC
	PRODUCT_SORT_OFFERS_DSC
)

type ProductFilter struct {
	CategoryId      uuid.UUID
	Query           *string
	Characteristics []struct {
		Key    string
		Values []string
		Range  *struct {
			Min float64
			Max float64
		}
	}
}

type IProductService interface {
	ListProducts(
		filter ProductFilter,
		sort ProductSort,
	) (Collection[models.Product], error)
	GetProductById(productId uuid.UUID) (models.Product, error)
}

type IProductCharacteristicsService interface {
	GetProductCharacteristics(
		productId uuid.UUID,
	) (models.ProductCharacteristics, error)
}

type IProductPhotoService interface {
	ListProductPhotos(productId uuid.UUID) (Collection[uuid.UUID], error)
}

