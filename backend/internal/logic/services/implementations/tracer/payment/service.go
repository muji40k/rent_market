package payment

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/payment"
	"rent_service/internal/logic/services/types/token"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"
)

type payMethodService struct {
	wrapped payment.IPayMethodService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPayMethod(
	wrapped payment.IPayMethodService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) payment.IPayMethodService {
	return &payMethodService{wrapped, hl, tracer}
}

func (self *payMethodService) GetPayMethods() (collection.Collection[payment.PayMethod], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPayMethods",
		func(span trace.Span) (collection.Collection[payment.PayMethod], error) {
			col, err := self.wrapped.GetPayMethods()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

type userPayMethodService struct {
	wrapped payment.IUserPayMethodService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewUserPayMethod(
	wrapped payment.IUserPayMethodService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) payment.IUserPayMethodService {
	return &userPayMethodService{wrapped, hl, tracer}
}

func (self *userPayMethodService) GetPayMethods(
	token token.Token,
) (collection.Collection[payment.UserPayMethod], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPayMethods",
		func(span trace.Span) (collection.Collection[payment.UserPayMethod], error) {
			col, err := self.wrapped.GetPayMethods(token)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *userPayMethodService) RegisterPayMethod(
	token token.Token,
	method payment.PayMethodRegistrationForm,
) (uuid.UUID, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.RegisterPayMethods",
		func(span trace.Span) (uuid.UUID, error) {
			return self.wrapped.RegisterPayMethod(token, method)
		},
		cmnerrors.Internal,
	)
}

func (self *userPayMethodService) UpdatePayMethodsPriority(
	token token.Token,
	methodsOrder []uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UpdatePayMethodsPriority",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("MethodsOrder", methodsOrder),
			)
			return self.wrapped.UpdatePayMethodsPriority(token, methodsOrder)
		},
		cmnerrors.Internal,
	)
}

func (self *userPayMethodService) RemovePayMethod(
	token token.Token,
	methodId uuid.UUID,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.RemovePayMethod",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("MethodId", methodId),
			)
			return self.wrapped.RemovePayMethod(token, methodId)
		},
		cmnerrors.Internal,
	)
}

type rentPaymentService struct {
	wrapped payment.IRentPaymentService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRentPayment(
	wrapped payment.IRentPaymentService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) payment.IRentPaymentService {
	return &rentPaymentService{wrapped, hl, tracer}
}

func (self *rentPaymentService) GetPaymentsByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (collection.Collection[payment.Payment], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPaymentsByInstance",
		func(span trace.Span) (collection.Collection[payment.Payment], error) {
			col, err := self.wrapped.GetPaymentsByInstance(token, instanceId)
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *rentPaymentService) GetPaymentsByRentId(
	token token.Token,
	rentId uuid.UUID,
) (collection.Collection[payment.Payment], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPaymentsByRentId",
		func(span trace.Span) (collection.Collection[payment.Payment], error) {
			col, err := self.wrapped.GetPaymentsByRentId(token, rentId)
			span.SetAttributes(
				attribute.Stringer("RentId", rentId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

