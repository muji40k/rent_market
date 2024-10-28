package category

import "rent_service/internal/repository/interfaces/category"

type IFactory interface {
	CreateCategoryRepository() category.IRepository
}

