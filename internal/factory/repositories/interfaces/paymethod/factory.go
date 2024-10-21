package paymethod

import "rent_service/internal/repository/interfaces/paymethod"

type IFactory interface {
	CreatePayMethodRepository() paymethod.IRepository
}

