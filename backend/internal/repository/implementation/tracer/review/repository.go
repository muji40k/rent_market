package review

import (
	"fmt"
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/review"
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped review.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped review.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) review.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(
	review models.Review,
) (models.Review, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Review.Create",
		func(span trace.Span) (models.Review, error) {
			out, err := self.wrapped.Create(review)
			span.SetAttributes(
				wrap.AttributeJSON("Review", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetWithFilter(
	filter review.Filter,
	sort review.Sort,
) (collection.Collection[models.Review], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Review.GetWithFilter",
		func(span trace.Span) (collection.Collection[models.Review], error) {
			col, err := self.wrapped.GetWithFilter(filter, sort)
			span.SetAttributes(
				wrap.AttributeJSON("Filter", filter),
				attribute.String("Sort", func() string {
					switch sort {
					case review.SORT_NONE:
						return "NONE"
					case review.SORT_DATE_ASC:
						return "DATE_ASC"
					case review.SORT_DATE_DSC:
						return "DATE_DSC"
					case review.SORT_RATING_ASC:
						return "RATING_ASC"
					case review.SORT_RATING_DSC:
						return "RATING_DSC"
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

