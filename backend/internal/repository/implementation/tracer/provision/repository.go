package provision

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/provision"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped provision.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped provision.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provision.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(
	provision records.Provision,
) (records.Provision, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Provision.Create",
		func(span trace.Span) (records.Provision, error) {
			out, err := self.wrapped.Create(provision)
			span.SetAttributes(
				wrap.AttributeJSON("Provision", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(provision records.Provision) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Provision.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Provision", provision),
			)
			return self.wrapped.Update(provision)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(
	provisionId uuid.UUID,
) (records.Provision, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Provision.GetById",
		func(span trace.Span) (records.Provision, error) {
			span.SetAttributes(
				attribute.Stringer("ProvisionId", provisionId),
			)
			return self.wrapped.GetById(provisionId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByRenterUserId(
	userId uuid.UUID,
) (collection.Collection[records.Provision], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Provision.GetByRenterUserId",
		func(span trace.Span) (collection.Collection[records.Provision], error) {
			col, err := self.wrapped.GetByRenterUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByInstanceId(
	instanceId uuid.UUID,
) (collection.Collection[records.Provision], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Provision.GetByInstanceId",
		func(span trace.Span) (collection.Collection[records.Provision], error) {
			col, err := self.wrapped.GetByInstanceId(instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (records.Provision, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Provision.GetActiveByInstanceId",
		func(span trace.Span) (records.Provision, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetActiveByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

type requestRepository struct {
	wrapped provision.IRequestRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRequest(
	wrapped provision.IRequestRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provision.IRequestRepository {
	return &requestRepository{wrapped, hl, tracer}
}

func (self *requestRepository) Create(
	request requests.Provide,
) (requests.Provide, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRequest.Create",
		func(span trace.Span) (requests.Provide, error) {
			out, err := self.wrapped.Create(request)
			span.SetAttributes(
				wrap.AttributeJSON("Request", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *requestRepository) GetById(
	requestId uuid.UUID,
) (requests.Provide, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRequest.GetById",
		func(span trace.Span) (requests.Provide, error) {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.GetById(requestId)
		},
		func(err error) error { return err },
	)
}

func (self *requestRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Provide], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRequest.GetByUserId",
		func(span trace.Span) (collection.Collection[requests.Provide], error) {
			col, err := self.wrapped.GetByUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *requestRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Provide, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRequest.GetByInstanceId",
		func(span trace.Span) (requests.Provide, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

func (self *requestRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Provide], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRequest.GetByPickUpPointId",
		func(span trace.Span) (collection.Collection[requests.Provide], error) {
			col, err := self.wrapped.GetByPickUpPointId(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *requestRepository) Remove(requestId uuid.UUID) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.ProvisionRequest.Remove",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.Remove(requestId)
		},
		func(err error) error { return err },
	)
}

type revokeRepository struct {
	wrapped provision.IRevokeRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRevoke(
	wrapped provision.IRevokeRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provision.IRevokeRepository {
	return &revokeRepository{wrapped, hl, tracer}
}

func (self *revokeRepository) Create(
	request requests.Revoke,
) (requests.Revoke, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRevoke.Create",
		func(span trace.Span) (requests.Revoke, error) {
			out, err := self.wrapped.Create(request)
			span.SetAttributes(
				wrap.AttributeJSON("Revoke", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *revokeRepository) GetById(
	requestId uuid.UUID,
) (requests.Revoke, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRevoke.GetById",
		func(span trace.Span) (requests.Revoke, error) {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.GetById(requestId)
		},
		func(err error) error { return err },
	)
}

func (self *revokeRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Revoke], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRevoke.GetByUserId",
		func(span trace.Span) (collection.Collection[requests.Revoke], error) {
			col, err := self.wrapped.GetByUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UsereId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *revokeRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Revoke, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRevoke.GetByInstanceId",
		func(span trace.Span) (requests.Revoke, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

func (self *revokeRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Revoke], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProvisionRevoke.GetByPickUpPointId",
		func(span trace.Span) (collection.Collection[requests.Revoke], error) {
			col, err := self.wrapped.GetByPickUpPointId(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *revokeRepository) Remove(requestId uuid.UUID) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.ProvisionRevoke.Remove",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RevokeId", requestId),
			)
			return self.wrapped.Remove(requestId)
		},
		func(err error) error { return err },
	)
}

