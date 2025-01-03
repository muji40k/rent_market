// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../implementation/mock/user/repository.go
//

// Package mock_user is a generated GoMock package.
package mock_user

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
func (m *MockIRepository) Create(user models.User) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIRepositoryMockRecorder) Create(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIRepository)(nil).Create), user)
}

// GetByEmail mocks base method.
func (m *MockIRepository) GetByEmail(email string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockIRepositoryMockRecorder) GetByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockIRepository)(nil).GetByEmail), email)
}

// GetById mocks base method.
func (m *MockIRepository) GetById(userId uuid.UUID) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", userId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRepositoryMockRecorder) GetById(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRepository)(nil).GetById), userId)
}

// GetByToken mocks base method.
func (m *MockIRepository) GetByToken(token models.Token) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByToken", token)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByToken indicates an expected call of GetByToken.
func (mr *MockIRepositoryMockRecorder) GetByToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByToken", reflect.TypeOf((*MockIRepository)(nil).GetByToken), token)
}

// Update mocks base method.
func (m *MockIRepository) Update(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIRepositoryMockRecorder) Update(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRepository)(nil).Update), user)
}

// MockIProfileRepository is a mock of IProfileRepository interface.
type MockIProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIProfileRepositoryMockRecorder
	isgomock struct{}
}

// MockIProfileRepositoryMockRecorder is the mock recorder for MockIProfileRepository.
type MockIProfileRepositoryMockRecorder struct {
	mock *MockIProfileRepository
}

// NewMockIProfileRepository creates a new mock instance.
func NewMockIProfileRepository(ctrl *gomock.Controller) *MockIProfileRepository {
	mock := &MockIProfileRepository{ctrl: ctrl}
	mock.recorder = &MockIProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProfileRepository) EXPECT() *MockIProfileRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIProfileRepository) Create(profile models.UserProfile) (models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", profile)
	ret0, _ := ret[0].(models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIProfileRepositoryMockRecorder) Create(profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIProfileRepository)(nil).Create), profile)
}

// GetByUserId mocks base method.
func (m *MockIProfileRepository) GetByUserId(userId uuid.UUID) (models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].(models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockIProfileRepositoryMockRecorder) GetByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockIProfileRepository)(nil).GetByUserId), userId)
}

// Update mocks base method.
func (m *MockIProfileRepository) Update(profile models.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIProfileRepositoryMockRecorder) Update(profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIProfileRepository)(nil).Update), profile)
}

// MockIFavoriteRepository is a mock of IFavoriteRepository interface.
type MockIFavoriteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIFavoriteRepositoryMockRecorder
	isgomock struct{}
}

// MockIFavoriteRepositoryMockRecorder is the mock recorder for MockIFavoriteRepository.
type MockIFavoriteRepositoryMockRecorder struct {
	mock *MockIFavoriteRepository
}

// NewMockIFavoriteRepository creates a new mock instance.
func NewMockIFavoriteRepository(ctrl *gomock.Controller) *MockIFavoriteRepository {
	mock := &MockIFavoriteRepository{ctrl: ctrl}
	mock.recorder = &MockIFavoriteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFavoriteRepository) EXPECT() *MockIFavoriteRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIFavoriteRepository) Create(profile models.UserFavoritePickUpPoint) (models.UserFavoritePickUpPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", profile)
	ret0, _ := ret[0].(models.UserFavoritePickUpPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIFavoriteRepositoryMockRecorder) Create(profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIFavoriteRepository)(nil).Create), profile)
}

// GetByUserId mocks base method.
func (m *MockIFavoriteRepository) GetByUserId(userId uuid.UUID) (models.UserFavoritePickUpPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].(models.UserFavoritePickUpPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockIFavoriteRepositoryMockRecorder) GetByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockIFavoriteRepository)(nil).GetByUserId), userId)
}

// Update mocks base method.
func (m *MockIFavoriteRepository) Update(profile models.UserFavoritePickUpPoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIFavoriteRepositoryMockRecorder) Update(profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIFavoriteRepository)(nil).Update), profile)
}

// MockIPayMethodsRepository is a mock of IPayMethodsRepository interface.
type MockIPayMethodsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIPayMethodsRepositoryMockRecorder
	isgomock struct{}
}

// MockIPayMethodsRepositoryMockRecorder is the mock recorder for MockIPayMethodsRepository.
type MockIPayMethodsRepositoryMockRecorder struct {
	mock *MockIPayMethodsRepository
}

// NewMockIPayMethodsRepository creates a new mock instance.
func NewMockIPayMethodsRepository(ctrl *gomock.Controller) *MockIPayMethodsRepository {
	mock := &MockIPayMethodsRepository{ctrl: ctrl}
	mock.recorder = &MockIPayMethodsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPayMethodsRepository) EXPECT() *MockIPayMethodsRepositoryMockRecorder {
	return m.recorder
}

// CreatePayMethod mocks base method.
func (m *MockIPayMethodsRepository) CreatePayMethod(userId uuid.UUID, payMethod models.UserPayMethod) (models.UserPayMethods, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayMethod", userId, payMethod)
	ret0, _ := ret[0].(models.UserPayMethods)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayMethod indicates an expected call of CreatePayMethod.
func (mr *MockIPayMethodsRepositoryMockRecorder) CreatePayMethod(userId, payMethod any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayMethod", reflect.TypeOf((*MockIPayMethodsRepository)(nil).CreatePayMethod), userId, payMethod)
}

// GetByUserId mocks base method.
func (m *MockIPayMethodsRepository) GetByUserId(userId uuid.UUID) (models.UserPayMethods, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].(models.UserPayMethods)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockIPayMethodsRepositoryMockRecorder) GetByUserId(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockIPayMethodsRepository)(nil).GetByUserId), userId)
}

// Update mocks base method.
func (m *MockIPayMethodsRepository) Update(payMethods models.UserPayMethods) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", payMethods)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIPayMethodsRepositoryMockRecorder) Update(payMethods any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIPayMethodsRepository)(nil).Update), payMethods)
}
