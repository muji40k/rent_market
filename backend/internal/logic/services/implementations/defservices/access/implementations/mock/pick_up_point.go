// Code generated by MockGen. DO NOT EDIT.
// Source: pick_up_point.go
//
// Generated by this command:
//
//	mockgen -source=pick_up_point.go -destination=implementations/mock/pick_up_point.go
//

// Package mock_access is a generated GoMock package.
package mock_access

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockIPickUpPoint is a mock of IPickUpPoint interface.
type MockIPickUpPoint struct {
	ctrl     *gomock.Controller
	recorder *MockIPickUpPointMockRecorder
	isgomock struct{}
}

// MockIPickUpPointMockRecorder is the mock recorder for MockIPickUpPoint.
type MockIPickUpPointMockRecorder struct {
	mock *MockIPickUpPoint
}

// NewMockIPickUpPoint creates a new mock instance.
func NewMockIPickUpPoint(ctrl *gomock.Controller) *MockIPickUpPoint {
	mock := &MockIPickUpPoint{ctrl: ctrl}
	mock.recorder = &MockIPickUpPointMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPickUpPoint) EXPECT() *MockIPickUpPointMockRecorder {
	return m.recorder
}

// Access mocks base method.
func (m *MockIPickUpPoint) Access(userId, pickUpPointId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Access", userId, pickUpPointId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Access indicates an expected call of Access.
func (mr *MockIPickUpPointMockRecorder) Access(userId, pickUpPointId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Access", reflect.TypeOf((*MockIPickUpPoint)(nil).Access), userId, pickUpPointId)
}