package delivery

import "rent_service/internal/repository/interfaces/delivery"

type IProvider interface {
	GetDeliveryRepository() delivery.IRepository
}

type ICompanyProvider interface {
	GetDeliveryCompanyRepository() delivery.ICompanyRepository
}

