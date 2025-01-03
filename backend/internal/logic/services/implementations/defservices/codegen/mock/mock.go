// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=mock/mock.go
//

// Package mock_codegen is a generated GoMock package.
package mock_codegen

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIGenerator is a mock of IGenerator interface.
type MockIGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockIGeneratorMockRecorder
	isgomock struct{}
}

// MockIGeneratorMockRecorder is the mock recorder for MockIGenerator.
type MockIGeneratorMockRecorder struct {
	mock *MockIGenerator
}

// NewMockIGenerator creates a new mock instance.
func NewMockIGenerator(ctrl *gomock.Controller) *MockIGenerator {
	mock := &MockIGenerator{ctrl: ctrl}
	mock.recorder = &MockIGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGenerator) EXPECT() *MockIGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockIGenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockIGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockIGenerator)(nil).Generate))
}
