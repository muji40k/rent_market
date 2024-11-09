// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/photo/repository.go
//

// Package mock_photo is a generated GoMock package.
package mock_photo

import (
	reflect "reflect"
	models "rent_service/internal/domain/models"

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

// Create mocks base method.
func (m *MockIRepository) Create(photo models.Photo) (models.Photo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", photo)
	ret0, _ := ret[0].(models.Photo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRepositoryMockRecorder) Create(photo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRepository)(nil).Create), photo)
}

// GetById mocks base method.
func (m *MockIRepository) GetById(photoId uuid.UUID) (models.Photo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", photoId)
	ret0, _ := ret[0].(models.Photo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRepositoryMockRecorder) GetById(photoId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRepository)(nil).GetById), photoId)
}

// MockITempRepository is a mock of ITempRepository interface.
type MockITempRepository struct {
	ctrl     *gomock.Controller
	recorder *MockITempRepositoryMockRecorder
	isgomock struct{}
}

// MockITempRepositoryMockRecorder is the mock recorder for MockITempRepository.
type MockITempRepositoryMockRecorder struct {
	mock *MockITempRepository
}

// NewMockITempRepository creates a new mock instance.
func NewMockITempRepository(ctrl *gomock.Controller) *MockITempRepository {
	mock := &MockITempRepository{ctrl: ctrl}
	mock.recorder = &MockITempRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITempRepository) EXPECT() *MockITempRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockITempRepository) Create(photo models.TempPhoto) (models.TempPhoto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", photo)
	ret0, _ := ret[0].(models.TempPhoto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockITempRepositoryMockRecorder) Create(photo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockITempRepository)(nil).Create), photo)
}

// GetById mocks base method.
func (m *MockITempRepository) GetById(photoId uuid.UUID) (models.TempPhoto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", photoId)
	ret0, _ := ret[0].(models.TempPhoto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockITempRepositoryMockRecorder) GetById(photoId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockITempRepository)(nil).GetById), photoId)
}

// Remove mocks base method.
func (m *MockITempRepository) Remove(photoId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", photoId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockITempRepositoryMockRecorder) Remove(photoId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockITempRepository)(nil).Remove), photoId)
}

// Update mocks base method.
func (m *MockITempRepository) Update(photo models.TempPhoto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", photo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockITempRepositoryMockRecorder) Update(photo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockITempRepository)(nil).Update), photo)
}