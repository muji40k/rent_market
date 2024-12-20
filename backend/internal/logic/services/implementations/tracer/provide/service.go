package provide

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/provide"
	"rent_service/internal/logic/services/types/token"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"
)

type service struct {
	wrapped provide.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped provide.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provide.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListProvisionsByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[provide.Provision], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProvisionsByUser",
		func(span trace.Span) (collection.Collection[provide.Provision], error) {
			col, err := self.wrapped.ListProvisionsByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetProvisionByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (provide.Provision, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetProvisionByInstance",
		func(span trace.Span) (provide.Provision, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetProvisionByInstance(token, instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) StartProvision(
	token token.Token,
	form provide.StartForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.StartProvision",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("StartForm", form),
			)
			return self.wrapped.StartProvision(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *service) RejectProvision(
	token token.Token,
	requestId uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.RejectProvision",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.RejectProvision(token, requestId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) StopProvision(
	token token.Token,
	form provide.StopForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.StopProvision",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("StopForm", form),
			)
			return self.wrapped.StopProvision(token, form)
		},
		cmnerrors.Internal,
	)
}

type requestService struct {
	wrapped provide.IRequestService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRequest(
	wrapped provide.IRequestService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provide.IRequestService {
	return &requestService{wrapped, hl, tracer}
}

func (self *requestService) ListProvisionRequstsByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[provide.ProvideRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProvisionRequestsByUser",
		func(span trace.Span) (collection.Collection[provide.ProvideRequest], error) {
			col, err := self.wrapped.ListProvisionRequstsByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) GetProvisionRequestByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (provide.ProvideRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetProvisionRequestByInstance",
		func(span trace.Span) (provide.ProvideRequest, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetProvisionRequestByInstance(token, instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) ListProvisionRequstsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[provide.ProvideRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProvisionRequstsByPickUpPoint",
		func(span trace.Span) (collection.Collection[provide.ProvideRequest], error) {
			col, err := self.wrapped.ListProvisionRequstsByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) CreateProvisionRequest(
	token token.Token,
	form provide.RequestCreateForm,
) (provide.ProvideRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateProvisionRequest",
		func(span trace.Span) (provide.ProvideRequest, error) {
			span.SetAttributes(
				wrap.AttributeJSON("CreateForm", form),
			)
			return self.wrapped.CreateProvisionRequest(token, form)
		},
		cmnerrors.Internal,
	)
}

type revokeService struct {
	wrapped provide.IRevokeService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRevoke(
	wrapped provide.IRevokeService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) provide.IRevokeService {
	return &revokeService{wrapped, hl, tracer}
}

func (self *revokeService) ListProvisionRevokesByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[provide.RevokeRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProvisionRevokesByUser",
		func(span trace.Span) (collection.Collection[provide.RevokeRequest], error) {
			col, err := self.wrapped.ListProvisionRevokesByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *revokeService) GetProvisionRevokeByInstance(
	token token.Token,
	instance uuid.UUID,
) (provide.RevokeRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetProvisionRevokeByInstance",
		func(span trace.Span) (provide.RevokeRequest, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instance),
			)
			return self.wrapped.GetProvisionRevokeByInstance(token, instance)
		},
		cmnerrors.Internal,
	)
}

func (self *revokeService) ListProvisionRetvokesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[provide.RevokeRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProvisionRetvokesByPickUpPoint",
		func(span trace.Span) (collection.Collection[provide.RevokeRequest], error) {
			col, err := self.wrapped.ListProvisionRetvokesByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *revokeService) CreateProvisionRevoke(
	token token.Token,
	form provide.RevokeCreateForm,
) (provide.RevokeRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateProvisionRevoke",
		func(span trace.Span) (provide.RevokeRequest, error) {
			span.SetAttributes(
				wrap.AttributeJSON("CreateForm", form),
			)
			return self.wrapped.CreateProvisionRevoke(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *revokeService) CancelProvisionRevoke(
	token token.Token,
	requestId uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.CancelProvisionRevoke",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.CancelProvisionRevoke(token, requestId)
		},
		cmnerrors.Internal,
	)
}

