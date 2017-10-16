// Code generated by MockGen. DO NOT EDIT.
// Source: ../flowcontrol/interface.go

package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/flowcontrol"
	handshake "github.com/lucas-clemente/quic-go/internal/handshake"
	protocol "github.com/lucas-clemente/quic-go/internal/protocol"
)

// MockFlowControlManager is a mock of FlowControlManager interface
type MockFlowControlManager struct {
	ctrl     *gomock.Controller
	recorder *MockFlowControlManagerMockRecorder
}

// MockFlowControlManagerMockRecorder is the mock recorder for MockFlowControlManager
type MockFlowControlManagerMockRecorder struct {
	mock *MockFlowControlManager
}

// NewMockFlowControlManager creates a new mock instance
func NewMockFlowControlManager(ctrl *gomock.Controller) *MockFlowControlManager {
	mock := &MockFlowControlManager{ctrl: ctrl}
	mock.recorder = &MockFlowControlManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockFlowControlManager) EXPECT() *MockFlowControlManagerMockRecorder {
	return _m.recorder
}

// NewStream mocks base method
func (_m *MockFlowControlManager) NewStream(streamID protocol.StreamID, contributesToConnectionFlow bool) {
	_m.ctrl.Call(_m, "NewStream", streamID, contributesToConnectionFlow)
}

// NewStream indicates an expected call of NewStream
func (_mr *MockFlowControlManagerMockRecorder) NewStream(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "NewStream", reflect.TypeOf((*MockFlowControlManager)(nil).NewStream), arg0, arg1)
}

// RemoveStream mocks base method
func (_m *MockFlowControlManager) RemoveStream(streamID protocol.StreamID) {
	_m.ctrl.Call(_m, "RemoveStream", streamID)
}

// RemoveStream indicates an expected call of RemoveStream
func (_mr *MockFlowControlManagerMockRecorder) RemoveStream(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "RemoveStream", reflect.TypeOf((*MockFlowControlManager)(nil).RemoveStream), arg0)
}

// UpdateTransportParameters mocks base method
func (_m *MockFlowControlManager) UpdateTransportParameters(_param0 *handshake.TransportParameters) {
	_m.ctrl.Call(_m, "UpdateTransportParameters", _param0)
}

// UpdateTransportParameters indicates an expected call of UpdateTransportParameters
func (_mr *MockFlowControlManagerMockRecorder) UpdateTransportParameters(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateTransportParameters", reflect.TypeOf((*MockFlowControlManager)(nil).UpdateTransportParameters), arg0)
}

// ResetStream mocks base method
func (_m *MockFlowControlManager) ResetStream(streamID protocol.StreamID, byteOffset protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "ResetStream", streamID, byteOffset)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetStream indicates an expected call of ResetStream
func (_mr *MockFlowControlManagerMockRecorder) ResetStream(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ResetStream", reflect.TypeOf((*MockFlowControlManager)(nil).ResetStream), arg0, arg1)
}

// UpdateHighestReceived mocks base method
func (_m *MockFlowControlManager) UpdateHighestReceived(streamID protocol.StreamID, byteOffset protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "UpdateHighestReceived", streamID, byteOffset)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHighestReceived indicates an expected call of UpdateHighestReceived
func (_mr *MockFlowControlManagerMockRecorder) UpdateHighestReceived(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateHighestReceived", reflect.TypeOf((*MockFlowControlManager)(nil).UpdateHighestReceived), arg0, arg1)
}

// AddBytesRead mocks base method
func (_m *MockFlowControlManager) AddBytesRead(streamID protocol.StreamID, n protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "AddBytesRead", streamID, n)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBytesRead indicates an expected call of AddBytesRead
func (_mr *MockFlowControlManagerMockRecorder) AddBytesRead(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddBytesRead", reflect.TypeOf((*MockFlowControlManager)(nil).AddBytesRead), arg0, arg1)
}

// GetWindowUpdates mocks base method
func (_m *MockFlowControlManager) GetWindowUpdates() []flowcontrol.WindowUpdate {
	ret := _m.ctrl.Call(_m, "GetWindowUpdates")
	ret0, _ := ret[0].([]flowcontrol.WindowUpdate)
	return ret0
}

// GetWindowUpdates indicates an expected call of GetWindowUpdates
func (_mr *MockFlowControlManagerMockRecorder) GetWindowUpdates() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetWindowUpdates", reflect.TypeOf((*MockFlowControlManager)(nil).GetWindowUpdates))
}

// GetReceiveWindow mocks base method
func (_m *MockFlowControlManager) GetReceiveWindow(streamID protocol.StreamID) (protocol.ByteCount, error) {
	ret := _m.ctrl.Call(_m, "GetReceiveWindow", streamID)
	ret0, _ := ret[0].(protocol.ByteCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceiveWindow indicates an expected call of GetReceiveWindow
func (_mr *MockFlowControlManagerMockRecorder) GetReceiveWindow(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetReceiveWindow", reflect.TypeOf((*MockFlowControlManager)(nil).GetReceiveWindow), arg0)
}

// AddBytesSent mocks base method
func (_m *MockFlowControlManager) AddBytesSent(streamID protocol.StreamID, n protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "AddBytesSent", streamID, n)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBytesSent indicates an expected call of AddBytesSent
func (_mr *MockFlowControlManagerMockRecorder) AddBytesSent(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddBytesSent", reflect.TypeOf((*MockFlowControlManager)(nil).AddBytesSent), arg0, arg1)
}

// SendWindowSize mocks base method
func (_m *MockFlowControlManager) SendWindowSize(streamID protocol.StreamID) (protocol.ByteCount, error) {
	ret := _m.ctrl.Call(_m, "SendWindowSize", streamID)
	ret0, _ := ret[0].(protocol.ByteCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendWindowSize indicates an expected call of SendWindowSize
func (_mr *MockFlowControlManagerMockRecorder) SendWindowSize(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SendWindowSize", reflect.TypeOf((*MockFlowControlManager)(nil).SendWindowSize), arg0)
}

// RemainingConnectionWindowSize mocks base method
func (_m *MockFlowControlManager) RemainingConnectionWindowSize() protocol.ByteCount {
	ret := _m.ctrl.Call(_m, "RemainingConnectionWindowSize")
	ret0, _ := ret[0].(protocol.ByteCount)
	return ret0
}

// RemainingConnectionWindowSize indicates an expected call of RemainingConnectionWindowSize
func (_mr *MockFlowControlManagerMockRecorder) RemainingConnectionWindowSize() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "RemainingConnectionWindowSize", reflect.TypeOf((*MockFlowControlManager)(nil).RemainingConnectionWindowSize))
}

// UpdateStreamWindow mocks base method
func (_m *MockFlowControlManager) UpdateStreamWindow(streamID protocol.StreamID, offset protocol.ByteCount) (bool, error) {
	ret := _m.ctrl.Call(_m, "UpdateStreamWindow", streamID, offset)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStreamWindow indicates an expected call of UpdateStreamWindow
func (_mr *MockFlowControlManagerMockRecorder) UpdateStreamWindow(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateStreamWindow", reflect.TypeOf((*MockFlowControlManager)(nil).UpdateStreamWindow), arg0, arg1)
}

// UpdateConnectionWindow mocks base method
func (_m *MockFlowControlManager) UpdateConnectionWindow(offset protocol.ByteCount) bool {
	ret := _m.ctrl.Call(_m, "UpdateConnectionWindow", offset)
	ret0, _ := ret[0].(bool)
	return ret0
}

// UpdateConnectionWindow indicates an expected call of UpdateConnectionWindow
func (_mr *MockFlowControlManagerMockRecorder) UpdateConnectionWindow(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateConnectionWindow", reflect.TypeOf((*MockFlowControlManager)(nil).UpdateConnectionWindow), arg0)
}
