package handshake

import (
	"math/rand"
	"time"

	"github.com/lucas-clemente/quic-go/internal/protocol"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("QUIC TLS Extension", func() {
	Context("Client Hello Transport Parameters", func() {
		It("marshals and unmarshals", func() {
			chtp := &clientHelloTransportParameters{
				InitialVersion: 0x123456,
				Parameters: TransportParameters{
					StreamFlowControlWindow: 0x42,
					IdleTimeout:             0x1337 * time.Second,
				},
			}
			chtp2 := &clientHelloTransportParameters{}
			Expect(chtp2.Unmarshal(chtp.Marshal())).To(Succeed())
			Expect(chtp2.InitialVersion).To(Equal(chtp.InitialVersion))
			Expect(chtp2.Parameters.StreamFlowControlWindow).To(Equal(chtp.Parameters.StreamFlowControlWindow))
			Expect(chtp2.Parameters.IdleTimeout).To(Equal(chtp.Parameters.IdleTimeout))
		})

		It("fuzzes", func() {
			rand := rand.New(rand.NewSource(GinkgoRandomSeed()))
			b := make([]byte, 100)
			for i := 0; i < 1000; i++ {
				rand.Read(b)
				chtp := &clientHelloTransportParameters{}
				chtp.Unmarshal(b[:int(rand.Int31n(100))])
			}
		})
	})

	Context("Encrypted Extensions Transport Parameters", func() {
		It("marshals and unmarshals", func() {
			eetp := &encryptedExtensionsTransportParameters{
				NegotiatedVersion: 0x123456,
				SupportedVersions: []protocol.VersionNumber{0x42, 0x4242},
				Parameters: TransportParameters{
					StreamFlowControlWindow: 0x42,
					IdleTimeout:             0x1337 * time.Second,
				},
			}
			eetp2 := &encryptedExtensionsTransportParameters{}
			Expect(eetp2.Unmarshal(eetp.Marshal())).To(Succeed())
			Expect(eetp2.NegotiatedVersion).To(Equal(eetp.NegotiatedVersion))
			Expect(eetp2.SupportedVersions).To(Equal(eetp.SupportedVersions))
			Expect(eetp2.Parameters.StreamFlowControlWindow).To(Equal(eetp.Parameters.StreamFlowControlWindow))
			Expect(eetp2.Parameters.IdleTimeout).To(Equal(eetp.Parameters.IdleTimeout))
		})

		It("fuzzes", func() {
			rand := rand.New(rand.NewSource(GinkgoRandomSeed()))
			b := make([]byte, 100)
			for i := 0; i < 1000; i++ {
				rand.Read(b)
				chtp := &encryptedExtensionsTransportParameters{}
				chtp.Unmarshal(b[:int(rand.Int31n(100))])
			}
		})
	})
})
