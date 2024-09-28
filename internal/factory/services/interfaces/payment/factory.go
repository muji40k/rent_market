package payment

import "rent_service/internal/logic/services/interfaces/payment"

type IPayMethodFactory interface {
	CreatePayMethodService() payment.IPayMethodService
}

type IUserPayMethodFactory interface {
	CreateUserPayMethodService() payment.IUserPayMethodService
}

type IRentPaymentFactory interface {
	CreateRentPaymentService() payment.IRentPaymentService
}

