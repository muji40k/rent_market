package delivery

import "rent_service/internal/logic/services/interfaces/delivery"

type IProvider interface {
	GetDeliveryService() delivery.IService
}

type ICompanyProvider interface {
	GetDeliveryCompanyService() delivery.ICompanyService
}

