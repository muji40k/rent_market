package product

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/product"
	. "rent_service/internal/misc/types/collection"
	product_provider "rent_service/internal/repository/context/providers/product"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
	product_repository "rent_service/internal/repository/interfaces/product"

	"github.com/google/uuid"
)

type repoproviders struct {
	product product_provider.IProvider
}

type service struct {
	repos repoproviders
}

func New(product product_provider.IProvider) product.IService {
	return &service{repoproviders{product}}
}

func mapf(value *models.Product) product.Product {
	return product.Product{
		Id:          value.Id,
		Name:        value.Name,
		CategoryId:  value.CategoryId,
		Description: value.Description,
	}
}

func mapSort(value *product.Sort) (product_repository.Sort, error) {
	switch *value {
	case product.SORT_NONE:
		return product_repository.SORT_NONE, nil
	case product.SORT_OFFERS_ASC:
		return product_repository.SORT_OFFERS_ASC, nil
	case product.SORT_OFFERS_DSC:
		return product_repository.SORT_OFFERS_DSC, nil
	default:
		return product_repository.SORT_NONE, cmnerrors.Unknown("sort")
	}
}

func mapFilter(value *product.Filter) (product_repository.Filter, error) {
	var empty product_repository.Filter
	var out = product_repository.Filter{
		CategoryId: value.CategoryId,
	}

	if nil != value.Query {
		out.Query = new(string)
		*out.Query = *value.Query
	}

	for _, f := range value.Characteristics {
		if "" == f.Key {
			return empty, cmnerrors.Incorrect("characteristic.key")
		}

		if nil != f.Range && nil == f.Values {
			out.Ranges = append(out.Ranges, product_repository.Range{
				Characteristic: product_repository.Characteristic{Key: f.Key},
				Min:            f.Range.Min,
				Max:            f.Range.Max,
			})
		} else if nil == f.Range && nil != f.Values {
			selector := product_repository.Selector{
				Characteristic: product_repository.Characteristic{Key: f.Key},
				Values:         make([]string, len(f.Values)),
			}

			for i, v := range f.Values {
				if "" == v {
					return empty, cmnerrors.Incorrect("characteristic.value")
				}

				selector.Values[i] = v
			}

			out.Selectors = append(out.Selectors, selector)
		} else {
			return empty, cmnerrors.Incorrect("characteristic.content")
		}
	}

	return out, nil
}

func (self *service) ListProducts(
	filter product.Filter,
	sort product.Sort,
) (Collection[product.Product], error) {
	var products Collection[models.Product]
	var sortr product_repository.Sort
	filterr, err := mapFilter(&filter)

	if nil == err {
		sortr, err = mapSort(&sort)
	}

	if nil == err {
		repo := self.repos.product.GetProductRepository()
		products, err = repo.GetWithFilter(filterr, sortr)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return MapCollection(mapf, products), err
}

func (self *service) GetProductById(
	productId uuid.UUID,
) (product.Product, error) {
	repo := self.repos.product.GetProductRepository()
	product, err := repo.GetById(productId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return mapf(&product), err
}

type characteristicsRepoProviders struct {
	characteristics product_provider.ICharacteristicsProvider
}

type characteristicsService struct {
	repos characteristicsRepoProviders
}

func NewCharacteristics(
	characteristics product_provider.ICharacteristicsProvider,
) product.ICharacteristicsService {
	return &characteristicsService{
		characteristicsRepoProviders{characteristics},
	}
}

func (self *characteristicsService) ListProductCharacteristics(
	productId uuid.UUID,
) (Collection[product.Charachteristic], error) {
	var chars Collection[product.Charachteristic]
	repo := self.repos.characteristics.GetProductCharacteristicsRepository()
	productChars, err := repo.GetByProductId(productId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err {
		buf := make([]product.Charachteristic, len(productChars.Map))
		i := 0

		for _, char := range productChars.Map {
			buf[i] = product.Charachteristic{
				Id:    char.Id,
				Name:  char.Name,
				Value: char.Value,
			}
			i++
		}

		chars = SliceCollection(buf)
	}

	return chars, err
}

type photoRepoProviders struct {
	photo product_provider.IPhotoProvider
}

type photoService struct {
	repos photoRepoProviders
}

func NewPhoto(
	photo product_provider.IPhotoProvider,
) product.IPhotoService {
	return &photoService{photoRepoProviders{photo}}
}

func (self *photoService) ListProductPhotos(
	productId uuid.UUID,
) (Collection[uuid.UUID], error) {
	repo := self.repos.photo.GetProductPhotoRepository()
	photos, err := repo.GetByProductId(productId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return photos, err
}

