package rent

import (
	"rent_service/internal/domain/records"
	"rent_service/internal/domain/requests"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/rent"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped rent.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped rent.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(rent records.Rent) (records.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Rent.Create",
		func(span trace.Span) (records.Rent, error) {
			out, err := self.wrapped.Create(rent)
			span.SetAttributes(
				wrap.AttributeJSON("Rent", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) Update(rent records.Rent) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.Rent.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("Rent", rent),
			)
			return self.wrapped.Update(rent)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(rentId uuid.UUID) (records.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Rent.GetById",
		func(span trace.Span) (records.Rent, error) {
			span.SetAttributes(
				attribute.Stringer("RentId", rentId),
			)
			return self.wrapped.GetById(rentId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[records.Rent], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Rent.GetByUserId",
		func(span trace.Span) (collection.Collection[records.Rent], error) {
			col, err := self.wrapped.GetByUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetActiveByInstanceId(
	instanceId uuid.UUID,
) (records.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Rent.GetActiveByInstanceId",
		func(span trace.Span) (records.Rent, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetActiveByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetPastByUserId(
	userId uuid.UUID,
) (collection.Collection[records.Rent], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Rent.GetPastByUserId",
		func(span trace.Span) (collection.Collection[records.Rent], error) {
			col, err := self.wrapped.GetPastByUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

type requestRepository struct {
	wrapped rent.IRequestRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewRequest(
	wrapped rent.IRequestRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IRequestRepository {
	return &requestRepository{wrapped, hl, tracer}
}

func (self *requestRepository) Create(
	request requests.Rent,
) (requests.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentRequest.Create",
		func(span trace.Span) (requests.Rent, error) {
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
) (requests.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentRequest.GetById",
		func(span trace.Span) (requests.Rent, error) {
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
) (collection.Collection[requests.Rent], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentRequest.GetByUserId",
		func(span trace.Span) (collection.Collection[requests.Rent], error) {
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
) (requests.Rent, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentRequest.GetByInstanceId",
		func(span trace.Span) (requests.Rent, error) {
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
) (collection.Collection[requests.Rent], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentRequest.GetByPickUpPointId",
		func(span trace.Span) (collection.Collection[requests.Rent], error) {
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
	return wrap.SpanCall(self.hl, self.tracer, "Repository.RentRequest.Remove",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("RequestId", requestId),
			)
			return self.wrapped.Remove(requestId)
		},
		func(err error) error { return err },
	)
}

type returnRepository struct {
	wrapped rent.IReturnRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewReturn(
	wrapped rent.IReturnRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) rent.IReturnRepository {
	return &returnRepository{wrapped, hl, tracer}
}

func (self *returnRepository) Create(
	request requests.Return,
) (requests.Return, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentReturn.Create",
		func(span trace.Span) (requests.Return, error) {
			out, err := self.wrapped.Create(request)
			span.SetAttributes(
				wrap.AttributeJSON("Return", request),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *returnRepository) GetById(
	requestId uuid.UUID,
) (requests.Return, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentReturn.GetById",
		func(span trace.Span) (requests.Return, error) {
			span.SetAttributes(
				attribute.Stringer("ReturnId", requestId),
			)
			return self.wrapped.GetById(requestId)
		},
		func(err error) error { return err },
	)
}

func (self *returnRepository) GetByUserId(
	userId uuid.UUID,
) (collection.Collection[requests.Return], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentReturn.GetByUserId",
		func(span trace.Span) (collection.Collection[requests.Return], error) {
			col, err := self.wrapped.GetByUserId(userId)
			span.SetAttributes(
				attribute.Stringer("UserId", userId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *returnRepository) GetByInstanceId(
	instanceId uuid.UUID,
) (requests.Return, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentReturn.GetByInstanceId",
		func(span trace.Span) (requests.Return, error) {
			span.SetAttributes(
				attribute.Stringer("InstanceId", instanceId),
			)
			return self.wrapped.GetByInstanceId(instanceId)
		},
		func(err error) error { return err },
	)
}

func (self *returnRepository) GetByPickUpPointId(
	pickUpPointId uuid.UUID,
) (collection.Collection[requests.Return], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.RentReturn.GetByPickUpPointId",
		func(span trace.Span) (collection.Collection[requests.Return], error) {
			col, err := self.wrapped.GetByPickUpPointId(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

func (self *returnRepository) Remove(requestId uuid.UUID) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.RentReturn.Remove",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("ReturnId", requestId),
			)
			return self.wrapped.Remove(requestId)
		},
		func(err error) error { return err },
	)
}

