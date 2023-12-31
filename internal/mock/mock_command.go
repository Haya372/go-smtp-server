// Code generated by MockGen. DO NOT EDIT.
// Source: internal/command/handler.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	session "github.com/Haya372/smtp-server/internal/session"
	gomock "github.com/golang/mock/gomock"
)

// MockCommandHandler is a mock of CommandHandler interface.
type MockCommandHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCommandHandlerMockRecorder
}

// MockCommandHandlerMockRecorder is the mock recorder for MockCommandHandler.
type MockCommandHandlerMockRecorder struct {
	mock *MockCommandHandler
}

// NewMockCommandHandler creates a new mock instance.
func NewMockCommandHandler(ctrl *gomock.Controller) *MockCommandHandler {
	mock := &MockCommandHandler{ctrl: ctrl}
	mock.recorder = &MockCommandHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandHandler) EXPECT() *MockCommandHandlerMockRecorder {
	return m.recorder
}

// Command mocks base method.
func (m *MockCommandHandler) Command() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Command")
	ret0, _ := ret[0].(string)
	return ret0
}

// Command indicates an expected call of Command.
func (mr *MockCommandHandlerMockRecorder) Command() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Command", reflect.TypeOf((*MockCommandHandler)(nil).Command))
}

// HandleCommand mocks base method.
func (m *MockCommandHandler) HandleCommand(ctx context.Context, s *session.Session, arg []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleCommand", ctx, s, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleCommand indicates an expected call of HandleCommand.
func (mr *MockCommandHandlerMockRecorder) HandleCommand(ctx, s, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleCommand", reflect.TypeOf((*MockCommandHandler)(nil).HandleCommand), ctx, s, arg)
}
