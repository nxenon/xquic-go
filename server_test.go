package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"reflect"
	"time"

	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/qerr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockSession struct {
	runner         sessionRunner
	connectionID   protocol.ConnectionID
	handledPackets []*receivedPacket
	closed         bool
	closeReason    error
	closedRemote   bool
	stopRunLoop    chan struct{} // run returns as soon as this channel receives a value
}

func (s *mockSession) handlePacket(p *receivedPacket) {
	s.handledPackets = append(s.handledPackets, p)
}

func (s *mockSession) run() error {
	<-s.stopRunLoop
	return s.closeReason
}
func (s *mockSession) Close(e error) error {
	if s.closed {
		return nil
	}
	s.closeReason = e
	s.closed = true
	close(s.stopRunLoop)
	return nil
}
func (s *mockSession) closeRemote(e error) {
	s.closeReason = e
	s.closed = true
	s.closedRemote = true
	close(s.stopRunLoop)
}
func (s *mockSession) OpenStream() (Stream, error) {
	return &stream{}, nil
}
func (s *mockSession) AcceptStream() (Stream, error)           { panic("not implemented") }
func (s *mockSession) AcceptUniStream() (ReceiveStream, error) { panic("not implemented") }
func (s *mockSession) OpenStreamSync() (Stream, error)         { panic("not implemented") }
func (s *mockSession) OpenUniStream() (SendStream, error)      { panic("not implemented") }
func (s *mockSession) OpenUniStreamSync() (SendStream, error)  { panic("not implemented") }
func (s *mockSession) LocalAddr() net.Addr                     { panic("not implemented") }
func (s *mockSession) RemoteAddr() net.Addr                    { panic("not implemented") }
func (*mockSession) Context() context.Context                  { panic("not implemented") }
func (*mockSession) ConnectionState() ConnectionState          { panic("not implemented") }
func (*mockSession) GetVersion() protocol.VersionNumber        { return protocol.VersionWhatever }
func (*mockSession) getCryptoStream() cryptoStreamI            { panic("not implemented") }

var _ Session = &mockSession{}

func newMockSession(
	_ connection,
	runner sessionRunner,
	_ protocol.VersionNumber,
	connectionID protocol.ConnectionID,
	_ *handshake.ServerConfig,
	_ *tls.Config,
	_ *Config,
	_ utils.Logger,
) (packetHandler, error) {
	s := mockSession{
		runner:       runner,
		connectionID: connectionID,
		stopRunLoop:  make(chan struct{}),
	}
	return &s, nil
}

var _ = Describe("Server", func() {
	var (
		conn    *mockPacketConn
		config  *Config
		udpAddr = &net.UDPAddr{IP: net.IPv4(192, 168, 100, 200), Port: 1337}
	)

	BeforeEach(func() {
		conn = newMockPacketConn()
		conn.addr = &net.UDPAddr{}
		config = &Config{Versions: protocol.SupportedVersions}
	})

	Context("with mock session", func() {
		var (
			serv        *server
			firstPacket []byte // a valid first packet for a new connection with connectionID 0x4cfa9f9b668619f6 (= connID)
			connID      = protocol.ConnectionID{0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6}
		)

		BeforeEach(func() {
			serv = &server{
				sessions:     make(map[string]packetHandler),
				newSession:   newMockSession,
				conn:         conn,
				config:       config,
				sessionQueue: make(chan Session, 5),
				errorChan:    make(chan struct{}),
				logger:       utils.DefaultLogger,
			}
			b := &bytes.Buffer{}
			utils.BigEndian.WriteUint32(b, uint32(protocol.SupportedVersions[0]))
			firstPacket = []byte{0x09, 0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6}
			firstPacket = append(append(firstPacket, b.Bytes()...), 0x01)
			firstPacket = append(firstPacket, bytes.Repeat([]byte{0}, protocol.MinClientHelloSize)...) // add padding
		})

		It("setups with the right values", func() {
			config := &Config{
				HandshakeTimeout:            1337 * time.Minute,
				IdleTimeout:                 42 * time.Hour,
				RequestConnectionIDOmission: true,
				MaxIncomingStreams:          1234,
				MaxIncomingUniStreams:       4321,
			}
			c := populateServerConfig(config)
			Expect(c.HandshakeTimeout).To(Equal(1337 * time.Minute))
			Expect(c.IdleTimeout).To(Equal(42 * time.Hour))
			Expect(c.RequestConnectionIDOmission).To(BeFalse())
			Expect(c.MaxIncomingStreams).To(Equal(1234))
			Expect(c.MaxIncomingUniStreams).To(Equal(4321))
		})

		It("disables bidirectional streams", func() {
			config := &Config{
				MaxIncomingStreams:    -1,
				MaxIncomingUniStreams: 4321,
			}
			c := populateServerConfig(config)
			Expect(c.MaxIncomingStreams).To(BeZero())
			Expect(c.MaxIncomingUniStreams).To(Equal(4321))
		})

		It("disables unidirectional streams", func() {
			config := &Config{
				MaxIncomingStreams:    1234,
				MaxIncomingUniStreams: -1,
			}
			c := populateServerConfig(config)
			Expect(c.MaxIncomingStreams).To(Equal(1234))
			Expect(c.MaxIncomingUniStreams).To(BeZero())
		})

		It("returns the address", func() {
			conn.addr = &net.UDPAddr{
				IP:   net.IPv4(192, 168, 13, 37),
				Port: 1234,
			}
			Expect(serv.Addr().String()).To(Equal("192.168.13.37:1234"))
		})

		It("creates new sessions", func() {
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			sess := serv.sessions[string(connID)].(*mockSession)
			Expect(sess.connectionID).To(Equal(connID))
			Expect(sess.handledPackets).To(HaveLen(1))
		})

		It("accepts new TLS sessions", func() {
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			sess, err := newMockSession(nil, nil, protocol.VersionTLS, connID, nil, nil, nil, nil)
			Expect(err).ToNot(HaveOccurred())
			err = serv.setupTLS()
			Expect(err).ToNot(HaveOccurred())
			serv.serverTLS.sessionChan <- tlsSession{
				connID: connID,
				sess:   sess,
			}
			Eventually(func() packetHandler {
				serv.sessionsMutex.Lock()
				defer serv.sessionsMutex.Unlock()
				return serv.sessions[string(connID)]
			}).Should(Equal(sess))
		})

		It("only accepts one new TLS sessions for one connection ID", func() {
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			sess1, err := newMockSession(nil, nil, protocol.VersionTLS, connID, nil, nil, nil, nil)
			Expect(err).ToNot(HaveOccurred())
			sess2, err := newMockSession(nil, nil, protocol.VersionTLS, connID, nil, nil, nil, nil)
			Expect(err).ToNot(HaveOccurred())
			err = serv.setupTLS()
			Expect(err).ToNot(HaveOccurred())
			serv.serverTLS.sessionChan <- tlsSession{
				connID: connID,
				sess:   sess1,
			}
			Eventually(func() packetHandler {
				serv.sessionsMutex.Lock()
				defer serv.sessionsMutex.Unlock()
				return serv.sessions[string(connID)]
			}).Should(Equal(sess1))
			serv.serverTLS.sessionChan <- tlsSession{
				connID: connID,
				sess:   sess2,
			}
			Eventually(func() packetHandler {
				serv.sessionsMutex.Lock()
				defer serv.sessionsMutex.Unlock()
				return serv.sessions[string(connID)]
			}).Should(Equal(sess1))
		})

		It("accepts a session once the connection it is forward secure", func() {
			var acceptedSess Session
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				var err error
				acceptedSess, err = serv.Accept()
				Expect(err).ToNot(HaveOccurred())
				close(done)
			}()
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			sess := serv.sessions[string(connID)].(*mockSession)
			Consistently(func() Session { return acceptedSess }).Should(BeNil())
			serv.sessionQueue <- sess
			Eventually(func() Session { return acceptedSess }).Should(Equal(sess))
			Eventually(done).Should(BeClosed())
		})

		It("doesn't accept sessions that error during the handshake", func() {
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				serv.Accept()
				close(done)
			}()
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			sess := serv.sessions[string(connID)].(*mockSession)
			sess.closeReason = errors.New("handshake failed")
			close(sess.stopRunLoop)
			Consistently(done).ShouldNot(BeClosed())
			// make the go routine return
			serv.removeConnection(connID)
			close(serv.errorChan)
			serv.Close()
			Eventually(done).Should(BeClosed())
		})

		It("assigns packets to existing sessions", func() {
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			err = serv.handlePacket(nil, []byte{0x08, 0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6, 0x01})
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			Expect(serv.sessions[string(connID)].(*mockSession).connectionID).To(Equal(connID))
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(2))
		})

		It("deletes sessions", func() {
			serv.deleteClosedSessionsAfter = time.Second // make sure that the nil value for the closed session doesn't get deleted in this test
			serv.sessions[string(connID)] = &mockSession{}
			serv.removeConnection(connID)
			// The server should now have closed the session, leaving a nil value in the sessions map
			Consistently(func() map[string]packetHandler { return serv.sessions }).Should(HaveLen(1))
			Expect(serv.sessions[string(connID)]).To(BeNil())
		})

		It("deletes nil session entries after a wait time", func() {
			serv.deleteClosedSessionsAfter = 25 * time.Millisecond
			serv.sessions[string(connID)] = &mockSession{}
			// make session.run() return
			serv.removeConnection(connID)
			Eventually(func() bool {
				serv.sessionsMutex.Lock()
				_, ok := serv.sessions[string(connID)]
				serv.sessionsMutex.Unlock()
				return ok
			}).Should(BeFalse())
		})

		It("closes sessions and the connection when Close is called", func() {
			go serv.serve()
			session, _ := newMockSession(nil, nil, 0, connID, nil, nil, nil, nil)
			serv.sessions[string(connID)] = session
			err := serv.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(session.(*mockSession).closed).To(BeTrue())
			Expect(conn.closed).To(BeTrue())
		})

		It("ignores packets for closed sessions", func() {
			serv.sessions[string(connID)] = nil
			err := serv.handlePacket(nil, []byte{0x08, 0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6, 0x01})
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			Expect(serv.sessions[string(connID)]).To(BeNil())
		})

		It("works if no quic.Config is given", func(done Done) {
			ln, err := ListenAddr("127.0.0.1:0", testdata.GetTLSConfig(), nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(ln.Close()).To(Succeed())
			close(done)
		}, 1)

		It("closes properly", func() {
			ln, err := ListenAddr("127.0.0.1:0", testdata.GetTLSConfig(), config)
			Expect(err).ToNot(HaveOccurred())

			var returned bool
			go func() {
				defer GinkgoRecover()
				_, err := ln.Accept()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("use of closed network connection"))
				returned = true
			}()
			ln.Close()
			Eventually(func() bool { return returned }).Should(BeTrue())
		})

		It("errors when encountering a connection error", func(done Done) {
			testErr := errors.New("connection error")
			conn.readErr = testErr
			go serv.serve()
			_, err := serv.Accept()
			Expect(err).To(MatchError(testErr))
			Expect(serv.Close()).To(Succeed())
			close(done)
		}, 0.5)

		It("closes all sessions when encountering a connection error", func() {
			session, _ := newMockSession(nil, nil, 0, connID, nil, nil, nil, nil)
			serv.sessions[string(connID)] = session
			Expect(serv.sessions[string(connID)].(*mockSession).closed).To(BeFalse())
			testErr := errors.New("connection error")
			conn.readErr = testErr
			go serv.serve()
			Eventually(func() bool { return session.(*mockSession).closed }).Should(BeTrue())
			Expect(serv.Close()).To(Succeed())
		})

		It("ignores delayed packets with mismatching versions", func() {
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(1))
			b := &bytes.Buffer{}
			// add an unsupported version
			data := []byte{0x09, 0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6}
			utils.BigEndian.WriteUint32(b, uint32(protocol.SupportedVersions[0]+1))
			data = append(append(data, b.Bytes()...), 0x01)
			err = serv.handlePacket(nil, data)
			Expect(err).ToNot(HaveOccurred())
			// if we didn't ignore the packet, the server would try to send a version negotiation packet, which would make the test panic because it doesn't have a udpConn
			Expect(conn.dataWritten.Bytes()).To(BeEmpty())
			// make sure the packet was *not* passed to session.handlePacket()
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(1))
		})

		It("errors on invalid public header", func() {
			err := serv.handlePacket(nil, nil)
			Expect(err.(*qerr.QuicError).ErrorCode).To(Equal(qerr.InvalidPacketHeader))
		})

		It("errors on packets that are smaller than the Payload Length in the packet header", func() {
			serv.supportsTLS = true
			b := &bytes.Buffer{}
			hdr := &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeHandshake,
				PayloadLen:       1000,
				SrcConnectionID:  protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8},
				DestConnectionID: protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8},
				Version:          versionIETFFrames,
			}
			Expect(hdr.Write(b, protocol.PerspectiveClient, versionIETFFrames)).To(Succeed())
			err := serv.handlePacket(nil, append(b.Bytes(), make([]byte, 456)...))
			Expect(err).To(MatchError("packet payload (456 bytes) is smaller than the expected payload length (1000 bytes)"))
		})

		It("cuts packets at the payload length", func() {
			serv.supportsTLS = true
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			b := &bytes.Buffer{}
			hdr := &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeHandshake,
				PayloadLen:       123,
				SrcConnectionID:  connID,
				DestConnectionID: connID,
				Version:          versionIETFFrames,
			}
			Expect(hdr.Write(b, protocol.PerspectiveClient, versionIETFFrames)).To(Succeed())
			err = serv.handlePacket(nil, append(b.Bytes(), make([]byte, 456)...))
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(2))
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets[1].data).To(HaveLen(123))
		})

		It("drops packets with invalid packet types", func() {
			serv.supportsTLS = true
			b := &bytes.Buffer{}
			hdr := &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeRetry,
				PayloadLen:       123,
				SrcConnectionID:  connID,
				DestConnectionID: connID,
				Version:          versionIETFFrames,
			}
			Expect(hdr.Write(b, protocol.PerspectiveClient, versionIETFFrames)).To(Succeed())
			err := serv.handlePacket(nil, append(b.Bytes(), make([]byte, 456)...))
			Expect(err).To(MatchError("Received unsupported packet type: Retry"))
		})

		It("ignores Public Resets", func() {
			err := serv.handlePacket(nil, firstPacket)
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(1))
			err = serv.handlePacket(nil, wire.WritePublicReset(connID, 1, 1337))
			Expect(err).ToNot(HaveOccurred())
			Expect(serv.sessions).To(HaveLen(1))
			Expect(serv.sessions[string(connID)].(*mockSession).handledPackets).To(HaveLen(1))
		})

		It("doesn't try to process a packet after sending a gQUIC Version Negotiation Packet", func() {
			config.Versions = []protocol.VersionNumber{99}
			b := &bytes.Buffer{}
			hdr := wire.Header{
				VersionFlag:      true,
				DestConnectionID: connID,
				SrcConnectionID:  connID,
				PacketNumber:     1,
				PacketNumberLen:  protocol.PacketNumberLen2,
			}
			hdr.Write(b, protocol.PerspectiveClient, 13 /* not a valid QUIC version */)
			b.Write(bytes.Repeat([]byte{0}, protocol.MinClientHelloSize)) // add a fake CHLO
			serv.conn = conn
			err := serv.handlePacket(nil, b.Bytes())
			Expect(conn.dataWritten.Bytes()).ToNot(BeEmpty())
			Expect(err).ToNot(HaveOccurred())
		})

		It("doesn't respond with a version negotiation packet if the first packet is too small", func() {
			b := &bytes.Buffer{}
			hdr := wire.Header{
				VersionFlag:      true,
				DestConnectionID: connID,
				SrcConnectionID:  connID,
				PacketNumber:     1,
				PacketNumberLen:  protocol.PacketNumberLen2,
			}
			hdr.Write(b, protocol.PerspectiveClient, 13 /* not a valid QUIC version */)
			b.Write(bytes.Repeat([]byte{0}, protocol.MinClientHelloSize-1)) // this packet is 1 byte too small
			serv.conn = conn
			err := serv.handlePacket(udpAddr, b.Bytes())
			Expect(err).To(MatchError("dropping small packet with unknown version"))
			Expect(conn.dataWritten.Len()).Should(BeZero())
		})
	})

	It("setups with the right values", func() {
		supportedVersions := []protocol.VersionNumber{protocol.VersionTLS, protocol.Version39}
		acceptCookie := func(_ net.Addr, _ *Cookie) bool { return true }
		config := Config{
			Versions:         supportedVersions,
			AcceptCookie:     acceptCookie,
			HandshakeTimeout: 1337 * time.Hour,
			IdleTimeout:      42 * time.Minute,
			KeepAlive:        true,
		}
		ln, err := Listen(conn, &tls.Config{}, &config)
		Expect(err).ToNot(HaveOccurred())
		server := ln.(*server)
		Expect(server.deleteClosedSessionsAfter).To(Equal(protocol.ClosedSessionDeleteTimeout))
		Expect(server.sessions).ToNot(BeNil())
		Expect(server.scfg).ToNot(BeNil())
		Expect(server.config.Versions).To(Equal(supportedVersions))
		Expect(server.config.HandshakeTimeout).To(Equal(1337 * time.Hour))
		Expect(server.config.IdleTimeout).To(Equal(42 * time.Minute))
		Expect(reflect.ValueOf(server.config.AcceptCookie)).To(Equal(reflect.ValueOf(acceptCookie)))
		Expect(server.config.KeepAlive).To(BeTrue())
	})

	It("errors when the Config contains an invalid version", func() {
		version := protocol.VersionNumber(0x1234)
		_, err := Listen(conn, &tls.Config{}, &Config{Versions: []protocol.VersionNumber{version}})
		Expect(err).To(MatchError("0x1234 is not a valid QUIC version"))
	})

	It("fills in default values if options are not set in the Config", func() {
		ln, err := Listen(conn, &tls.Config{}, &Config{})
		Expect(err).ToNot(HaveOccurred())
		server := ln.(*server)
		Expect(server.config.Versions).To(Equal(protocol.SupportedVersions))
		Expect(server.config.HandshakeTimeout).To(Equal(protocol.DefaultHandshakeTimeout))
		Expect(server.config.IdleTimeout).To(Equal(protocol.DefaultIdleTimeout))
		Expect(reflect.ValueOf(server.config.AcceptCookie)).To(Equal(reflect.ValueOf(defaultAcceptCookie)))
		Expect(server.config.KeepAlive).To(BeFalse())
	})

	It("listens on a given address", func() {
		addr := "127.0.0.1:13579"
		ln, err := ListenAddr(addr, nil, config)
		Expect(err).ToNot(HaveOccurred())
		serv := ln.(*server)
		Expect(serv.Addr().String()).To(Equal(addr))
	})

	It("errors if given an invalid address", func() {
		addr := "127.0.0.1"
		_, err := ListenAddr(addr, nil, config)
		Expect(err).To(BeAssignableToTypeOf(&net.AddrError{}))
	})

	It("errors if given an invalid address", func() {
		addr := "1.1.1.1:1111"
		_, err := ListenAddr(addr, nil, config)
		Expect(err).To(BeAssignableToTypeOf(&net.OpError{}))
	})

	It("sends a gQUIC Version Negotaion Packet, if the client sent a gQUIC Public Header", func() {
		connID := protocol.ConnectionID{8, 7, 6, 5, 4, 3, 2, 1}
		b := &bytes.Buffer{}
		hdr := wire.Header{
			VersionFlag:      true,
			DestConnectionID: connID,
			SrcConnectionID:  connID,
			PacketNumber:     1,
			PacketNumberLen:  protocol.PacketNumberLen2,
		}
		hdr.Write(b, protocol.PerspectiveClient, 13 /* not a valid QUIC version */)
		b.Write(bytes.Repeat([]byte{0}, protocol.MinClientHelloSize)) // add a fake CHLO
		conn.dataToRead <- b.Bytes()
		conn.dataReadFrom = udpAddr
		ln, err := Listen(conn, nil, config)
		Expect(err).ToNot(HaveOccurred())

		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			ln.Accept()
			close(done)
		}()

		Eventually(func() int { return conn.dataWritten.Len() }).ShouldNot(BeZero())
		Expect(conn.dataWrittenTo).To(Equal(udpAddr))
		r := bytes.NewReader(conn.dataWritten.Bytes())
		packet, err := wire.ParseHeaderSentByServer(r)
		Expect(err).ToNot(HaveOccurred())
		Expect(packet.VersionFlag).To(BeTrue())
		Expect(packet.DestConnectionID).To(Equal(connID))
		Expect(packet.SrcConnectionID).To(Equal(connID))
		Expect(r.Len()).To(BeZero())
		Consistently(done).ShouldNot(BeClosed())
		// make the go routine return
		ln.Close()
		Eventually(done).Should(BeClosed())
	})

	It("sends an IETF draft style Version Negotaion Packet, if the client sent a IETF draft style header", func() {
		connID := protocol.ConnectionID{8, 7, 6, 5, 4, 3, 2, 1}
		config.Versions = append(config.Versions, protocol.VersionTLS)
		b := &bytes.Buffer{}
		hdr := wire.Header{
			Type:             protocol.PacketTypeInitial,
			IsLongHeader:     true,
			DestConnectionID: connID,
			SrcConnectionID:  connID,
			PacketNumber:     0x55,
			Version:          0x1234,
			PayloadLen:       protocol.MinInitialPacketSize,
		}
		err := hdr.Write(b, protocol.PerspectiveClient, protocol.VersionTLS)
		Expect(err).ToNot(HaveOccurred())
		b.Write(bytes.Repeat([]byte{0}, protocol.MinInitialPacketSize)) // add a fake CHLO
		conn.dataToRead <- b.Bytes()
		conn.dataReadFrom = udpAddr
		ln, err := Listen(conn, testdata.GetTLSConfig(), config)
		Expect(err).ToNot(HaveOccurred())

		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			ln.Accept()
			close(done)
		}()

		Eventually(func() int { return conn.dataWritten.Len() }).ShouldNot(BeZero())
		Expect(conn.dataWrittenTo).To(Equal(udpAddr))
		r := bytes.NewReader(conn.dataWritten.Bytes())
		packet, err := wire.ParseHeaderSentByServer(r)
		Expect(err).ToNot(HaveOccurred())
		Expect(packet.IsVersionNegotiation).To(BeTrue())
		Expect(packet.DestConnectionID).To(Equal(connID))
		Expect(packet.SrcConnectionID).To(Equal(connID))
		Expect(r.Len()).To(BeZero())
		Consistently(done).ShouldNot(BeClosed())
		// make the go routine return
		ln.Close()
		Eventually(done).Should(BeClosed())
	})

	It("ignores IETF draft style Initial packets, if it doesn't support TLS", func() {
		connID := protocol.ConnectionID{8, 7, 6, 5, 4, 3, 2, 1}
		b := &bytes.Buffer{}
		hdr := wire.Header{
			Type:             protocol.PacketTypeInitial,
			IsLongHeader:     true,
			DestConnectionID: connID,
			SrcConnectionID:  connID,
			PacketNumber:     0x55,
			Version:          protocol.VersionTLS,
		}
		err := hdr.Write(b, protocol.PerspectiveClient, protocol.VersionTLS)
		Expect(err).ToNot(HaveOccurred())
		b.Write(bytes.Repeat([]byte{0}, protocol.MinClientHelloSize)) // add a fake CHLO
		conn.dataToRead <- b.Bytes()
		conn.dataReadFrom = udpAddr
		ln, err := Listen(conn, testdata.GetTLSConfig(), config)
		Expect(err).ToNot(HaveOccurred())
		defer ln.Close()
		Consistently(func() int { return conn.dataWritten.Len() }).Should(BeZero())
	})

	It("ignores non-Initial Long Header packets for unknown connections", func() {
		connID := protocol.ConnectionID{8, 7, 6, 5, 4, 3, 2, 1}
		b := &bytes.Buffer{}
		hdr := wire.Header{
			Type:             protocol.PacketTypeHandshake,
			IsLongHeader:     true,
			DestConnectionID: connID,
			SrcConnectionID:  connID,
			PacketNumber:     0x55,
			Version:          protocol.VersionTLS,
		}
		err := hdr.Write(b, protocol.PerspectiveClient, protocol.VersionTLS)
		Expect(err).ToNot(HaveOccurred())
		conn.dataToRead <- b.Bytes()
		conn.dataReadFrom = udpAddr
		ln, err := Listen(conn, testdata.GetTLSConfig(), config)
		Expect(err).ToNot(HaveOccurred())
		defer ln.Close()
		Consistently(func() int { return conn.dataWritten.Len() }).Should(BeZero())
	})

	It("sends a PublicReset for new connections that don't have the VersionFlag set", func() {
		conn.dataReadFrom = udpAddr
		conn.dataToRead <- []byte{0x08, 0x4c, 0xfa, 0x9f, 0x9b, 0x66, 0x86, 0x19, 0xf6, 0x01}
		ln, err := Listen(conn, nil, config)
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			_, err := ln.Accept()
			Expect(err).ToNot(HaveOccurred())
		}()

		Eventually(func() int { return conn.dataWritten.Len() }).ShouldNot(BeZero())
		Expect(conn.dataWrittenTo).To(Equal(udpAddr))
		Expect(conn.dataWritten.Bytes()[0] & 0x02).ToNot(BeZero()) // check that the ResetFlag is set
		Expect(ln.(*server).sessions).To(BeEmpty())
	})
})

var _ = Describe("default source address verification", func() {
	It("accepts a token", func() {
		remoteAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 0, 1)}
		cookie := &Cookie{
			RemoteAddr: "192.168.0.1",
			SentTime:   time.Now().Add(-protocol.CookieExpiryTime).Add(time.Second), // will expire in 1 second
		}
		Expect(defaultAcceptCookie(remoteAddr, cookie)).To(BeTrue())
	})

	It("requests verification if no token is provided", func() {
		remoteAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 0, 1)}
		Expect(defaultAcceptCookie(remoteAddr, nil)).To(BeFalse())
	})

	It("rejects a token if the address doesn't match", func() {
		remoteAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 0, 1)}
		cookie := &Cookie{
			RemoteAddr: "127.0.0.1",
			SentTime:   time.Now(),
		}
		Expect(defaultAcceptCookie(remoteAddr, cookie)).To(BeFalse())
	})

	It("accepts a token for a remote address is not a UDP address", func() {
		remoteAddr := &net.TCPAddr{IP: net.IPv4(192, 168, 0, 1), Port: 1337}
		cookie := &Cookie{
			RemoteAddr: "192.168.0.1:1337",
			SentTime:   time.Now(),
		}
		Expect(defaultAcceptCookie(remoteAddr, cookie)).To(BeTrue())
	})

	It("rejects an invalid token for a remote address is not a UDP address", func() {
		remoteAddr := &net.TCPAddr{IP: net.IPv4(192, 168, 0, 1), Port: 1337}
		cookie := &Cookie{
			RemoteAddr: "192.168.0.1:7331", // mismatching port
			SentTime:   time.Now(),
		}
		Expect(defaultAcceptCookie(remoteAddr, cookie)).To(BeFalse())
	})

	It("rejects an expired token", func() {
		remoteAddr := &net.UDPAddr{IP: net.IPv4(192, 168, 0, 1)}
		cookie := &Cookie{
			RemoteAddr: "192.168.0.1",
			SentTime:   time.Now().Add(-protocol.CookieExpiryTime).Add(-time.Second), // expired 1 second ago
		}
		Expect(defaultAcceptCookie(remoteAddr, cookie)).To(BeFalse())
	})
})
