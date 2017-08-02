package frames

import (
	"bytes"

	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/protocol"
)

// A WindowUpdateFrame in QUIC
type WindowUpdateFrame struct {
	StreamID   protocol.StreamID
	ByteOffset protocol.ByteCount
}

//Write writes a RST_STREAM frame
func (f *WindowUpdateFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	typeByte := uint8(0x04)
	b.WriteByte(typeByte)

	utils.LittleEndian.WriteUint32(b, uint32(f.StreamID))
	utils.LittleEndian.WriteUint64(b, uint64(f.ByteOffset))
	return nil
}

// MinLength of a written frame
func (f *WindowUpdateFrame) MinLength(version protocol.VersionNumber) (protocol.ByteCount, error) {
	return 1 + 4 + 8, nil
}

// ParseWindowUpdateFrame parses a RST_STREAM frame
func ParseWindowUpdateFrame(r *bytes.Reader, version protocol.VersionNumber) (*WindowUpdateFrame, error) {
	frame := &WindowUpdateFrame{}

	// read the TypeByte
	_, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	sid, err := utils.LittleEndian.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	frame.StreamID = protocol.StreamID(sid)

	byteOffset, err := utils.LittleEndian.ReadUint64(r)
	if err != nil {
		return nil, err
	}
	frame.ByteOffset = protocol.ByteCount(byteOffset)

	return frame, nil
}
