package instance

import (
	"fmt"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/instance"
	"rent_service/internal/logic/services/types/token"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped instance.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped instance.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListInstances(
	filter instance.Filter,
	sort instance.Sort,
) (collection.Collection[instance.Instance], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListInstances",
		func(span trace.Span) (collection.Collection[instance.Instance], error) {
			col, err := self.wrapped.ListInstances(filter, sort)
			span.SetAttributes(
				wrap.AttributeJSON("Filter", filter),
				attribute.String("Sort", func() string {
					switch sort {
					case instance.SORT_NONE:
						return "NONE"
					case instance.SORT_RATING_ASC:
						return "RATING_ASC"
					case instance.SORT_RATING_DSC:
						return "RATING_DSC"
					case instance.SORT_DATE_ASC:
						return "DATE_ASC"
					case instance.SORT_DATE_DSC:
						return "DATE_DSC"
					case instance.SORT_PRICE_ASC:
						return "PRICE_ASC"
					case instance.SORT_PRICE_DSC:
						return "PRICE_DSC"
					case instance.SORT_USAGE_ASC:
						return "USAGE_ASC"
					case instance.SORT_USAGE_DSC:
						return "USAGE_DSC"
					default:
						return "UNKNOWN: " + fmt.Sprint(sort)
					}
				}()),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetInstanceById(instanceId uuid.UUID) (instance.Instance, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetInstanceById",
		func(span trace.Span) (instance.Instance, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetInstanceById(instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) UpdateInstance(token token.Token, instance instance.Instance) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateInstance",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Instance", instance),
			)
			return self.wrapped.UpdateInstance(token, instance)
		},
		cmnerrors.Internal,
	)
}

type payPlansService struct {
	wrapped instance.IPayPlansService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPayPlans(
	wrapped instance.IPayPlansService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IPayPlansService {
	return &payPlansService{wrapped, hl, tracer}
}

func (self *payPlansService) GetInstancePayPlans(
	instanceId uuid.UUID,
) (collection.Collection[instance.PayPlan], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetInstancePayPlans",
		func(span trace.Span) (collection.Collection[instance.PayPlan], error) {
			col, err := self.wrapped.GetInstancePayPlans(instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *payPlansService) UpdateInstancePayPlans(
	token token.Token,
	instanceId uuid.UUID,
	payPlans instance.PayPlansUpdateForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdateInstancePayPlans",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
				wrap.AttributeJSON("UpdateForm", payPlans),
			)
			return self.wrapped.UpdateInstancePayPlans(token, instanceId, payPlans)
		},
		cmnerrors.Internal,
	)
}

type photoService struct {
	wrapped instance.IPhotoService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped instance.IPhotoService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IPhotoService {
	return &photoService{wrapped, hl, tracer}
}

func (self *photoService) ListInstancePhotos(
	instanceId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListInstancePhotos",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.ListInstancePhotos(instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *photoService) AddInstancePhotos(
	token token.Token,
	instanceId uuid.UUID,
	tempPhotos []uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.AddInstancePhotos",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
				wrap.AttributeJSON("TempPhotos", tempPhotos),
			)
			return self.wrapped.AddInstancePhotos(token, instanceId, tempPhotos)
		},
		cmnerrors.Internal,
	)
}

type reviewService struct {
	wrapped instance.IReviewService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewReview(
	wrapped instance.IReviewService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IReviewService {
	return &reviewService{wrapped, hl, tracer}
}

func (self *reviewService) ListInstanceReviews(
	filter instance.ReviewFilter,
	sort instance.ReviewSort,
) (collection.Collection[instance.Review], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListInstanceReviews",
		func(span trace.Span) (collection.Collection[instance.Review], error) {
			col, err := self.wrapped.ListInstanceReviews(filter, sort)
			span.SetAttributes(
				wrap.AttributeJSON("Filter", filter),
				attribute.String("Sort", func() string {
					switch sort {
					case instance.REVIEW_SORT_NONE:
						return "NONE"
					case instance.REVIEW_SORT_DATE_ASC:
						return "DATE_ASC"
					case instance.REVIEW_SORT_DATE_DSC:
						return "DATE_DSC"
					case instance.REVIEW_SORT_RATING_ASC:
						return "RATING_ASC"
					case instance.REVIEW_SORT_RATING_DSC:
						return "RATING_DSC"
					default:
						return "UNKNOWN: " + fmt.Sprint(sort)
					}
				}()),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *reviewService) PostInstanceReview(
	token token.Token,
	instanceId uuid.UUID,
	review instance.ReviewPostForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.PostInstanceReview",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
				wrap.AttributeJSON("PostForm", review),
			)
			return self.wrapped.PostInstanceReview(token, instanceId, review)
		},
		cmnerrors.Internal,
	)
}

