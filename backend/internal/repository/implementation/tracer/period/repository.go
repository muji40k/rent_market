package period

import (
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/period"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped period.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped period.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) period.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetById(periodId uuid.UUID) (models.Period, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Period.GetById",
		func(span trace.Span) (models.Period, error) {
			span.SetAttributes(
				attribute.Stringer("PeriodId", periodId),
			)
			return self.wrapped.GetById(periodId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetAll() (collection.Collection[models.Period], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Period.GetAll",
		func(span trace.Span) (collection.Collection[models.Period], error) {
			col, err := self.wrapped.GetAll()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

