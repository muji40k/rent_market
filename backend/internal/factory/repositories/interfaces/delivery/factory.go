package delivery

import "rent_service/internal/repository/interfaces/delivery"

type IFactory interface {
	CreateDeliveryRepository() delivery.IRepository
}

type ICompanyFactory interface {
	CreateDeliveryCompanyRepository() delivery.ICompanyRepository
}

