package review

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"

	"github.com/google/uuid"
)

type Sort uint

const (
	SORT_NONE Sort = iota
	SORT_DATE_ASC
	SORT_DATE_DSC
	SORT_RATING_ASC
	SORT_RATING_DSC
)

type Rating uint

type Filter struct {
	InstanceId uuid.UUID
	Ratings    []Rating
}

type IRepository interface {
	Create(review models.Review) (models.Review, error)

	GetWithFilter(filter Filter, sort Sort) (Collection[models.Review], error)
}
