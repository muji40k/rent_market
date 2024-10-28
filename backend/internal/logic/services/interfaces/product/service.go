package product

import (
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type Sort uint

const (
	SORT_NONE Sort = iota
	SORT_OFFERS_ASC
	SORT_OFFERS_DSC
)

type Filter struct {
	CategoryId      uuid.UUID
	Query           *string
	Characteristics []FilterCharachteristic
}

type FilterCharachteristic struct {
	Key    string
	Values []string
	Range  *struct {
		Min float64
		Max float64
	}
}

type IService interface {
	ListProducts(filter Filter, sort Sort) (Collection[Product], error)
	GetProductById(productId uuid.UUID) (Product, error)
}

type ICharacteristicsService interface {
	ListProductCharacteristics(
		productId uuid.UUID,
	) (Collection[Charachteristic], error)
}

type IPhotoService interface {
	ListProductPhotos(productId uuid.UUID) (Collection[uuid.UUID], error)
}

