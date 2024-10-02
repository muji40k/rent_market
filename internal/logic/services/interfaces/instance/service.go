package instance

import (
	"fmt"
	"rent_service/internal/logic/services/types/token"
	. "rent_service/internal/misc/types/collection"

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
	ListInstances(filter Filter, sort Sort) (Collection[Instance], error)
	GetInstanceById(instanceId uuid.UUID) (Instance, error)
	UpdateInstance(token token.Token, instance Instance) error
}

type IPayPlansService interface {
	GetInstancePayPlans(instanceId uuid.UUID) (Collection[PayPlan], error)
	UpdateInstancePayPlans(
		token token.Token,
		instanceId uuid.UUID,
		payPlans PayPlansUpdateForm,
	) error
}

type IPhotoService interface {
	ListInstancePhotos(instanceId uuid.UUID) (Collection[uuid.UUID], error)
	AddInstancePhotos(
		token token.Token,
		instanceId uuid.UUID,
		tempPhotos []uuid.UUID,
	) error
}

type ReviewSort uint

const (
	REVIEW_SORT_NONE ReviewSort = iota
	REVIEW_SORT_DATE_ASC
	REVIEW_SORT_DATE_DSC
	REVIEW_SORT_RATING_ASC
	REVIEW_SORT_RATING_DSC
)

type ReviewRating uint

type ReviewFilter struct {
	InstanceId uuid.UUID
	Ratings    []ReviewRating
}

type IReviewService interface {
	ListInstanceReviews(filter ReviewFilter, sort ReviewSort) (Collection[Review], error)
	PostInstanceReview(
		token token.Token,
		instanceId uuid.UUID,
		review ReviewPostForm,
	) error
}

type ErrorRatingIncorrectValue struct{ Value ReviewRating }

func (e ErrorRatingIncorrectValue) Error() string {
	return fmt.Sprintf(
		"Incorrect rating value '%v' exceeded range [0; 5]",
		e.Value,
	)
}

