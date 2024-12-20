package instance

import (
	"fmt"
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/instance"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped instance.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped instance.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(
	instance models.Instance,
) (models.Instance, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Instance.Create",
		func(span trace.Span) (models.Instance, error) {
			out, err := self.wrapped.Create(instance)
			span.SetAttributes(
				wrap.AttributeJSON("Instance", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(instance models.Instance) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Instance.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Instance", instance),
			)
			return self.wrapped.Update(instance)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(instanceId uuid.UUID) (models.Instance, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Instance.GetById",
		func(span trace.Span) (models.Instance, error) {
			span.SetAttributes(
				wrap.AttributeJSON("InstanceId", instanceId),
			)
			return self.wrapped.GetById(instanceId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetWithFilter(
	filter instance.Filter,
	sort instance.Sort,
) (collection.Collection[models.Instance], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Instance.GetWithFilter",
		func(span trace.Span) (collection.Collection[models.Instance], error) {
			col, err := self.wrapped.GetWithFilter(filter, sort)
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
		func(err error) error { return err },
	)
}

type payPlansRepository struct {
	wrapped instance.IPayPlansRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPayPlans(
	wrapped instance.IPayPlansRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IPayPlansRepository {
	return &payPlansRepository{wrapped, hl, tracer}
}

func (self *payPlansRepository) Create(
	payPlans models.InstancePayPlans,
) (models.InstancePayPlans, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.InstancePayPlans.Create",
		func(span trace.Span) (models.InstancePayPlans, error) {
			out, err := self.wrapped.Create(payPlans)
			span.SetAttributes(
				wrap.AttributeJSON("PayPlans", payPlans),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *payPlansRepository) AddPayPlan(
	instanceId uuid.UUID,
	plan models.PayPlan,
) (models.InstancePayPlans, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.InstancePayPlans.Create",
		func(span trace.Span) (models.InstancePayPlans, error) {
			out, err := self.wrapped.AddPayPlan(instanceId, plan)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
				wrap.AttributeJSON("PayPlan", plan),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *payPlansRepository) Update(payPlans models.InstancePayPlans) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.InstancePayPlans.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("PayPlans", payPlans),
			)
			return self.wrapped.Update(payPlans)
		},
		func(err error) error { return err },
	)
}

func (self *payPlansRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (models.InstancePayPlans, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.InstancePayPlans.GetByInstanceId",
		func(span trace.Span) (models.InstancePayPlans, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

type photoRepository struct {
	wrapped instance.IPhotoRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped instance.IPhotoRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) instance.IPhotoRepository {
	return &photoRepository{wrapped, hl, tracer}
}

func (self *photoRepository) Create(
	instanceId uuid.UUID,
	photoId uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Instance.Create",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
				attribute.Stringer("PhotoId", photoId),
			)
			return self.wrapped.Create(instanceId, photoId)
		},
		func(err error) error { return err },
	)
}

func (self *photoRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Instance.GetByInstanceId",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.GetByInstanceId(instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

