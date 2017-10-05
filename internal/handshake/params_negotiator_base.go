package handshake

import (
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// The ParamsNegotiator negotiates and stores the connection parameters.
// It can be used for a server as well as a client.
type ParamsNegotiator interface {
	GetSendStreamFlowControlWindow() protocol.ByteCount
	GetSendConnectionFlowControlWindow() protocol.ByteCount
	GetReceiveStreamFlowControlWindow() protocol.ByteCount
	GetReceiveConnectionFlowControlWindow() protocol.ByteCount
	GetMaxOutgoingStreams() uint32
	GetMaxIncomingStreams() uint32
	// get the idle timeout that was sent by the peer
	GetRemoteIdleTimeout() time.Duration
	// determines if the client requests omission of connection IDs.
	OmitConnectionID() bool
}

// For the server:
// 1. call SetFromMap with the values received in the CHLO. This sets the corresponding values here, subject to negotiation
// 2. call GetHelloMap to get the values to send in the SHLO
// For the client:
// 1. call GetHelloMap to get the values to send in a CHLO
// 2. call SetFromMap with the values received in the SHLO
type paramsNegotiatorBase struct {
	mutex sync.RWMutex

	version     protocol.VersionNumber
	perspective protocol.Perspective

	flowControlNegotiated bool

	omitConnectionID            bool
	requestConnectionIDOmission bool

	maxStreamsPerConnection                uint32
	maxIncomingDynamicStreamsPerConnection uint32
	idleTimeout                            time.Duration
	remoteIdleTimeout                      time.Duration
	sendStreamFlowControlWindow            protocol.ByteCount
	sendConnectionFlowControlWindow        protocol.ByteCount
	receiveStreamFlowControlWindow         protocol.ByteCount
	receiveConnectionFlowControlWindow     protocol.ByteCount
}

func (h *paramsNegotiatorBase) init(params *TransportParameters) {
	h.sendStreamFlowControlWindow = protocol.InitialStreamFlowControlWindow         // can only be changed by the client
	h.sendConnectionFlowControlWindow = protocol.InitialConnectionFlowControlWindow // can only be changed by the client
	h.receiveStreamFlowControlWindow = protocol.ReceiveStreamFlowControlWindow
	h.receiveConnectionFlowControlWindow = protocol.ReceiveConnectionFlowControlWindow
	h.requestConnectionIDOmission = params.RequestConnectionIDOmission

	h.idleTimeout = params.IdleTimeout
	if h.perspective == protocol.PerspectiveServer {
		h.maxStreamsPerConnection = protocol.MaxStreamsPerConnection                // this is the value negotiated based on what the client sent
		h.maxIncomingDynamicStreamsPerConnection = protocol.MaxStreamsPerConnection // "incoming" seen from the client's perspective
	} else {
		h.maxStreamsPerConnection = protocol.MaxStreamsPerConnection                // this is the value negotiated based on what the client sent
		h.maxIncomingDynamicStreamsPerConnection = protocol.MaxStreamsPerConnection // "incoming" seen from the server's perspective
	}
}

func (h *paramsNegotiatorBase) negotiateMaxStreamsPerConnection(clientValue uint32) uint32 {
	return utils.MinUint32(clientValue, protocol.MaxStreamsPerConnection)
}

func (h *paramsNegotiatorBase) negotiateMaxIncomingDynamicStreamsPerConnection(clientValue uint32) uint32 {
	return utils.MinUint32(clientValue, protocol.MaxIncomingDynamicStreamsPerConnection)
}

// GetSendStreamFlowControlWindow gets the size of the stream-level flow control window for sending data
func (h *paramsNegotiatorBase) GetSendStreamFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.sendStreamFlowControlWindow
}

// GetSendConnectionFlowControlWindow gets the size of the stream-level flow control window for sending data
func (h *paramsNegotiatorBase) GetSendConnectionFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.sendConnectionFlowControlWindow
}

func (h *paramsNegotiatorBase) GetReceiveStreamFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.receiveStreamFlowControlWindow
}

// GetReceiveConnectionFlowControlWindow gets the size of the stream-level flow control window for receiving data
func (h *paramsNegotiatorBase) GetReceiveConnectionFlowControlWindow() protocol.ByteCount {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.receiveConnectionFlowControlWindow
}

func (h *paramsNegotiatorBase) GetMaxOutgoingStreams() uint32 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.maxIncomingDynamicStreamsPerConnection
}

func (h *paramsNegotiatorBase) GetMaxIncomingStreams() uint32 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	maxStreams := protocol.MaxIncomingDynamicStreamsPerConnection
	return utils.MaxUint32(uint32(maxStreams)+protocol.MaxStreamsMinimumIncrement, uint32(float64(maxStreams)*protocol.MaxStreamsMultiplier))
}

func (h *paramsNegotiatorBase) setRemoteIdleTimeout(t time.Duration) {
	h.remoteIdleTimeout = utils.MaxDuration(protocol.MinRemoteIdleTimeout, t)
}

func (h *paramsNegotiatorBase) GetRemoteIdleTimeout() time.Duration {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.remoteIdleTimeout
}
