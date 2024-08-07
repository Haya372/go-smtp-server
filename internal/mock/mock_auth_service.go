// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/auth.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	data "github.com/Haya372/smtp-server/internal/data"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockAuthService) Auth(ctx context.Context, mime data.MimeData) *data.AuthResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", ctx, mime)
	ret0, _ := ret[0].(*data.AuthResult)
	return ret0
}

// Auth indicates an expected call of Auth.
func (mr *MockAuthServiceMockRecorder) Auth(ctx, mime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockAuthService)(nil).Auth), ctx, mime)
}
