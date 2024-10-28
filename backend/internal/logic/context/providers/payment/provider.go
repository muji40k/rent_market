package payment

import "rent_service/internal/logic/services/interfaces/payment"

type IPayMethodProvider interface {
	GetPayMethodService() payment.IPayMethodService
}

type IUserPayMethodProvider interface {
	GetUserPayMethodService() payment.IUserPayMethodService
}

type IRentPaymentProvider interface {
	GetRentPaymentService() payment.IRentPaymentService
}

