package services

import (
	"rent_service/internal/domain/models"
	. "rent_service/internal/misc/types/collection"
	"rent_service/internal/misc/types/currency"

	"github.com/google/uuid"
)

type InstanceSort uint

const (
	INSTANCE_SORT_NONE ProductSort = iota
	INSTANCE_SORT_RATING_ASC
	INSTANCE_SORT_RATING_DSC
	INSTANCE_SORT_DATE_ASC
	INSTANCE_SORT_DATE_DSC
	INSTANCE_SORT_PRICE_ASC
	INSTANCE_SORT_PRICE_DSC
	INSTANCE_SORT_USAGE_ASC
	INSTANCE_SORT_USAGE_DSC
)

type InstanceFilter struct {
	ProductId uuid.UUID
}

type IInstanceService interface {
	ListInstances(
		filter InstanceFilter,
		sort InstanceSort,
	) (Collection[models.Instance], error)
	GetInstanceById(instanceId uuid.UUID) (models.Instance, error)
	UpdateInstance(token Token, instance models.Instance) error
}

type InstancePayPlans struct {
	InstanceId uuid.UUID
	Items      map[uuid.UUID]struct {
		PeriodId uuid.UUID // index
		Price    currency.Currency
	}
}

type IInstancePayPlansService interface {
	GetInstancePayPlans(instanceId uuid.UUID) (models.InstancePayPlans, error)
	UpdateInstancePayPlans(
		token Token,
		payPlans InstancePayPlans,
	) (models.InstancePayPlans, error)
}

type IInstancePhotoService interface {
	ListInstancePhotos(instanceId uuid.UUID) (Collection[uuid.UUID], error)
	AddInstancePhotos(
		token Token,
		instanceId uuid.UUID,
		tempPhotos Collection[uuid.UUID],
	) error
}

type InstanceReview struct {
	InstanceId uuid.UUID
	Content    string
	Rating     float64
}

type IInstanceReviewService interface {
	ListInstanceReviews(instanceId uuid.UUID) (Collection[models.Review], error)
	PostInstanceReview(token Token, review InstanceReview) error
}

