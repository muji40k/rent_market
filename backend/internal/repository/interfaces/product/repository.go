package product

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type Sort uint

const (
	SORT_NONE Sort = iota
	SORT_OFFERS_ASC
	SORT_OFFERS_DSC
)

type Filter struct {
	CategoryId uuid.UUID
	Query      *string
	Ranges     []Range
	Selectors  []Selector
}

type Characteristic struct {
	Key string
}

type Range struct {
	Characteristic
	Min float64
	Max float64
}

type Selector struct {
	Characteristic
	Values []string
}

type IRepository interface {
	GetById(productId uuid.UUID) (models.Product, error)
	GetWithFilter(
		filter Filter,
		sort Sort,
	) (collection.Collection[models.Product], error)
}

type ICharacteristicsRepository interface {
	GetByProductId(productId uuid.UUID) (models.ProductCharacteristics, error)
}

type IPhotoRepository interface {
	GetByProductId(productId uuid.UUID) (collection.Collection[uuid.UUID], error)
}

