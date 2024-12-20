package category

import (
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/category"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped category.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped category.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) category.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetAll() (collection.Collection[models.Category], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Category.GetAll",
		func(span trace.Span) (collection.Collection[models.Category], error) {
			col, err := self.wrapped.GetAll()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetPath(leaf uuid.UUID) (collection.Collection[models.Category], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Category.GetPath",
		func(span trace.Span) (collection.Collection[models.Category], error) {
			col, err := self.wrapped.GetPath(leaf)
			span.SetAttributes(
				attribute.Stringer("Leaf", leaf),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

