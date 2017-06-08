// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/lucas-clemente/quic-go/flowcontrol (interfaces: FlowControlManager)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	flowcontrol "github.com/lucas-clemente/quic-go/flowcontrol"
	protocol "github.com/lucas-clemente/quic-go/protocol"
)

// Mock of FlowControlManager interface
type MockFlowControlManager struct {
	ctrl     *gomock.Controller
	recorder *_MockFlowControlManagerRecorder
}

// Recorder for MockFlowControlManager (not exported)
type _MockFlowControlManagerRecorder struct {
	mock *MockFlowControlManager
}

func NewMockFlowControlManager(ctrl *gomock.Controller) *MockFlowControlManager {
	mock := &MockFlowControlManager{ctrl: ctrl}
	mock.recorder = &_MockFlowControlManagerRecorder{mock}
	return mock
}

func (_m *MockFlowControlManager) EXPECT() *_MockFlowControlManagerRecorder {
	return _m.recorder
}

func (_m *MockFlowControlManager) AddBytesRead(_param0 protocol.StreamID, _param1 protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "AddBytesRead", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) AddBytesRead(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddBytesRead", arg0, arg1)
}

func (_m *MockFlowControlManager) AddBytesSent(_param0 protocol.StreamID, _param1 protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "AddBytesSent", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) AddBytesSent(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddBytesSent", arg0, arg1)
}

func (_m *MockFlowControlManager) GetReceiveWindow(_param0 protocol.StreamID) (protocol.ByteCount, error) {
	ret := _m.ctrl.Call(_m, "GetReceiveWindow", _param0)
	ret0, _ := ret[0].(protocol.ByteCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFlowControlManagerRecorder) GetReceiveWindow(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetReceiveWindow", arg0)
}

func (_m *MockFlowControlManager) GetWindowUpdates() []flowcontrol.WindowUpdate {
	ret := _m.ctrl.Call(_m, "GetWindowUpdates")
	ret0, _ := ret[0].([]flowcontrol.WindowUpdate)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) GetWindowUpdates() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetWindowUpdates")
}

func (_m *MockFlowControlManager) NewStream(_param0 protocol.StreamID, _param1 bool) {
	_m.ctrl.Call(_m, "NewStream", _param0, _param1)
}

func (_mr *_MockFlowControlManagerRecorder) NewStream(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewStream", arg0, arg1)
}

func (_m *MockFlowControlManager) RemainingConnectionWindowSize() protocol.ByteCount {
	ret := _m.ctrl.Call(_m, "RemainingConnectionWindowSize")
	ret0, _ := ret[0].(protocol.ByteCount)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) RemainingConnectionWindowSize() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemainingConnectionWindowSize")
}

func (_m *MockFlowControlManager) RemoveStream(_param0 protocol.StreamID) {
	_m.ctrl.Call(_m, "RemoveStream", _param0)
}

func (_mr *_MockFlowControlManagerRecorder) RemoveStream(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveStream", arg0)
}

func (_m *MockFlowControlManager) ResetStream(_param0 protocol.StreamID, _param1 protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "ResetStream", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) ResetStream(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ResetStream", arg0, arg1)
}

func (_m *MockFlowControlManager) SendWindowSize(_param0 protocol.StreamID) (protocol.ByteCount, error) {
	ret := _m.ctrl.Call(_m, "SendWindowSize", _param0)
	ret0, _ := ret[0].(protocol.ByteCount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFlowControlManagerRecorder) SendWindowSize(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SendWindowSize", arg0)
}

func (_m *MockFlowControlManager) UpdateHighestReceived(_param0 protocol.StreamID, _param1 protocol.ByteCount) error {
	ret := _m.ctrl.Call(_m, "UpdateHighestReceived", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockFlowControlManagerRecorder) UpdateHighestReceived(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateHighestReceived", arg0, arg1)
}

func (_m *MockFlowControlManager) UpdateWindow(_param0 protocol.StreamID, _param1 protocol.ByteCount) (bool, error) {
	ret := _m.ctrl.Call(_m, "UpdateWindow", _param0, _param1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFlowControlManagerRecorder) UpdateWindow(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateWindow", arg0, arg1)
}
