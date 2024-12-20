package product

import (
	"fmt"
	"rent_service/internal/domain/models"
	tcollection "rent_service/internal/misc/tracer/collection"
	"rent_service/internal/misc/tracer/wrap"
	"rent_service/internal/misc/types/collection"
	"rent_service/internal/repository/interfaces/product"
	"rent_service/misc/contextholder"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type repository struct {
	wrapped product.IRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func New(
	wrapped product.IRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.IRepository {
	return &repository{wrapped, hl, tracer}
}

func (self *repository) GetById(productId uuid.UUID) (models.Product, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Product.GetById",
		func(span trace.Span) (models.Product, error) {
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return self.wrapped.GetById(productId)
		},
		func(err error) error { return err },
	)
}

func (self *repository) GetWithFilter(
	filter product.Filter,
	sort product.Sort,
) (collection.Collection[models.Product], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.Product.GetWithFilter",
		func(span trace.Span) (collection.Collection[models.Product], error) {
			col, err := self.wrapped.GetWithFilter(filter, sort)
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
		func(err error) error { return err },
	)
}

type characteristicsRepository struct {
	wrapped product.ICharacteristicsRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewCharacteristics(
	wrapped product.ICharacteristicsRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.ICharacteristicsRepository {
	return &characteristicsRepository{wrapped, hl, tracer}
}

func (self *characteristicsRepository) GetByProductId(
	productId uuid.UUID,
) (models.ProductCharacteristics, error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProductCharacteristics.GetByProductId",
		func(span trace.Span) (models.ProductCharacteristics, error) {
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return self.wrapped.GetByProductId(productId)
		},
		func(err error) error { return err },
	)
}

type photoRepository struct {
	wrapped product.IPhotoRepository
	hl      *contextholder.Holder
	tracer  trace.Tracer
}

func NewPhoto(
	wrapped product.IPhotoRepository,
	hl *contextholder.Holder,
	tracer trace.Tracer,
) product.IPhotoRepository {
	return &photoRepository{wrapped, hl, tracer}
}

func (self *photoRepository) GetByProductId(
	productId uuid.UUID,
) (collection.Collection[uuid.UUID], error) {
	return wrap.SpanCallValue(self.hl, self.tracer, "Repository.ProductPhoto.GetByProductId",
		func(span trace.Span) (collection.Collection[uuid.UUID], error) {
			col, err := self.wrapped.GetByProductId(productId)
			span.SetAttributes(
				attribute.Stringer("ProductId", productId),
			)
			return tcollection.TraceCollection(self.hl, self.tracer, col), err
		},
		func(err error) error { return err },
	)
}

