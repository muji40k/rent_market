package storage

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/storage"
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
	wrapped storage.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped storage.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) storage.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListStoragesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[storage.Storage], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListStoragesByPickUpPoint",
		func(span trace.Span) (collection.Collection[storage.Storage], error) {
			col, err := self.wrapped.ListStoragesByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetStorageByInstance(
	instanceId uuid.UUID,
) (storage.Storage, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetStorageByInstance",
		func(span trace.Span) (storage.Storage, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetStorageByInstance(instanceId)
		},
		cmnerrors.Internal,
	)
}

