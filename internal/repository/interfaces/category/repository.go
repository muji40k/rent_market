package category

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type IRepository interface {
	GetAll() (collection.Collection[models.Category], error)
	GetPath(leaf uuid.UUID) (collection.Collection[models.Category], error)
}

