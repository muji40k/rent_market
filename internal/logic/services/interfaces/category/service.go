package category

import (
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListCategories() (Collection[Category], error)
	GetFullCategory(categoryId uuid.UUID) (Collection[Category], error)
}

