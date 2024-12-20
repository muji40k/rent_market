package photo

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/repository/interfaces/photo"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped photo.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped photo.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) photo.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) Create(photo models.Photo) (models.Photo, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Photo.Create",
		func(span trace.Span) (models.Photo, error) {
			out, err := self.wrapped.Create(photo)
			span.SetAttributes(
				wrap.AttributeJSON("Photo", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetById(photoId uuid.UUID) (models.Photo, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Photo.GetById",
		func(span trace.Span) (models.Photo, error) {
			span.SetAttributes(
				attribute.Stringer("PhotoId", photoId),
			)
			return self.wrapped.GetById(photoId)
		},
		func(err error) error { return err },
	)
}

type tempRepository struct {
	wrapped photo.ITempRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewTemp(
	wrapped photo.ITempRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) photo.ITempRepository {
	return &tempRepository{wrapped, hl, tracer}
}

func (self *tempRepository) Create(
	photo models.TempPhoto,
) (models.TempPhoto, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.TempPhoto.Create",
		func(span trace.Span) (models.TempPhoto, error) {
			out, err := self.wrapped.Create(photo)
			span.SetAttributes(
				wrap.AttributeJSON("Photo", out),
			)
			return out, err
		},
		func(err error) error { return err },
	)
}

func (self *tempRepository) Update(photo models.TempPhoto) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.TempPhoto.Update",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("TempPhoto", photo),
			)
			return self.wrapped.Update(photo)
		},
		func(err error) error { return err },
	)
}

func (self *tempRepository) GetById(
	photoId uuid.UUID,
) (models.TempPhoto, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.TempPhoto.GetById",
		func(span trace.Span) (models.TempPhoto, error) {
			span.SetAttributes(
				attribute.Stringer("TempPhotoId", photoId),
			)
			return self.wrapped.GetById(photoId)
		},
		func(err error) error { return err },
	)
}

func (self *tempRepository) Remove(photoId uuid.UUID) error {
	return wrap.SpanCall(self.hl, self.tracer, "Repository.TempPhoto.Remove",
		func(span trace.Span) error {
			span.SetAttributes(
				wrap.AttributeJSON("TempPhoto", photoId),
			)
			return self.wrapped.Remove(photoId)
		},
		func(err error) error { return err },
	)
}

