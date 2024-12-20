package delivery

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/domain/requests"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/delivery"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped delivery.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped delivery.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) delivery.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(
	delivery requests.Delivery,
) (requests.Delivery, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Delivery.Create",
		func(span trace.Span) (requests.Delivery, error) {
			out, err := self.wrapped.Create(delivery)
			span.SetAttributes(
				wrap.AttributeJSON("Delivery", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(delivery requests.Delivery) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Delivery.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Delivery", delivery),
			)
			return self.wrapped.Update(delivery)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(
	deliveryId uuid.UUID,
) (requests.Delivery, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Delivery.GetById",
		func(span trace.Span) (requests.Delivery, error) {
			span.SetAttributes(
				attribute.Stringer("DeliveryId", deliveryId),
			)
			return self.wrapped.GetById(deliveryId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetActiveByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Delivery], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Delivery.GetActiveByPickUpPointId",
		func(span trace.Span) (collection.Collection[requests.Delivery], error) {
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
) (requests.Delivery, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Delivery.GetActiveByInstanceId",
		func(span trace.Span) (requests.Delivery, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetActiveByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

type companyRepository struct {
	wrapped delivery.ICompanyRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewCompany(
	wrapped delivery.ICompanyRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) delivery.ICompanyRepository {
	return &companyRepository{wrapped, hl, tracer}
}

func (self *companyRepository) GetById(
	companyId uuid.UUID,
) (models.DeliveryCompany, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.DeliveryCompany.GetById",
		func(span trace.Span) (models.DeliveryCompany, error) {
			span.SetAttributes(
				attribute.Stringer("CompanyId", companyId),
			)
			return self.wrapped.GetById(companyId)
		},
		func(err error) error { return err },
	)
}

func (self *companyRepository) GetAll() (collection.Collection[models.DeliveryCompany], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.DeliveryCompany.GetAll",
		func(span trace.Span) (collection.Collection[models.DeliveryCompany], error) {
			col, err := self.wrapped.GetAll()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

