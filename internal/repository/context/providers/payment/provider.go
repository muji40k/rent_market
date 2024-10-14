package payment

import "rent_service/internal/repository/interfaces/payment"

type IProvider interface {
	GetPaymentRepository() payment.IRepository
}

