package delivery

import "rent_service/internal/logic/services/interfaces/delivery"

type IFactory interface {
	CreateDeliveryService() delivery.IService
}

type ICompanyFactory interface {
	CreateDeliveryCompanyService() delivery.ICompanyService
}

