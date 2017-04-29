package quic

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/qerr"
	"github.com/lucas-clemente/quic-go/utils"
)

type client struct {
	mutex                    sync.Mutex
	connStateChangeOrErrCond sync.Cond
	listenErr                error

	conn     connection
	hostname string

	config    *Config
	connState ConnState

	connectionID protocol.ConnectionID
	version      protocol.VersionNumber

	session packetHandler
}

var (
	errCloseSessionForNewVersion = errors.New("closing session in order to recreate it with a new version")
)

// Dial establishes a new QUIC connection to a server using a net.PacketConn.
// The host parameter is used for SNI.
func Dial(pconn net.PacketConn, remoteAddr net.Addr, host string, config *Config) (Session, error) {
	connID, err := utils.GenerateConnectionID()
	if err != nil {
		return nil, err
	}

	hostname, _, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}

	clientConfig := populateClientConfig(config)
	c := &client{
		conn:         &conn{pconn: pconn, currentAddr: remoteAddr},
		connectionID: connID,
		hostname:     hostname,
		config:       clientConfig,
		version:      clientConfig.Versions[0],
	}

	c.connStateChangeOrErrCond.L = &c.mutex

	err = c.createNewSession(nil)
	if err != nil {
		return nil, err
	}

	utils.Infof("Starting new connection to %s (%s), connectionID %x, version %d", hostname, c.conn.RemoteAddr().String(), c.connectionID, c.version)

	return c.establishConnection()
}

func populateClientConfig(config *Config) *Config {
	versions := config.Versions
	if len(versions) == 0 {
		versions = protocol.SupportedVersions
	}

	return &Config{
		TLSConfig: config.TLSConfig,
		ConnState: config.ConnState,
		Versions:  versions,
	}
}

// DialAddr establishes a new QUIC connection to a server.
// The hostname for SNI is taken from the given address.
func DialAddr(addr string, config *Config) (Session, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	return Dial(udpConn, udpAddr, addr, config)
}

func (c *client) establishConnection() (Session, error) {
	go c.listen()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for {
		if c.listenErr != nil {
			return nil, c.listenErr
		}
		if c.config.ConnState != nil && c.connState >= ConnStateVersionNegotiated {
			break
		}
		if c.config.ConnState == nil && c.connState == ConnStateForwardSecure {
			break
		}
		c.connStateChangeOrErrCond.Wait()
	}

	return c.session, nil
}

// Listen listens
func (c *client) listen() {
	var err error

	for {
		var n int
		var addr net.Addr
		data := getPacketBuffer()
		data = data[:protocol.MaxReceivePacketSize]
		// The packet size should not exceed protocol.MaxReceivePacketSize bytes
		// If it does, we only read a truncated packet, which will then end up undecryptable
		n, addr, err = c.conn.Read(data)
		if err != nil {
			if !strings.HasSuffix(err.Error(), "use of closed network connection") {
				c.session.Close(err)
			}
			break
		}
		data = data[:n]

		err = c.handlePacket(addr, data)
		if err != nil {
			utils.Errorf("error handling packet: %s", err.Error())
			c.session.Close(err)
			break
		}
	}

	c.mutex.Lock()
	c.listenErr = err
	c.connStateChangeOrErrCond.Signal()
	c.mutex.Unlock()
}

func (c *client) handlePacket(remoteAddr net.Addr, packet []byte) error {
	rcvTime := time.Now()

	r := bytes.NewReader(packet)
	hdr, err := ParsePublicHeader(r, protocol.PerspectiveServer)
	if err != nil {
		return qerr.Error(qerr.InvalidPacketHeader, err.Error())
	}
	hdr.Raw = packet[:len(packet)-r.Len()]

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// ignore delayed / duplicated version negotiation packets
	if c.connState >= ConnStateVersionNegotiated && hdr.VersionFlag {
		return nil
	}

	// this is the first packet after the client sent a packet with the VersionFlag set
	// if the server doesn't send a version negotiation packet, it supports the suggested version
	if !hdr.VersionFlag && c.connState == ConnStateInitial {
		c.connState = ConnStateVersionNegotiated
		c.connStateChangeOrErrCond.Signal()
		if c.config.ConnState != nil {
			go c.config.ConnState(c.session, ConnStateVersionNegotiated)
		}
	}

	if hdr.VersionFlag {
		// version negotiation packets have no payload
		return c.handlePacketWithVersionFlag(hdr)
	}

	c.session.handlePacket(&receivedPacket{
		remoteAddr:   remoteAddr,
		publicHeader: hdr,
		data:         packet[len(packet)-r.Len():],
		rcvTime:      rcvTime,
	})
	return nil
}

func (c *client) handlePacketWithVersionFlag(hdr *PublicHeader) error {
	for _, v := range hdr.SupportedVersions {
		if v == c.version {
			// the version negotiation packet contains the version that we offered
			// this might be a packet sent by an attacker (or by a terribly broken server implementation)
			// ignore it
			return nil
		}
	}

	ok, highestSupportedVersion := protocol.HighestSupportedVersion(c.config.Versions, hdr.SupportedVersions)
	if !ok {
		return qerr.InvalidVersion
	}

	// switch to negotiated version
	c.version = highestSupportedVersion
	c.connState = ConnStateVersionNegotiated
	var err error
	c.connectionID, err = utils.GenerateConnectionID()
	if err != nil {
		return err
	}
	utils.Infof("Switching to QUIC version %d. New connection ID: %x", highestSupportedVersion, c.connectionID)

	c.session.Close(errCloseSessionForNewVersion)
	err = c.createNewSession(hdr.SupportedVersions)
	if err != nil {
		return err
	}
	if c.config.ConnState != nil {
		go c.config.ConnState(c.session, ConnStateVersionNegotiated)
	}

	return nil
}

func (c *client) cryptoChangeCallback(_ Session, isForwardSecure bool) {
	var state ConnState
	if isForwardSecure {
		state = ConnStateForwardSecure
	} else {
		state = ConnStateSecure
	}

	c.mutex.Lock()
	c.connState = state
	c.connStateChangeOrErrCond.Signal()
	c.mutex.Unlock()

	if c.config.ConnState != nil {
		go c.config.ConnState(c.session, state)
	}
}

func (c *client) createNewSession(negotiatedVersions []protocol.VersionNumber) error {
	var err error
	c.session, err = newClientSession(
		c.conn,
		c.hostname,
		c.version,
		c.connectionID,
		c.config.TLSConfig,
		c.cryptoChangeCallback,
		negotiatedVersions,
	)
	if err != nil {
		return err
	}

	go func() {
		// session.run() returns as soon as the session is closed
		err := c.session.run()
		if err == errCloseSessionForNewVersion {
			return
		}

		c.mutex.Lock()
		c.listenErr = err
		c.connStateChangeOrErrCond.Signal()
		c.mutex.Unlock()

		utils.Infof("Connection %x closed.", c.connectionID)
		c.conn.Close()
	}()
	return nil
}
