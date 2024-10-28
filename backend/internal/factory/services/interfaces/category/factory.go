package category

import "rent_service/internal/logic/services/interfaces/category"

type IFactory interface {
	CreateCategoryService() category.IService
}

