// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/category/repository.go
//

// Package mock_category is a generated GoMock package.
package mock_category

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

// GetAll mocks base method.
func (m *MockIRepository) GetAll() (collection.Collection[models.Category], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(collection.Collection[models.Category])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIRepository)(nil).GetAll))
}

// GetPath mocks base method.
func (m *MockIRepository) GetPath(leaf uuid.UUID) (collection.Collection[models.Category], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPath", leaf)
	ret0, _ := ret[0].(collection.Collection[models.Category])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPath indicates an expected call of GetPath.
func (mr *MockIRepositoryMockRecorder) GetPath(leaf any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPath", reflect.TypeOf((*MockIRepository)(nil).GetPath), leaf)
}
