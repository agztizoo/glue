// Code generated by MockGen. DO NOT EDIT.

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// EscapeTransaction mocks base method.
func (m *MockManager) EscapeTransaction(arg0 context.Context, arg1 func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EscapeTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EscapeTransaction indicates an expected call of EscapeTransaction.
func (mr *MockManagerMockRecorder) EscapeTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EscapeTransaction", reflect.TypeOf((*MockManager)(nil).EscapeTransaction), arg0, arg1)
}

// OnCommitted mocks base method.
func (m *MockManager) OnCommitted(arg0 context.Context, arg1 func(context.Context)) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnCommitted", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// OnCommitted indicates an expected call of OnCommitted.
func (mr *MockManagerMockRecorder) OnCommitted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnCommitted", reflect.TypeOf((*MockManager)(nil).OnCommitted), arg0, arg1)
}

// Transaction mocks base method.
func (m *MockManager) Transaction(arg0 context.Context, arg1 func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockManagerMockRecorder) Transaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockManager)(nil).Transaction), arg0, arg1)
}
