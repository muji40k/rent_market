package category

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/category"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped category.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped category.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) category.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListCategories() (collection.Collection[category.Category], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListCategories",
		func(_ trace.Span) (collection.Collection[category.Category], error) {
			col, err := self.wrapped.ListCategories()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetFullCategory(categoryId uuid.UUID) (collection.Collection[category.Category], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetFullCategory",
		func(span trace.Span) (collection.Collection[category.Category], error) {
			col, err := self.wrapped.GetFullCategory(categoryId)
			span.SetAttributes(
				attribute.Stringer("CategoryId", categoryId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

