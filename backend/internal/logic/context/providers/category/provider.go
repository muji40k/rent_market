package category

import "rent_service/internal/logic/services/interfaces/category"

type IProvider interface {
	GetCategoryService() category.IService
}

