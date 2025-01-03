// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/provision/repository.go
//

// Package mock_provision is a generated GoMock package.
package mock_provision

import (
	reflect "reflect"
	records "rent_service/internal/domain/records"
	requests "rent_service/internal/domain/requests"
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

// Create mocks base method.
func (m *MockIRepository) Create(provision records.Provision) (records.Provision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", provision)
	ret0, _ := ret[0].(records.Provision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRepositoryMockRecorder) Create(provision any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRepository)(nil).Create), provision)
}

// GetActiveByInstanceId mocks base method.
func (m *MockIRepository) GetActiveByInstanceId(instanceId uuid.UUID) (records.Provision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveByInstanceId", instanceId)
	ret0, _ := ret[0].(records.Provision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveByInstanceId indicates an expected call of GetActiveByInstanceId.
func (mr *MockIRepositoryMockRecorder) GetActiveByInstanceId(instanceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveByInstanceId", reflect.TypeOf((*MockIRepository)(nil).GetActiveByInstanceId), instanceId)
}

// GetById mocks base method.
func (m *MockIRepository) GetById(provisionId uuid.UUID) (records.Provision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", provisionId)
	ret0, _ := ret[0].(records.Provision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRepositoryMockRecorder) GetById(provisionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRepository)(nil).GetById), provisionId)
}

// GetByInstanceId mocks base method.
func (m *MockIRepository) GetByInstanceId(instanceId uuid.UUID) (collection.Collection[records.Provision], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByInstanceId", instanceId)
	ret0, _ := ret[0].(collection.Collection[records.Provision])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByInstanceId indicates an expected call of GetByInstanceId.
func (mr *MockIRepositoryMockRecorder) GetByInstanceId(instanceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByInstanceId", reflect.TypeOf((*MockIRepository)(nil).GetByInstanceId), instanceId)
}

// GetByRenterUserId mocks base method.
func (m *MockIRepository) GetByRenterUserId(userId uuid.UUID) (collection.Collection[records.Provision], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRenterUserId", userId)
	ret0, _ := ret[0].(collection.Collection[records.Provision])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRenterUserId indicates an expected call of GetByRenterUserId.
func (mr *MockIRepositoryMockRecorder) GetByRenterUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRenterUserId", reflect.TypeOf((*MockIRepository)(nil).GetByRenterUserId), userId)
}

// Update mocks base method.
func (m *MockIRepository) Update(provision records.Provision) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", provision)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIRepositoryMockRecorder) Update(provision any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRepository)(nil).Update), provision)
}

// MockIRequestRepository is a mock of IRequestRepository interface.
type MockIRequestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRequestRepositoryMockRecorder
	isgomock struct{}
}

// MockIRequestRepositoryMockRecorder is the mock recorder for MockIRequestRepository.
type MockIRequestRepositoryMockRecorder struct {
	mock *MockIRequestRepository
}

// NewMockIRequestRepository creates a new mock instance.
func NewMockIRequestRepository(ctrl *gomock.Controller) *MockIRequestRepository {
	mock := &MockIRequestRepository{ctrl: ctrl}
	mock.recorder = &MockIRequestRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRequestRepository) EXPECT() *MockIRequestRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIRequestRepository) Create(request requests.Provide) (requests.Provide, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", request)
	ret0, _ := ret[0].(requests.Provide)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRequestRepositoryMockRecorder) Create(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRequestRepository)(nil).Create), request)
}

// GetById mocks base method.
func (m *MockIRequestRepository) GetById(requestId uuid.UUID) (requests.Provide, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", requestId)
	ret0, _ := ret[0].(requests.Provide)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRequestRepositoryMockRecorder) GetById(requestId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRequestRepository)(nil).GetById), requestId)
}

// GetByInstanceId mocks base method.
func (m *MockIRequestRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Provide, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByInstanceId", instanceId)
	ret0, _ := ret[0].(requests.Provide)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByInstanceId indicates an expected call of GetByInstanceId.
func (mr *MockIRequestRepositoryMockRecorder) GetByInstanceId(instanceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByInstanceId", reflect.TypeOf((*MockIRequestRepository)(nil).GetByInstanceId), instanceId)
}

// GetByPickUpPointId mocks base method.
func (m *MockIRequestRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Provide], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPickUpPointId", pickUpPointId)
	ret0, _ := ret[0].(collection.Collection[requests.Provide])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPickUpPointId indicates an expected call of GetByPickUpPointId.
func (mr *MockIRequestRepositoryMockRecorder) GetByPickUpPointId(pickUpPointId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPickUpPointId", reflect.TypeOf((*MockIRequestRepository)(nil).GetByPickUpPointId), pickUpPointId)
}

// GetByUserId mocks base method.
func (m *MockIRequestRepository) GetByUserId(userId uuid.UUID) (collection.Collection[requests.Provide], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].(collection.Collection[requests.Provide])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockIRequestRepositoryMockRecorder) GetByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockIRequestRepository)(nil).GetByUserId), userId)
}

// Remove mocks base method.
func (m *MockIRequestRepository) Remove(requestId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", requestId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockIRequestRepositoryMockRecorder) Remove(requestId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockIRequestRepository)(nil).Remove), requestId)
}

// MockIRevokeRepository is a mock of IRevokeRepository interface.
type MockIRevokeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRevokeRepositoryMockRecorder
	isgomock struct{}
}

// MockIRevokeRepositoryMockRecorder is the mock recorder for MockIRevokeRepository.
type MockIRevokeRepositoryMockRecorder struct {
	mock *MockIRevokeRepository
}

// NewMockIRevokeRepository creates a new mock instance.
func NewMockIRevokeRepository(ctrl *gomock.Controller) *MockIRevokeRepository {
	mock := &MockIRevokeRepository{ctrl: ctrl}
	mock.recorder = &MockIRevokeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRevokeRepository) EXPECT() *MockIRevokeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIRevokeRepository) Create(request requests.Revoke) (requests.Revoke, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", request)
	ret0, _ := ret[0].(requests.Revoke)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRevokeRepositoryMockRecorder) Create(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRevokeRepository)(nil).Create), request)
}

// GetById mocks base method.
func (m *MockIRevokeRepository) GetById(requestId uuid.UUID) (requests.Revoke, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", requestId)
	ret0, _ := ret[0].(requests.Revoke)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRevokeRepositoryMockRecorder) GetById(requestId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRevokeRepository)(nil).GetById), requestId)
}

// GetByInstanceId mocks base method.
func (m *MockIRevokeRepository) GetByInstanceId(instanceId uuid.UUID) (requests.Revoke, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByInstanceId", instanceId)
	ret0, _ := ret[0].(requests.Revoke)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByInstanceId indicates an expected call of GetByInstanceId.
func (mr *MockIRevokeRepositoryMockRecorder) GetByInstanceId(instanceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByInstanceId", reflect.TypeOf((*MockIRevokeRepository)(nil).GetByInstanceId), instanceId)
}

// GetByPickUpPointId mocks base method.
func (m *MockIRevokeRepository) GetByPickUpPointId(pickUpPointId uuid.UUID) (collection.Collection[requests.Revoke], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPickUpPointId", pickUpPointId)
	ret0, _ := ret[0].(collection.Collection[requests.Revoke])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPickUpPointId indicates an expected call of GetByPickUpPointId.
func (mr *MockIRevokeRepositoryMockRecorder) GetByPickUpPointId(pickUpPointId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPickUpPointId", reflect.TypeOf((*MockIRevokeRepository)(nil).GetByPickUpPointId), pickUpPointId)
}

// GetByUserId mocks base method.
func (m *MockIRevokeRepository) GetByUserId(userId uuid.UUID) (collection.Collection[requests.Revoke], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].(collection.Collection[requests.Revoke])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockIRevokeRepositoryMockRecorder) GetByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockIRevokeRepository)(nil).GetByUserId), userId)
}

// Remove mocks base method.
func (m *MockIRevokeRepository) Remove(requestId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", requestId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockIRevokeRepositoryMockRecorder) Remove(requestId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockIRevokeRepository)(nil).Remove), requestId)
}
