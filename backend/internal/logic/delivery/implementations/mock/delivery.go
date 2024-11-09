// Code generated by MockGen. DO NOT EDIT.
// Source: delivery.go
//
// Generated by this command:
//
//	mockgen -source=delivery.go -destination=implementations/mock/delivery.go
//

// Package mock_delivery is a generated GoMock package.
package mock_delivery

import (
	reflect "reflect"
	models "rent_service/internal/domain/models"
	delivery "rent_service/internal/logic/delivery"

	gomock "go.uber.org/mock/gomock"
)

// MockICreator is a mock of ICreator interface.
type MockICreator struct {
	ctrl     *gomock.Controller
	recorder *MockICreatorMockRecorder
	isgomock struct{}
}

// MockICreatorMockRecorder is the mock recorder for MockICreator.
type MockICreatorMockRecorder struct {
	mock *MockICreator
}

// NewMockICreator creates a new mock instance.
func NewMockICreator(ctrl *gomock.Controller) *MockICreator {
	mock := &MockICreator{ctrl: ctrl}
	mock.recorder = &MockICreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICreator) EXPECT() *MockICreatorMockRecorder {
	return m.recorder
}

// CreateDelivery mocks base method.
func (m *MockICreator) CreateDelivery(from, to models.Address, verificationCode string) (delivery.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDelivery", from, to, verificationCode)
	ret0, _ := ret[0].(delivery.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDelivery indicates an expected call of CreateDelivery.
func (mr *MockICreatorMockRecorder) CreateDelivery(from, to, verificationCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDelivery", reflect.TypeOf((*MockICreator)(nil).CreateDelivery), from, to, verificationCode)
}
