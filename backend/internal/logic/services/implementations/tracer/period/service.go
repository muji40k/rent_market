package period

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/period"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped period.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped period.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) period.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) GetPeriods() (collection.Collection[period.Period], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPeriods",
		func(_ trace.Span) (collection.Collection[period.Period], error) {
			col, err := self.wrapped.GetPeriods()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

