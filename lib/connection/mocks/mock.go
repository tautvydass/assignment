// Code generated by MockGen. DO NOT EDIT.
// Source: assignment/lib/connection (interfaces: Connection,ReadWriteStream)

// Package mocks is a generated GoMock package.
package mocks

import (
	connection "assignment/lib/connection"
	entity "assignment/lib/entity"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockConnection is a mock of Connection interface.
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection.
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance.
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// AcceptReadStream mocks base method.
func (m *MockConnection) AcceptReadStream(arg0 context.Context, arg1 connection.MessageReceiver) (connection.ReadStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptReadStream", arg0, arg1)
	ret0, _ := ret[0].(connection.ReadStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptReadStream indicates an expected call of AcceptReadStream.
func (mr *MockConnectionMockRecorder) AcceptReadStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptReadStream", reflect.TypeOf((*MockConnection)(nil).AcceptReadStream), arg0, arg1)
}

// AcceptReadWriteStream mocks base method.
func (m *MockConnection) AcceptReadWriteStream(arg0 context.Context, arg1 connection.MessageReceiver) (connection.ReadWriteStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptReadWriteStream", arg0, arg1)
	ret0, _ := ret[0].(connection.ReadWriteStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptReadWriteStream indicates an expected call of AcceptReadWriteStream.
func (mr *MockConnectionMockRecorder) AcceptReadWriteStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptReadWriteStream", reflect.TypeOf((*MockConnection)(nil).AcceptReadWriteStream), arg0, arg1)
}

// OpenReadWriteStream mocks base method.
func (m *MockConnection) OpenReadWriteStream(arg0 context.Context, arg1 connection.MessageReceiver) (connection.ReadWriteStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenReadWriteStream", arg0, arg1)
	ret0, _ := ret[0].(connection.ReadWriteStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenReadWriteStream indicates an expected call of OpenReadWriteStream.
func (mr *MockConnectionMockRecorder) OpenReadWriteStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenReadWriteStream", reflect.TypeOf((*MockConnection)(nil).OpenReadWriteStream), arg0, arg1)
}

// OpenWriteStream mocks base method.
func (m *MockConnection) OpenWriteStream(arg0 context.Context) (connection.WriteStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenWriteStream", arg0)
	ret0, _ := ret[0].(connection.WriteStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenWriteStream indicates an expected call of OpenWriteStream.
func (mr *MockConnectionMockRecorder) OpenWriteStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenWriteStream", reflect.TypeOf((*MockConnection)(nil).OpenWriteStream), arg0)
}

// MockReadWriteStream is a mock of ReadWriteStream interface.
type MockReadWriteStream struct {
	ctrl     *gomock.Controller
	recorder *MockReadWriteStreamMockRecorder
}

// MockReadWriteStreamMockRecorder is the mock recorder for MockReadWriteStream.
type MockReadWriteStreamMockRecorder struct {
	mock *MockReadWriteStream
}

// NewMockReadWriteStream creates a new mock instance.
func NewMockReadWriteStream(ctrl *gomock.Controller) *MockReadWriteStream {
	mock := &MockReadWriteStream{ctrl: ctrl}
	mock.recorder = &MockReadWriteStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadWriteStream) EXPECT() *MockReadWriteStreamMockRecorder {
	return m.recorder
}

// CloseStream mocks base method.
func (m *MockReadWriteStream) CloseStream() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseStream")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseStream indicates an expected call of CloseStream.
func (mr *MockReadWriteStreamMockRecorder) CloseStream() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseStream", reflect.TypeOf((*MockReadWriteStream)(nil).CloseStream))
}

// SendMessage mocks base method.
func (m *MockReadWriteStream) SendMessage(arg0 entity.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockReadWriteStreamMockRecorder) SendMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockReadWriteStream)(nil).SendMessage), arg0)
}

// SetConnClosedCallback mocks base method.
func (m *MockReadWriteStream) SetConnClosedCallback(arg0 connection.ConnClosedCallback) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnClosedCallback", arg0)
}

// SetConnClosedCallback indicates an expected call of SetConnClosedCallback.
func (mr *MockReadWriteStreamMockRecorder) SetConnClosedCallback(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnClosedCallback", reflect.TypeOf((*MockReadWriteStream)(nil).SetConnClosedCallback), arg0)
}

// SetMessageReceiver mocks base method.
func (m *MockReadWriteStream) SetMessageReceiver(arg0 connection.MessageReceiver) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMessageReceiver", arg0)
}

// SetMessageReceiver indicates an expected call of SetMessageReceiver.
func (mr *MockReadWriteStreamMockRecorder) SetMessageReceiver(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMessageReceiver", reflect.TypeOf((*MockReadWriteStream)(nil).SetMessageReceiver), arg0)
}

// SetReadBufferSize mocks base method.
func (m *MockReadWriteStream) SetReadBufferSize(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReadBufferSize", arg0)
}

// SetReadBufferSize indicates an expected call of SetReadBufferSize.
func (mr *MockReadWriteStreamMockRecorder) SetReadBufferSize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadBufferSize", reflect.TypeOf((*MockReadWriteStream)(nil).SetReadBufferSize), arg0)
}

// SetSendMessageTimeout mocks base method.
func (m *MockReadWriteStream) SetSendMessageTimeout(arg0 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetSendMessageTimeout", arg0)
}

// SetSendMessageTimeout indicates an expected call of SetSendMessageTimeout.
func (mr *MockReadWriteStreamMockRecorder) SetSendMessageTimeout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSendMessageTimeout", reflect.TypeOf((*MockReadWriteStream)(nil).SetSendMessageTimeout), arg0)
}
