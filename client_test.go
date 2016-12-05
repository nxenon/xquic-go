package quic

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/qerr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var client *Client
	var session *mockSession

	BeforeEach(func() {
		client = &Client{}
		session = &mockSession{connectionID: 0x1337}
		client.connectionID = 0x1337
		client.session = session
		client.version = protocol.Version36
	})

	startUDPConn := func() {
		var err error
		client.addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:0")
		Expect(err).ToNot(HaveOccurred())
		client.conn, err = net.ListenUDP("udp", client.addr)
		Expect(err).NotTo(HaveOccurred())
	}

	It("errors on invalid public header", func() {
		err := client.handlePacket(nil)
		Expect(err.(*qerr.QuicError).ErrorCode).To(Equal(qerr.InvalidPacketHeader))
	})

	It("errors on large packets", func() {
		err := client.handlePacket(bytes.Repeat([]byte{'a'}, int(protocol.MaxPacketSize)+1))
		Expect(err).To(MatchError(qerr.PacketTooLarge))
	})

	It("closes sessions when Close is called", func() {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		Expect(err).ToNot(HaveOccurred())
		client.conn, err = net.ListenUDP("udp", addr)
		Expect(err).ToNot(HaveOccurred())
		err = client.Close()
		Expect(err).ToNot(HaveOccurred())
		Expect(session.closed).To(BeTrue())
		Expect(session.closeReason).To(BeNil())
	})

	It("creates new sessions with the right parameters", func() {
		startUDPConn()
		client.session = nil
		client.hostname = "hostname"
		err := client.createNewSession()
		Expect(err).ToNot(HaveOccurred())
		Expect(client.session).ToNot(BeNil())
		Expect(client.session.(*Session).connectionID).To(Equal(client.connectionID))
		Expect(client.session.(*Session).version).To(Equal(client.version))
	})

	Context("handling packets", func() {
		It("errors on too large packets", func() {
			err := client.handlePacket(bytes.Repeat([]byte{'f'}, int(protocol.MaxPacketSize+1)))
			Expect(err).To(MatchError(qerr.PacketTooLarge))
		})

		It("handles packets", func(done Done) {
			startUDPConn()
			serverConn, err := net.DialUDP("udp", nil, client.conn.LocalAddr().(*net.UDPAddr))
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				listenErr := client.Listen()
				Expect(listenErr).ToNot(HaveOccurred())
				close(done)
			}()

			Expect(session.packetCount).To(BeZero())
			ph := PublicHeader{
				PacketNumber:    1,
				PacketNumberLen: protocol.PacketNumberLen2,
				ConnectionID:    0x1337,
			}
			b := &bytes.Buffer{}
			err = ph.Write(b, protocol.Version36, protocol.PerspectiveServer)
			Expect(err).ToNot(HaveOccurred())
			_, err = serverConn.Write(b.Bytes())
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() int { return session.packetCount }).Should(Equal(1))
			Expect(session.closed).To(BeFalse())

			err = client.Close()
			Expect(err).ToNot(HaveOccurred())
		})

		It("closes the session when encountering an error while handling a packet", func(done Done) {
			startUDPConn()
			serverConn, err := net.DialUDP("udp", nil, client.conn.LocalAddr().(*net.UDPAddr))
			Expect(err).NotTo(HaveOccurred())

			var listenErr error
			go func() {
				defer GinkgoRecover()
				listenErr = client.Listen()
				Expect(listenErr).To(HaveOccurred())
				close(done)
			}()

			// cause a PacketTooLarge error
			_, err = serverConn.Write(bytes.Repeat([]byte{'f'}, 100))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool { return session.closed }).Should(BeTrue())
			Expect(session.closeReason).To(MatchError(listenErr))

			err = client.Close()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("version negotiation", func() {
		getVersionNegotiation := func(versions []protocol.VersionNumber) []byte {
			oldVersionNegotiationPacket := composeVersionNegotiation(0x1337)
			oldSupportVersionTags := protocol.SupportedVersionsAsTags
			var b bytes.Buffer
			for _, v := range versions {
				s := make([]byte, 4)
				binary.LittleEndian.PutUint32(s, protocol.VersionNumberToTag(v))
				b.Write(s)
			}
			protocol.SupportedVersionsAsTags = b.Bytes()
			packet := composeVersionNegotiation(client.connectionID)
			protocol.SupportedVersionsAsTags = oldSupportVersionTags
			Expect(composeVersionNegotiation(0x1337)).To(Equal(oldVersionNegotiationPacket))
			return packet
		}

		It("changes the version after receiving a version negotiation packet", func() {
			startUDPConn()
			newVersion := protocol.Version35
			Expect(newVersion).ToNot(Equal(client.version))
			Expect(session.packetCount).To(BeZero())
			err := client.handlePacket(getVersionNegotiation([]protocol.VersionNumber{newVersion}))
			Expect(client.version).To(Equal(newVersion))
			// it swapped the sessions
			Expect(client.session).ToNot(Equal(session))
			Expect(err).ToNot(HaveOccurred())
			// it didn't pass the version negoation packet to the session (since it has no payload)
			Expect(session.packetCount).To(BeZero())

			err = client.Close()
			Expect(err).ToNot(HaveOccurred())
		})

		It("errors if no matching version is found", func() {
			err := client.handlePacket(getVersionNegotiation([]protocol.VersionNumber{1}))
			Expect(err).To(MatchError(qerr.VersionNegotiationMismatch))
		})

		It("errors if the server should have accepted the offered version", func() {
			err := client.handlePacket(getVersionNegotiation([]protocol.VersionNumber{client.version}))
			Expect(err).To(MatchError(errInvalidVersionNegotiation))
		})
	})
})
