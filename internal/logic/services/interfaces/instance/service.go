package instance

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type Sort uint

const (
	SORT_NONE Sort = iota
	SORT_RATING_ASC
	SORT_RATING_DSC
	SORT_DATE_ASC
	SORT_DATE_DSC
	SORT_PRICE_ASC
	SORT_PRICE_DSC
	SORT_USAGE_ASC
	SORT_USAGE_DSC
)

type Filter struct {
	ProductId uuid.UUID
}

type IService interface {
	ListInstances(filter Filter, sort Sort) (Collection[models.Instance], error)
	GetInstanceById(instanceId uuid.UUID) (models.Instance, error)
	UpdateInstance(token models.Token, instance models.Instance) error
}

type PayPlans struct {
	InstanceId uuid.UUID
	Items      map[uuid.UUID]struct {
		PeriodId uuid.UUID // index
		Price    currency.Currency
	}
}

type IPayPlansService interface {
	GetInstancePayPlans(instanceId uuid.UUID) (models.InstancePayPlans, error)
	UpdateInstancePayPlans(
		token models.Token,
		payPlans PayPlans,
	) (models.InstancePayPlans, error)
}

type IPhotoService interface {
	ListInstancePhotos(instanceId uuid.UUID) (Collection[uuid.UUID], error)
	AddInstancePhotos(
		token models.Token,
		instanceId uuid.UUID,
		tempPhotos Collection[uuid.UUID],
	) error
}

type Review struct {
	InstanceId uuid.UUID
	Content    string
	Rating     float64
}

type IReviewService interface {
	ListInstanceReviews(instanceId uuid.UUID) (Collection[models.Review], error)
	PostInstanceReview(token models.Token, review Review) error
}

