package protocol

import "math"

// A PacketNumber in QUIC
type PacketNumber uint64

// PacketNumberLen is the length of the packet number in bytes
type PacketNumberLen uint8

const (
	// PacketNumberLenInvalid is the default value and not a valid length for a packet number
	PacketNumberLenInvalid PacketNumberLen = 0
	// PacketNumberLen1 is a packet number length of 1 byte
	PacketNumberLen1 PacketNumberLen = 1
	// PacketNumberLen2 is a packet number length of 2 bytes
	PacketNumberLen2 PacketNumberLen = 2
	// PacketNumberLen4 is a packet number length of 4 bytes
	PacketNumberLen4 PacketNumberLen = 4
	// PacketNumberLen6 is a packet number length of 6 bytes
	PacketNumberLen6 PacketNumberLen = 6
)

// The PacketType is the Long Header Type (only used for the IETF draft header format)
type PacketType uint8

const (
	// PacketTypeVersionNegotiation is the packet type of a Version Negotiation packet
	PacketTypeVersionNegotiation PacketType = 1
	// PacketTypeClientInitial is the packet type of a Client Initial packet
	PacketTypeClientInitial PacketType = 2
	// PacketTypeServerStatelessRetry is the packet type of a Server Stateless Retry packet
	PacketTypeServerStatelessRetry PacketType = 3
	// PacketTypeServerCleartext is the packet type of a Server Cleartext packet
	PacketTypeServerCleartext PacketType = 4
	// PacketTypeClientCleartext is the packet type of a Client Cleartext packet
	PacketTypeClientCleartext PacketType = 5
	// PacketType0RTT is the packet type of a 0-RTT packet
	PacketType0RTT PacketType = 6
)

// A ConnectionID in QUIC
type ConnectionID uint64

// A StreamID in QUIC
type StreamID uint32

// A ByteCount in QUIC
type ByteCount uint64

// MaxByteCount is the maximum value of a ByteCount
const MaxByteCount = ByteCount(math.MaxUint64)

// MaxReceivePacketSize maximum packet size of any QUIC packet, based on
// ethernet's max size, minus the IP and UDP headers. IPv6 has a 40 byte header,
// UDP adds an additional 8 bytes.  This is a total overhead of 48 bytes.
// Ethernet's max packet size is 1500 bytes,  1500 - 48 = 1452.
const MaxReceivePacketSize ByteCount = 1452

// DefaultTCPMSS is the default maximum packet size used in the Linux TCP implementation.
// Used in QUIC for congestion window computations in bytes.
const DefaultTCPMSS ByteCount = 1460

// ClientHelloMinimumSize is the minimum size the server expects an inchoate CHLO to have.
const ClientHelloMinimumSize = 1024

// MaxClientHellos is the maximum number of times we'll send a client hello
// The value 3 accounts for:
// * one failure due to an incorrect or missing source-address token
// * one failure due the server's certificate chain being unavailible and the server being unwilling to send it without a valid source-address token
const MaxClientHellos = 3
