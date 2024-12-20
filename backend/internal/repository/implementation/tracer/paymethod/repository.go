package paymethod

import (
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/paymethod"
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped paymethod.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped paymethod.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) paymethod.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetAll() (collection.Collection[models.PayMethod], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.PayMethod.GetAll",
		func(span trace.Span) (collection.Collection[models.PayMethod], error) {
			col, err := self.wrapped.GetAll()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

