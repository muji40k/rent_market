package payment

import (
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/payment"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped payment.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped payment.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) payment.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetByInstanceId(
	instanceId uuid.UUID,
) (collection.Collection[models.Payment], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Payment.GetByInstanceId",
		func(span trace.Span) (collection.Collection[models.Payment], error) {
			col, err := self.wrapped.GetByInstanceId(instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByRentId(
	rentId uuid.UUID,
) (collection.Collection[models.Payment], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Payment.GetByRentId",
		func(span trace.Span) (collection.Collection[models.Payment], error) {
			col, err := self.wrapped.GetByRentId(rentId)
			span.SetAttributes(
				attribute.Stringer("RentId", rentId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

