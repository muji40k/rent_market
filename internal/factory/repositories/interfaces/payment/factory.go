package payment

import "rent_service/internal/repository/interfaces/payment"

type IFactory interface {
	CreatePaymentRepository() payment.IRepository
}

