package paymethod

import "rent_service/internal/repository/interfaces/paymethod"

type IProvider interface {
	GetPayMethodProvider() paymethod.IRepository
}

