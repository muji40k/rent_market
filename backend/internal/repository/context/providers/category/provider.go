package category

import "rent_service/internal/repository/interfaces/category"

type IProvider interface {
	GetCategoryRepository() category.IRepository
}

