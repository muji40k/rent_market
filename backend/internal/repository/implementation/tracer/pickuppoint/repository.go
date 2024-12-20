package pickuppoint

import (
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/pickuppoint"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped pickuppoint.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped pickuppoint.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetById(
	pickUpPointId uuid.UUID,
) (models.PickUpPoint, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.PickUpPoint.GetById",
		func(span trace.Span) (models.PickUpPoint, error) {
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return self.wrapped.GetById(pickUpPointId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetAll() (collection.Collection[models.PickUpPoint], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.PickUpPoint.GetAll",
		func(span trace.Span) (collection.Collection[models.PickUpPoint], error) {
			col, err := self.wrapped.GetAll()
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

type photoRepository struct {
	wrapped pickuppoint.IPhotoRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped pickuppoint.IPhotoRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IPhotoRepository {
	return &photoRepository{wrapped, hl, tracer}
}

func (self *photoRepository) GetById(
	pickUpPointId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.PickUpPointPhoto.GetById",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.GetById(pickUpPointId)
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

type workingHoursRepository struct {
	wrapped pickuppoint.IWorkingHoursRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewWorkingHours(
	wrapped pickuppoint.IWorkingHoursRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) pickuppoint.IWorkingHoursRepository {
	return &workingHoursRepository{wrapped, hl, tracer}
}

func (self *workingHoursRepository) GetById(
	pickUpPointId uuid.UUID,
) (models.PickUpPointWorkingHours, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.PickUpPointWorkingHours.GetById",
		func(span trace.Span) (models.PickUpPointWorkingHours, error) {
			span.SetAttributes(
				attribute.Stringer("PickUpPointId", pickUpPointId),
			)
			return self.wrapped.GetById(pickUpPointId)
		},
		func(err error) error { return err },
	)
}

