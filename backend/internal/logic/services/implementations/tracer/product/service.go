package product

import (
	"fmt"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/product"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	wrapped product.IService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped product.IService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.IService {
	return &service{wrapped, hl, tracer}
}

func (self *service) ListProducts(
	filter product.Filter,
	sort product.Sort,
) (collection.Collection[product.Product], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProducts",
		func(span trace.Span) (collection.Collection[product.Product], error) {
			col, err := self.wrapped.ListProducts(filter, sort)
			span.SetAttributes(
				wrap.AttributeJSON("Filter", filter),
				attribute.String("Sort", func() string {
					switch sort {
					case product.SORT_NONE:
						return "NONE"
					case product.SORT_OFFERS_ASC:
						return "OFFERS_ASC"
					case product.SORT_OFFERS_DSC:
						return "OFFERS_DSC"
					default:
						return "UNKNOWN: " + fmt.Sprint(sort)
					}
				}()),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

func (self *service) GetProductById(
	productId uuid.UUID,
) (product.Product, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.GetProductById",
		func(span trace.Span) (product.Product, error) {
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return self.wrapped.GetProductById(productId)
		},
		cmnerrors.Internal,
	)
}

type characteristicsService struct {
	wrapped product.ICharacteristicsService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewCharacteristics(
	wrapped product.ICharacteristicsService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.ICharacteristicsService {
	return &characteristicsService{wrapped, hl, tracer}
}

func (self *characteristicsService) ListProductCharacteristics(
	productId uuid.UUID,
) (collection.Collection[product.Charachteristic], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProductCharacteristics",
		func(span trace.Span) (collection.Collection[product.Charachteristic], error) {
			col, err := self.wrapped.ListProductCharacteristics(productId)
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

type photoService struct {
	wrapped product.IPhotoService
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped product.IPhotoService,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.IPhotoService {
	return &photoService{wrapped, hl, tracer}
}

func (self *photoService) ListProductPhotos(
	productId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Service.ListProductPhotos",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.ListProductPhotos(productId)
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		cmnerrors.Internal,
	)
}

