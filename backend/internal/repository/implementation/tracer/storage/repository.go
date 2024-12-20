package storage

import (
	"rent_service/internal/domain/records"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/storage"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped storage.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped storage.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) storage.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(
	storage records.Storage,
) (records.Storage, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Storage.Create",
		func(span trace.Span) (records.Storage, error) {
			out, err := self.wrapped.Create(storage)
			span.SetAttributes(
				wrap.AttributeJSON("Storage", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(
	storage records.Storage,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Storage.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Storage", storage),
			)
			return self.wrapped.Update(storage)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetActiveByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[records.Storage], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Storage.GetActiveByPickUpPointId",
		func(span trace.Span) (collection.Collection[records.Storage], error) {
			col, err := self.wrapped.GetActiveByPickUpPointId(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (records.Storage, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Storage.GetActiveByInstanceId",
		func(span trace.Span) (records.Storage, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetActiveByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

