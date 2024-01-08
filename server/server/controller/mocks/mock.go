// Code generated by MockGen. DO NOT EDIT.
// Source: assignment/server/server/controller (interfaces: CommsController)

// Package mocks is a generated GoMock package.
package mocks

import (
	connection "assignment/lib/connection"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommsController is a mock of CommsController interface.
type MockCommsController struct {
	ctrl     *gomock.Controller
	recorder *MockCommsControllerMockRecorder
}

// MockCommsControllerMockRecorder is the mock recorder for MockCommsController.
type MockCommsControllerMockRecorder struct {
	mock *MockCommsController
}

// NewMockCommsController creates a new mock instance.
func NewMockCommsController(ctrl *gomock.Controller) *MockCommsController {
	mock := &MockCommsController{ctrl: ctrl}
	mock.recorder = &MockCommsControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommsController) EXPECT() *MockCommsControllerMockRecorder {
	return m.recorder
}

// AddPublisher mocks base method.
func (m *MockCommsController) AddPublisher(arg0 connection.ReadWriteStream) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddPublisher", arg0)
}

// AddPublisher indicates an expected call of AddPublisher.
func (mr *MockCommsControllerMockRecorder) AddPublisher(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPublisher", reflect.TypeOf((*MockCommsController)(nil).AddPublisher), arg0)
}

// AddSubscriber mocks base method.
func (m *MockCommsController) AddSubscriber(arg0 connection.WriteStream) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddSubscriber", arg0)
}

// AddSubscriber indicates an expected call of AddSubscriber.
func (mr *MockCommsControllerMockRecorder) AddSubscriber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubscriber", reflect.TypeOf((*MockCommsController)(nil).AddSubscriber), arg0)
}

// Close mocks base method.
func (m *MockCommsController) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockCommsControllerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCommsController)(nil).Close))
}

// MessageReceiver mocks base method.
func (m *MockCommsController) MessageReceiver() connection.MessageReceiver {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MessageReceiver")
	ret0, _ := ret[0].(connection.MessageReceiver)
	return ret0
}

// MessageReceiver indicates an expected call of MessageReceiver.
func (mr *MockCommsControllerMockRecorder) MessageReceiver() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MessageReceiver", reflect.TypeOf((*MockCommsController)(nil).MessageReceiver))
}
