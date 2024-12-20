package delivery

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/delivery"
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
	wrapped delivery.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped delivery.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) delivery.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListDeliveriesByPickUpPoint(
	token token.Token,
	pickUpPointId uuid.UUID,
) (collection.Collection[delivery.Delivery], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListDeliveriesByPickUpPoint",
		func(span trace.Span) (collection.Collection[delivery.Delivery], error) {
			col, err := self.wrapped.ListDeliveriesByPickUpPoint(token, pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetDeliveryByInstance(
	token token.Token,
	instanceId uuid.UUID,
) (delivery.Delivery, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetDeliveryByInstance",
		func(span trace.Span) (delivery.Delivery, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetDeliveryByInstance(token, instanceId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) CreateDelivery(
	token token.Token,
	form delivery.CreateForm,
) (delivery.Delivery, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateDelivery",
		func(span trace.Span) (delivery.Delivery, error) {
			span.SetAttributes(
				wrap.AttributeJSON("CreateForm", form),
			)
			return self.wrapped.CreateDelivery(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *service) SendDelivery(
	token token.Token,
	form delivery.SendForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.SendDelivery",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("SendForm", form),
			)
			return self.wrapped.SendDelivery(token, form)
		},
		cmnerrors.Internal,
	)
}

func (self *service) AcceptDelivery(
	token token.Token,
	form delivery.AcceptForm,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.AcceptDelivery",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("AcceptForm", form),
			)
			return self.wrapped.AcceptDelivery(token, form)
		},
		cmnerrors.Internal,
	)
}

type companyService struct {
	wrapped delivery.ICompanyService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewCompany(
	wrapped delivery.ICompanyService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) delivery.ICompanyService {
	return &companyService{wrapped, hl, tracer}
}

func (self *companyService) ListDeliveryCompanies(
	token token.Token,
) (collection.Collection[delivery.DeliveryCompany], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListDeliveryCompanies",
		func(_ trace.Span) (collection.Collection[delivery.DeliveryCompany], error) {
			col, err := self.wrapped.ListDeliveryCompanies(token)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *companyService) GetDeliveryCompanyById(
	token token.Token,
	companyId uuid.UUID,
) (delivery.DeliveryCompany, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetDeliveryCompanyById",
		func(span trace.Span) (delivery.DeliveryCompany, error) {
			span.SetAttributes(
				attribute.Stringer("CompanyId", companyId),
			)
			return self.wrapped.GetDeliveryCompanyById(token, companyId)
		},
		cmnerrors.Internal,
	)
}

