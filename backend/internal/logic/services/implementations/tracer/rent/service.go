package rent

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/rent"
	"rent_service/internal/logic/services/types/token"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"
)

type service struct {
	wrapped rent.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped rent.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListRentsByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[rent.Rent], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListRentsByUser",
		func(span trace.Span) (collection.Collection[rent.Rent], error) {
			col, err := self.wrapped.ListRentsByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetRentByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (rent.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetRentByInstance",
		func(span trace.Span) (rent.Rent, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetRentByInstance(token, instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) StartRent(token token.Token, form rent.StartForm) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.StartRent",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("StartForm", form),
			)
			return self.wrapped.StartRent(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *service) RejectRent(token token.Token, requestId uuid.UUID) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.RejectRent",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.RejectRent(token, requestId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) StopRent(token token.Token, form rent.StopForm) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.StopRent",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("StopForm", form),
			)
			return self.wrapped.StopRent(token, form)
		},
		cmnerrors.Internal,
	)
}

type requestService struct {
	wrapped rent.IRequestService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRequest(
	wrapped rent.IRequestService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IRequestService {
	return &requestService{wrapped, hl, tracer}
}

func (self *requestService) ListRentRequstsByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[rent.RentRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListRentRequestsByUser",
		func(span trace.Span) (collection.Collection[rent.RentRequest], error) {
			col, err := self.wrapped.ListRentRequstsByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) GetRentRequestByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (rent.RentRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetRentRequestByInstance",
		func(span trace.Span) (rent.RentRequest, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetRentRequestByInstance(token, instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) ListRentRequstsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[rent.RentRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListRentRequstsByPickUpPoint",
		func(span trace.Span) (collection.Collection[rent.RentRequest], error) {
			col, err := self.wrapped.ListRentRequstsByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *requestService) CreateRentRequest(
	token token.Token,
	form rent.RequestCreateForm,
) (rent.RentRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateRentRequest",
		func(span trace.Span) (rent.RentRequest, error) {
			span.SetAttributes(
				wrap.AttributeJSON("CreateForm", form),
			)
			return self.wrapped.CreateRentRequest(token, form)
		},
		cmnerrors.Internal,
	)
}

type returnService struct {
	wrapped rent.IReturnService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewReturn(
	wrapped rent.IReturnService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IReturnService {
	return &returnService{wrapped, hl, tracer}
}

func (self *returnService) ListRentReturnsByUser(
	token token.Token,
	userId uuid.UUID,
) (collection.Collection[rent.ReturnRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListRentReturnsByUser",
		func(span trace.Span) (collection.Collection[rent.ReturnRequest], error) {
			col, err := self.wrapped.ListRentReturnsByUser(token, userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *returnService) GetRentReturnByInstance(
	token token.Token,
	instance uuid.UUID,
) (rent.ReturnRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetRentReturnsByInstance",
		func(span trace.Span) (rent.ReturnRequest, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instance),
			)
			return self.wrapped.GetRentReturnByInstance(token, instance)
		},
		cmnerrors.Internal,
	)
}

func (self *returnService) ListRentReturnsByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[rent.ReturnRequest], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListRentReturnsByPickUpPoint",
		func(span trace.Span) (collection.Collection[rent.ReturnRequest], error) {
			col, err := self.wrapped.ListRentReturnsByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *returnService) CreateRentReturn(
	token token.Token,
	form rent.ReturnCreateForm,
) (rent.ReturnRequest, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateProvisionRevoke",
		func(span trace.Span) (rent.ReturnRequest, error) {
			span.SetAttributes(
				wrap.AttributeJSON("CreateForm", form),
			)
			return self.wrapped.CreateRentReturn(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *returnService) CancelRentReturn(
	token token.Token,
	requestId uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.CancelRentReturn",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.CancelRentReturn(token, requestId)
		},
		cmnerrors.Internal,
	)
}

