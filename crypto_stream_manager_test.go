package quic

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crypto Stream Manager", func() {
	var (
		csm *cryptoStreamManager
		cs  *MockCryptoDataHandler

		initialStream   *MockCryptoStream
		handshakeStream *MockCryptoStream
	)

	BeforeEach(func() {
		initialStream = NewMockCryptoStream(mockCtrl)
		handshakeStream = NewMockCryptoStream(mockCtrl)
		cs = NewMockCryptoDataHandler(mockCtrl)
		csm = newCryptoStreamManager(cs, initialStream, handshakeStream)
	})

	It("passes messages to the initial stream", func() {
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		initialStream.EXPECT().HandleCryptoFrame(cf)
		initialStream.EXPECT().GetCryptoData().Return([]byte("foobar"))
		initialStream.EXPECT().GetCryptoData()
		cs.EXPECT().HandleMessage([]byte("foobar"), protocol.EncryptionInitial)
		encLevelChanged, err := csm.HandleCryptoFrame(cf, protocol.EncryptionInitial)
		Expect(err).ToNot(HaveOccurred())
		Expect(encLevelChanged).To(BeFalse())
	})

	It("passes messages to the handshake stream", func() {
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		handshakeStream.EXPECT().HandleCryptoFrame(cf)
		handshakeStream.EXPECT().GetCryptoData().Return([]byte("foobar"))
		handshakeStream.EXPECT().GetCryptoData()
		cs.EXPECT().HandleMessage([]byte("foobar"), protocol.EncryptionHandshake)
		encLevelChanged, err := csm.HandleCryptoFrame(cf, protocol.EncryptionHandshake)
		Expect(err).ToNot(HaveOccurred())
		Expect(encLevelChanged).To(BeFalse())
	})

	It("doesn't call the message handler, if there's no message", func() {
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		handshakeStream.EXPECT().HandleCryptoFrame(cf)
		handshakeStream.EXPECT().GetCryptoData() // don't return any data to handle
		// don't EXPECT any calls to HandleMessage()
		encLevelChanged, err := csm.HandleCryptoFrame(cf, protocol.EncryptionHandshake)
		Expect(err).ToNot(HaveOccurred())
		Expect(encLevelChanged).To(BeFalse())
	})

	It("processes all messages", func() {
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		handshakeStream.EXPECT().HandleCryptoFrame(cf)
		handshakeStream.EXPECT().GetCryptoData().Return([]byte("foo"))
		handshakeStream.EXPECT().GetCryptoData().Return([]byte("bar"))
		handshakeStream.EXPECT().GetCryptoData()
		cs.EXPECT().HandleMessage([]byte("foo"), protocol.EncryptionHandshake)
		cs.EXPECT().HandleMessage([]byte("bar"), protocol.EncryptionHandshake)
		encLevelChanged, err := csm.HandleCryptoFrame(cf, protocol.EncryptionHandshake)
		Expect(err).ToNot(HaveOccurred())
		Expect(encLevelChanged).To(BeFalse())
	})

	It("finishes the crypto stream, when the crypto setup is done with this encryption level", func() {
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		gomock.InOrder(
			handshakeStream.EXPECT().HandleCryptoFrame(cf),
			handshakeStream.EXPECT().GetCryptoData().Return([]byte("foobar")),
			cs.EXPECT().HandleMessage([]byte("foobar"), protocol.EncryptionHandshake).Return(true),
			handshakeStream.EXPECT().Finish(),
		)
		encLevelChanged, err := csm.HandleCryptoFrame(cf, protocol.EncryptionHandshake)
		Expect(err).ToNot(HaveOccurred())
		Expect(encLevelChanged).To(BeTrue())
	})

	It("returns errors that occur when finishing a stream", func() {
		testErr := errors.New("test error")
		cf := &wire.CryptoFrame{Data: []byte("foobar")}
		gomock.InOrder(
			handshakeStream.EXPECT().HandleCryptoFrame(cf),
			handshakeStream.EXPECT().GetCryptoData().Return([]byte("foobar")),
			cs.EXPECT().HandleMessage([]byte("foobar"), protocol.EncryptionHandshake).Return(true),
			handshakeStream.EXPECT().Finish().Return(testErr),
		)
		_, err := csm.HandleCryptoFrame(cf, protocol.EncryptionHandshake)
		Expect(err).To(MatchError(err))
	})

	It("errors for unknown encryption levels", func() {
		_, err := csm.HandleCryptoFrame(&wire.CryptoFrame{}, protocol.Encryption1RTT)
		Expect(err).To(MatchError("received CRYPTO frame with unexpected encryption level: 1-RTT"))
	})
})
