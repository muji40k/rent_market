package pickuppoint

import (
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/pickuppoint"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped pickuppoint.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped pickuppoint.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListPickUpPoints() (collection.Collection[pickuppoint.PickUpPoint], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListPickUpPoints",
		func(span trace.Span) (collection.Collection[pickuppoint.PickUpPoint], error) {
			col, err := self.wrapped.ListPickUpPoints()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetPickUpPointById(
	pickUpPointId uuid.UUID,
) (pickuppoint.PickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPickUpPointById",
		func(span trace.Span) (pickuppoint.PickUpPoint, error) {
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return self.wrapped.GetPickUpPointById(pickUpPointId)
		},
		cmnerrors.Internal,
	)
}

type photoService struct {
	wrapped pickuppoint.IPhotoService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped pickuppoint.IPhotoService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IPhotoService {
	return &photoService{wrapped, hl, tracer}
}

func (self *photoService) ListPickUpPointPhotos(
	pickUpPointId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListPickUpPointPhotos",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.ListPickUpPointPhotos(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

type workingHoursService struct {
	wrapped pickuppoint.IWorkingHoursService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewWorkingHours(
	wrapped pickuppoint.IWorkingHoursService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IWorkingHoursService {
	return &workingHoursService{wrapped, hl, tracer}
}

func (self *workingHoursService) ListPickUpPointWorkingHours(
	pickUpPointId uuid.UUID,
) (collection.Collection[pickuppoint.WorkingHours], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListPickUpPointWorkingHours",
		func(span trace.Span) (collection.Collection[pickuppoint.WorkingHours], error) {
			col, err := self.wrapped.ListPickUpPointWorkingHours(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

