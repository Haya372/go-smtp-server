// Code generated by MockGen. DO NOT EDIT.
// Source: internal/session/factory.go

// Package mock is a generated GoMock package.
package mock

import (
	net "net"
	reflect "reflect"

	session "github.com/Haya372/smtp-server/internal/session"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionFactory is a mock of SessionFactory interface.
type MockSessionFactory struct {
	ctrl     *gomock.Controller
	recorder *MockSessionFactoryMockRecorder
}

// MockSessionFactoryMockRecorder is the mock recorder for MockSessionFactory.
type MockSessionFactoryMockRecorder struct {
	mock *MockSessionFactory
}

// NewMockSessionFactory creates a new mock instance.
func NewMockSessionFactory(ctrl *gomock.Controller) *MockSessionFactory {
	mock := &MockSessionFactory{ctrl: ctrl}
	mock.recorder = &MockSessionFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionFactory) EXPECT() *MockSessionFactoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionFactory) CreateSession(conn net.Conn) session.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", conn)
	ret0, _ := ret[0].(session.Session)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionFactoryMockRecorder) CreateSession(conn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionFactory)(nil).CreateSession), conn)
}