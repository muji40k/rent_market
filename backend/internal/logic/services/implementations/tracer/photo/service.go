package photo

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/photo"
	"rent_service/internal/logic/services/types/token"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/misc/contextholder"
)

type service struct {
	wrapped photo.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped photo.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) photo.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) CreateTempPhoto(
	token token.Token,
	photo photo.Description,
) (uuid.UUID, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.CreateTempPhoto",
		func(span trace.Span) (uuid.UUID, error) {
			span.SetAttributes(
				wrap.AttributeJSON("Photo", photo),
			)
			return self.wrapped.CreateTempPhoto(token, photo)
		},
		cmnerrors.Internal,
	)
}

func (self *service) UploadTempPhoto(
	token token.Token,
	photoId uuid.UUID,
	content []byte,
) error {
	return wrap.SpanCall(self.hl, self.tracer, "Service.UploadTempPhoto",
		func(span trace.Span) error {
			span.SetAttributes(
				attribute.Stringer("PhotoId", photoId),
			)
			return self.wrapped.UploadTempPhoto(token, photoId, content)
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetTempPhoto(
	token token.Token,
	photoId uuid.UUID,
) (photo.TempPhoto, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetTempPhoto",
		func(span trace.Span) (photo.TempPhoto, error) {
			span.SetAttributes(
				attribute.Stringer("PhotoId", photoId),
			)
			return self.wrapped.GetTempPhoto(token, photoId)
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetPhoto(photoId uuid.UUID) (photo.Photo, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetPhoto",
		func(span trace.Span) (photo.Photo, error) {
			span.SetAttributes(
				attribute.Stringer("PhotoId", photoId),
			)
			return self.wrapped.GetPhoto(photoId)
		},
		cmnerrors.Internal,
	)
}

