package category

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IService interface {
	ListCategories() (Collection[models.Category], error)
	GetFullCategory(categoryId uuid.UUID) (Collection[models.Category], error)
}

