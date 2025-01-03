// Code generated by MockGen. DO NOT EDIT.
// Source: authenticator.go
//
// Generated by this command:
//
//	mockgen -source=authenticator.go -destination=implementations/mock/authenticator.go
//

// Package mock_authenticator is a generated GoMock package.
package mock_authenticator

import (
	reflect "reflect"
	models "rent_service/internal/domain/models"
	token "rent_service/internal/logic/services/types/token"

	gomock "go.uber.org/mock/gomock"
)

// MockIAuthenticator is a mock of IAuthenticator interface.
type MockIAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthenticatorMockRecorder
	isgomock struct{}
}

// MockIAuthenticatorMockRecorder is the mock recorder for MockIAuthenticator.
type MockIAuthenticatorMockRecorder struct {
	mock *MockIAuthenticator
}

// NewMockIAuthenticator creates a new mock instance.
func NewMockIAuthenticator(ctrl *gomock.Controller) *MockIAuthenticator {
	mock := &MockIAuthenticator{ctrl: ctrl}
	mock.recorder = &MockIAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthenticator) EXPECT() *MockIAuthenticatorMockRecorder {
	return m.recorder
}

// LoginWithToken mocks base method.
func (m *MockIAuthenticator) LoginWithToken(token token.Token) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginWithToken", token)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginWithToken indicates an expected call of LoginWithToken.
func (mr *MockIAuthenticatorMockRecorder) LoginWithToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginWithToken", reflect.TypeOf((*MockIAuthenticator)(nil).LoginWithToken), token)
}
