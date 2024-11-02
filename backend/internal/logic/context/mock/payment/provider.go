package payment

import "rent_service/internal/logic/services/interfaces/payment"

type MockPayMethodProvider struct {
	service payment.IPayMethodService
}

func NewPayMethod(service payment.IPayMethodService) *MockPayMethodProvider {
	return &MockPayMethodProvider{service}
}

func (self *MockPayMethodProvider) GetPayMethodService() payment.IPayMethodService {
	return self.service
}

type MockUserPayMethodProvider struct {
	service payment.IUserPayMethodService
}

func NewUserPayMethod(service payment.IUserPayMethodService) *MockUserPayMethodProvider {
	return &MockUserPayMethodProvider{service}
}

func (self *MockUserPayMethodProvider) GetUserPayMethodService() payment.IUserPayMethodService {
	return self.service
}

type MockRentPaymentProvider struct {
	service payment.IRentPaymentService
}

func NewRentPayment(service payment.IRentPaymentService) *MockRentPaymentProvider {
	return &MockRentPaymentProvider{service}
}

func (self *MockRentPaymentProvider) GetRentPaymentService() payment.IRentPaymentService {
	return self.service
}

