package wire

import (
	"io"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/quicvarint"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NEW_TOKEN frame", func() {
	Context("parsing", func() {
		It("accepts a sample frame", func() {
			token := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
			data := encodeVarInt(uint64(len(token)))
			data = append(data, token...)
			f, l, err := parseNewTokenFrame(data, protocol.Version1)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(f.Token)).To(Equal(token))
			Expect(l).To(Equal(len(data)))
		})

		It("rejects empty tokens", func() {
			data := encodeVarInt(0)
			_, _, err := parseNewTokenFrame(data, protocol.Version1)
			Expect(err).To(MatchError("token must not be empty"))
		})

		It("errors on EOFs", func() {
			token := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
			data := encodeVarInt(uint64(len(token)))
			data = append(data, token...)
			_, l, err := parseNewTokenFrame(data, protocol.Version1)
			Expect(err).NotTo(HaveOccurred())
			Expect(l).To(Equal(len(data)))
			for i := range data {
				_, _, err := parseNewTokenFrame(data[:i], protocol.Version1)
				Expect(err).To(MatchError(io.EOF))
			}
		})
	})

	Context("writing", func() {
		It("writes a sample frame", func() {
			token := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
			f := &NewTokenFrame{Token: []byte(token)}
			b, err := f.Append(nil, protocol.Version1)
			Expect(err).ToNot(HaveOccurred())
			expected := []byte{newTokenFrameType}
			expected = append(expected, encodeVarInt(uint64(len(token)))...)
			expected = append(expected, token...)
			Expect(b).To(Equal(expected))
		})

		It("has the correct min length", func() {
			frame := &NewTokenFrame{Token: []byte("foobar")}
			Expect(frame.Length(protocol.Version1)).To(BeEquivalentTo(1 + quicvarint.Len(6) + 6))
		})
	})
})
