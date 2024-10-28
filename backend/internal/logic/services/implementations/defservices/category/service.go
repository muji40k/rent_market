package category

import (
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/interfaces/category"
	. "rent_service/internal/misc/types/collection"
	category_provider "rent_service/internal/repository/context/providers/category"

	"github.com/google/uuid"
)

type repoproviders struct {
	category category_provider.IProvider
}

type service struct {
	repos repoproviders
}

func New(category category_provider.IProvider) category.IService {
	return &service{repoproviders{category}}
}

func mapf(value *models.Category) category.Category {
	return category.Category{
		Id:       value.Id,
		ParentId: value.ParentId,
		Name:     value.Name,
	}
}

func (self *service) ListCategories() (Collection[category.Category], error) {
	var categories Collection[category.Category]
	repo := self.repos.category.GetCategoryRepository()
	col, err := repo.GetAll()

	if nil == err {
		categories = MapCollection(mapf, col)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return categories, err
}

func (self *service) GetFullCategory(
	categoryId uuid.UUID,
) (Collection[category.Category], error) {
	var path Collection[category.Category]
	repo := self.repos.category.GetCategoryRepository()
	col, err := repo.GetPath(categoryId)

	if nil == err {
		path = MapCollection(mapf, col)
	} else {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return path, err
}

