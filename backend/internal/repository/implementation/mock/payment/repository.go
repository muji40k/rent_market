// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/payment/repository.go
//

// Package mock_payment is a generated GoMock package.
package mock_payment

import (
	reflect "reflect"
	models "rent_service/internal/domain/models"
	collection "rent_service/internal/misc/types/collection"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockIRepository is a mock of IRepository interface.
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
	isgomock struct{}
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository.
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance.
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// GetByInstanceId mocks base method.
func (m *MockIRepository) GetByInstanceId(instanceId uuid.UUID) (collection.Collection[models.Payment], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByInstanceId", instanceId)
	ret0, _ := ret[0].(collection.Collection[models.Payment])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByInstanceId indicates an expected call of GetByInstanceId.
func (mr *MockIRepositoryMockRecorder) GetByInstanceId(instanceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByInstanceId", reflect.TypeOf((*MockIRepository)(nil).GetByInstanceId), instanceId)
}

// GetByRentId mocks base method.
func (m *MockIRepository) GetByRentId(rentId uuid.UUID) (collection.Collection[models.Payment], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRentId", rentId)
	ret0, _ := ret[0].(collection.Collection[models.Payment])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRentId indicates an expected call of GetByRentId.
func (mr *MockIRepositoryMockRecorder) GetByRentId(rentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRentId", reflect.TypeOf((*MockIRepository)(nil).GetByRentId), rentId)
}
