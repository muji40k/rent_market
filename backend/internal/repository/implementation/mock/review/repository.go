// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/review/repository.go
//

// Package mock_review is a generated GoMock package.
package mock_review

import (
	reflect "reflect"
	models "rent_service/internal/domain/models"
	collection "rent_service/internal/misc/types/collection"
	review "rent_service/internal/repository/interfaces/review"

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

// Create mocks base method.
func (m *MockIRepository) Create(review models.Review) (models.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", review)
	ret0, _ := ret[0].(models.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRepositoryMockRecorder) Create(review any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRepository)(nil).Create), review)
}

// GetWithFilter mocks base method.
func (m *MockIRepository) GetWithFilter(filter review.Filter, sort review.Sort) (collection.Collection[models.Review], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithFilter", filter, sort)
	ret0, _ := ret[0].(collection.Collection[models.Review])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithFilter indicates an expected call of GetWithFilter.
func (mr *MockIRepositoryMockRecorder) GetWithFilter(filter, sort any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithFilter", reflect.TypeOf((*MockIRepository)(nil).GetWithFilter), filter, sort)
}