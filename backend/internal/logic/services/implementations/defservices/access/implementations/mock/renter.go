// Code generated by MockGen. DO NOT EDIT.
// Source: renter.go
//
// Generated by this command:
//
//	mockgen -source=renter.go -destination=implementations/mock/renter.go
//

// Package mock_access is a generated GoMock package.
package mock_access

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockIRenter is a mock of IRenter interface.
type MockIRenter struct {
	ctrl     *gomock.Controller
	recorder *MockIRenterMockRecorder
	isgomock struct{}
}

// MockIRenterMockRecorder is the mock recorder for MockIRenter.
type MockIRenterMockRecorder struct {
	mock *MockIRenter
}

// NewMockIRenter creates a new mock instance.
func NewMockIRenter(ctrl *gomock.Controller) *MockIRenter {
	mock := &MockIRenter{ctrl: ctrl}
	mock.recorder = &MockIRenterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRenter) EXPECT() *MockIRenterMockRecorder {
	return m.recorder
}

// Access mocks base method.
func (m *MockIRenter) Access(userId, renterUserId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Access", userId, renterUserId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Access indicates an expected call of Access.
func (mr *MockIRenterMockRecorder) Access(userId, renterUserId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Access", reflect.TypeOf((*MockIRenter)(nil).Access), userId, renterUserId)
}