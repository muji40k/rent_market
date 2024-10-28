package paymethod

import "rent_service/internal/repository/interfaces/paymethod"

type IProvider interface {
	GetPayMethodRepository() paymethod.IRepository
}

