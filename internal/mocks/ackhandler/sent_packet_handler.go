// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lucas-clemente/quic-go/internal/ackhandler (interfaces: SentPacketHandler)

// Package mockackhandler is a generated GoMock package.
package mockackhandler

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	ackhandler "github.com/lucas-clemente/quic-go/internal/ackhandler"
	protocol "github.com/lucas-clemente/quic-go/internal/protocol"
	wire "github.com/lucas-clemente/quic-go/internal/wire"
	quictrace "github.com/lucas-clemente/quic-go/quictrace"
)

// MockSentPacketHandler is a mock of SentPacketHandler interface
type MockSentPacketHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSentPacketHandlerMockRecorder
}

// MockSentPacketHandlerMockRecorder is the mock recorder for MockSentPacketHandler
type MockSentPacketHandlerMockRecorder struct {
	mock *MockSentPacketHandler
}

// NewMockSentPacketHandler creates a new mock instance
func NewMockSentPacketHandler(ctrl *gomock.Controller) *MockSentPacketHandler {
	mock := &MockSentPacketHandler{ctrl: ctrl}
	mock.recorder = &MockSentPacketHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSentPacketHandler) EXPECT() *MockSentPacketHandlerMockRecorder {
	return m.recorder
}

// DequeuePacketForRetransmission mocks base method
func (m *MockSentPacketHandler) DequeuePacketForRetransmission() *ackhandler.Packet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DequeuePacketForRetransmission")
	ret0, _ := ret[0].(*ackhandler.Packet)
	return ret0
}

// DequeuePacketForRetransmission indicates an expected call of DequeuePacketForRetransmission
func (mr *MockSentPacketHandlerMockRecorder) DequeuePacketForRetransmission() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DequeuePacketForRetransmission", reflect.TypeOf((*MockSentPacketHandler)(nil).DequeuePacketForRetransmission))
}

// DequeueProbePacket mocks base method
func (m *MockSentPacketHandler) DequeueProbePacket() (*ackhandler.Packet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DequeueProbePacket")
	ret0, _ := ret[0].(*ackhandler.Packet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DequeueProbePacket indicates an expected call of DequeueProbePacket
func (mr *MockSentPacketHandlerMockRecorder) DequeueProbePacket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DequeueProbePacket", reflect.TypeOf((*MockSentPacketHandler)(nil).DequeueProbePacket))
}

// DropPackets mocks base method
func (m *MockSentPacketHandler) DropPackets(arg0 protocol.EncryptionLevel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DropPackets", arg0)
}

// DropPackets indicates an expected call of DropPackets
func (mr *MockSentPacketHandlerMockRecorder) DropPackets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropPackets", reflect.TypeOf((*MockSentPacketHandler)(nil).DropPackets), arg0)
}

// GetAlarmTimeout mocks base method
func (m *MockSentPacketHandler) GetAlarmTimeout() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlarmTimeout")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetAlarmTimeout indicates an expected call of GetAlarmTimeout
func (mr *MockSentPacketHandlerMockRecorder) GetAlarmTimeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlarmTimeout", reflect.TypeOf((*MockSentPacketHandler)(nil).GetAlarmTimeout))
}

// GetLowestPacketNotConfirmedAcked mocks base method
func (m *MockSentPacketHandler) GetLowestPacketNotConfirmedAcked() protocol.PacketNumber {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLowestPacketNotConfirmedAcked")
	ret0, _ := ret[0].(protocol.PacketNumber)
	return ret0
}

// GetLowestPacketNotConfirmedAcked indicates an expected call of GetLowestPacketNotConfirmedAcked
func (mr *MockSentPacketHandlerMockRecorder) GetLowestPacketNotConfirmedAcked() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLowestPacketNotConfirmedAcked", reflect.TypeOf((*MockSentPacketHandler)(nil).GetLowestPacketNotConfirmedAcked))
}

// GetStats mocks base method
func (m *MockSentPacketHandler) GetStats() *quictrace.TransportState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStats")
	ret0, _ := ret[0].(*quictrace.TransportState)
	return ret0
}

// GetStats indicates an expected call of GetStats
func (mr *MockSentPacketHandlerMockRecorder) GetStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStats", reflect.TypeOf((*MockSentPacketHandler)(nil).GetStats))
}

// OnAlarm mocks base method
func (m *MockSentPacketHandler) OnAlarm() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnAlarm")
	ret0, _ := ret[0].(error)
	return ret0
}

// OnAlarm indicates an expected call of OnAlarm
func (mr *MockSentPacketHandlerMockRecorder) OnAlarm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnAlarm", reflect.TypeOf((*MockSentPacketHandler)(nil).OnAlarm))
}

// PeekPacketNumber mocks base method
func (m *MockSentPacketHandler) PeekPacketNumber(arg0 protocol.EncryptionLevel) (protocol.PacketNumber, protocol.PacketNumberLen) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekPacketNumber", arg0)
	ret0, _ := ret[0].(protocol.PacketNumber)
	ret1, _ := ret[1].(protocol.PacketNumberLen)
	return ret0, ret1
}

// PeekPacketNumber indicates an expected call of PeekPacketNumber
func (mr *MockSentPacketHandlerMockRecorder) PeekPacketNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekPacketNumber", reflect.TypeOf((*MockSentPacketHandler)(nil).PeekPacketNumber), arg0)
}

// PopPacketNumber mocks base method
func (m *MockSentPacketHandler) PopPacketNumber(arg0 protocol.EncryptionLevel) protocol.PacketNumber {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopPacketNumber", arg0)
	ret0, _ := ret[0].(protocol.PacketNumber)
	return ret0
}

// PopPacketNumber indicates an expected call of PopPacketNumber
func (mr *MockSentPacketHandlerMockRecorder) PopPacketNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopPacketNumber", reflect.TypeOf((*MockSentPacketHandler)(nil).PopPacketNumber), arg0)
}

// ReceivedAck mocks base method
func (m *MockSentPacketHandler) ReceivedAck(arg0 *wire.AckFrame, arg1 protocol.PacketNumber, arg2 protocol.EncryptionLevel, arg3 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceivedAck", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReceivedAck indicates an expected call of ReceivedAck
func (mr *MockSentPacketHandlerMockRecorder) ReceivedAck(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceivedAck", reflect.TypeOf((*MockSentPacketHandler)(nil).ReceivedAck), arg0, arg1, arg2, arg3)
}

// ResetForRetry mocks base method
func (m *MockSentPacketHandler) ResetForRetry() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetForRetry")
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetForRetry indicates an expected call of ResetForRetry
func (mr *MockSentPacketHandlerMockRecorder) ResetForRetry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetForRetry", reflect.TypeOf((*MockSentPacketHandler)(nil).ResetForRetry))
}

// SendMode mocks base method
func (m *MockSentPacketHandler) SendMode() ackhandler.SendMode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMode")
	ret0, _ := ret[0].(ackhandler.SendMode)
	return ret0
}

// SendMode indicates an expected call of SendMode
func (mr *MockSentPacketHandlerMockRecorder) SendMode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMode", reflect.TypeOf((*MockSentPacketHandler)(nil).SendMode))
}

// SentPacket mocks base method
func (m *MockSentPacketHandler) SentPacket(arg0 *ackhandler.Packet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SentPacket", arg0)
}

// SentPacket indicates an expected call of SentPacket
func (mr *MockSentPacketHandlerMockRecorder) SentPacket(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentPacket", reflect.TypeOf((*MockSentPacketHandler)(nil).SentPacket), arg0)
}

// SentPacketsAsRetransmission mocks base method
func (m *MockSentPacketHandler) SentPacketsAsRetransmission(arg0 []*ackhandler.Packet, arg1 protocol.PacketNumber) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SentPacketsAsRetransmission", arg0, arg1)
}

// SentPacketsAsRetransmission indicates an expected call of SentPacketsAsRetransmission
func (mr *MockSentPacketHandlerMockRecorder) SentPacketsAsRetransmission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentPacketsAsRetransmission", reflect.TypeOf((*MockSentPacketHandler)(nil).SentPacketsAsRetransmission), arg0, arg1)
}

// SetMaxAckDelay mocks base method
func (m *MockSentPacketHandler) SetMaxAckDelay(arg0 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMaxAckDelay", arg0)
}

// SetMaxAckDelay indicates an expected call of SetMaxAckDelay
func (mr *MockSentPacketHandlerMockRecorder) SetMaxAckDelay(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxAckDelay", reflect.TypeOf((*MockSentPacketHandler)(nil).SetMaxAckDelay), arg0)
}

// ShouldSendNumPackets mocks base method
func (m *MockSentPacketHandler) ShouldSendNumPackets() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldSendNumPackets")
	ret0, _ := ret[0].(int)
	return ret0
}

// ShouldSendNumPackets indicates an expected call of ShouldSendNumPackets
func (mr *MockSentPacketHandlerMockRecorder) ShouldSendNumPackets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldSendNumPackets", reflect.TypeOf((*MockSentPacketHandler)(nil).ShouldSendNumPackets))
}

// TimeUntilSend mocks base method
func (m *MockSentPacketHandler) TimeUntilSend() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeUntilSend")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// TimeUntilSend indicates an expected call of TimeUntilSend
func (mr *MockSentPacketHandlerMockRecorder) TimeUntilSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeUntilSend", reflect.TypeOf((*MockSentPacketHandler)(nil).TimeUntilSend))
}
