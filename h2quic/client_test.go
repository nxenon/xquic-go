package h2quic

import (
	"bytes"
	"compress/gzip"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/qerr"
	"github.com/lucas-clemente/quic-go/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockQuicClient struct {
	streams  map[protocol.StreamID]*mockStream
	closeErr error
}

func (m *mockQuicClient) Close(e error) error { m.closeErr = e; return nil }
func (m *mockQuicClient) Listen() error       { panic("not implemented") }
func (m *mockQuicClient) OpenStream(id protocol.StreamID) (utils.Stream, error) {
	_, ok := m.streams[id]
	if ok {
		panic("Stream already exists")
	}
	ms := &mockStream{id: id}
	m.streams[id] = ms
	return ms, nil
}

func newMockQuicClient() *mockQuicClient {
	return &mockQuicClient{
		streams: make(map[protocol.StreamID]*mockStream),
	}
}

var _ quicClient = &mockQuicClient{}

var _ = Describe("Client", func() {
	var (
		client        *Client
		qClient       *mockQuicClient
		headerStream  *mockStream
		quicTransport *QuicRoundTripper
	)

	BeforeEach(func() {
		var err error
		quicTransport = &QuicRoundTripper{}
		hostname := "quic.clemente.io:1337"
		client, err = NewClient(quicTransport, hostname)
		Expect(err).ToNot(HaveOccurred())
		Expect(client.hostname).To(Equal(hostname))
		qClient = newMockQuicClient()
		client.client = qClient

		headerStream = &mockStream{}
		qClient.streams[3] = headerStream
		client.headerStream = headerStream
		client.requestWriter = newRequestWriter(headerStream)
	})

	It("adds the port to the hostname, if none is given", func() {
		var err error
		client, err = NewClient(quicTransport, "quic.clemente.io")
		Expect(err).ToNot(HaveOccurred())
		Expect(client.hostname).To(Equal("quic.clemente.io:443"))
	})

	It("opens the header stream only after the version has been negotiated", func() {
		// delete the headerStream openend in the BeforeEach
		client.headerStream = nil
		delete(qClient.streams, 3)
		Expect(client.headerStream).To(BeNil()) // header stream not yet opened
		// now start the actual test
		err := client.versionNegotiateCallback()
		Expect(err).ToNot(HaveOccurred())
		Expect(client.headerStream).ToNot(BeNil())
		Expect(client.headerStream.StreamID()).To(Equal(protocol.StreamID(3)))
	})

	It("sets the correct crypto level", func() {
		Expect(client.encryptionLevel).To(Equal(protocol.Unencrypted))
		client.cryptoChangeCallback(false)
		Expect(client.encryptionLevel).To(Equal(protocol.EncryptionSecure))
		client.cryptoChangeCallback(true)
		Expect(client.encryptionLevel).To(Equal(protocol.EncryptionForwardSecure))
	})

	Context("Doing requests", func() {
		var request *http.Request

		getRequest := func(data []byte) *http2.MetaHeadersFrame {
			r := bytes.NewReader(data)
			decoder := hpack.NewDecoder(4096, func(hf hpack.HeaderField) {})
			h2framer := http2.NewFramer(nil, r)
			frame, err := h2framer.ReadFrame()
			Expect(err).ToNot(HaveOccurred())
			mhframe := &http2.MetaHeadersFrame{HeadersFrame: frame.(*http2.HeadersFrame)}
			mhframe.Fields, err = decoder.DecodeFull(mhframe.HeadersFrame.HeaderBlockFragment())
			Expect(err).ToNot(HaveOccurred())
			return mhframe
		}

		getHeaderFields := func(f *http2.MetaHeadersFrame) map[string]string {
			fields := make(map[string]string)
			for _, hf := range f.Fields {
				fields[hf.Name] = hf.Value
			}
			return fields
		}

		BeforeEach(func() {
			var err error
			client.encryptionLevel = protocol.EncryptionForwardSecure
			request, err = http.NewRequest("https", "https://quic.clemente.io:1337/file1.dat", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("does a request", func(done Done) {
			var doRsp *http.Response
			var doErr error
			var doReturned bool
			go func() {
				doRsp, doErr = client.Do(request)
				doReturned = true
			}()

			Eventually(func() []byte { return headerStream.dataWritten.Bytes() }).ShouldNot(BeEmpty())
			Expect(client.highestOpenedStream).To(Equal(protocol.StreamID(5)))
			Expect(qClient.streams).Should(HaveKey(protocol.StreamID(5)))
			Expect(client.responses).To(HaveKey(protocol.StreamID(5)))
			rsp := &http.Response{
				Status:     "418 I'm a teapot",
				StatusCode: 418,
			}
			client.responses[5] <- rsp
			Eventually(func() bool { return doReturned }).Should(BeTrue())
			Expect(doErr).ToNot(HaveOccurred())
			Expect(doRsp).To(Equal(rsp))
			Expect(doRsp.Body).ToNot(BeNil())
			Expect(doRsp.ContentLength).To(BeEquivalentTo(-1))
			Expect(doRsp.Request).To(Equal(request))
			close(done)
		})

		It("closes the quic client when encountering an error on the header stream", func() {
			headerStream.dataToRead.Write([]byte("invalid response"))
			go client.handleHeaderStream()

			var doRsp *http.Response
			var doErr error
			var doReturned bool
			go func() {
				doRsp, doErr = client.Do(request)
				doReturned = true
			}()

			Eventually(func() bool { return doReturned }).Should(BeTrue())
			Expect(client.headerErr).To(HaveOccurred())
			Expect(doErr).To(MatchError(client.headerErr))
			Expect(doRsp).To(BeNil())
			Expect(client.client.(*mockQuicClient).closeErr).To(MatchError(client.headerErr))
		})

		Context("validating the address", func() {
			It("refuses to do requests for the wrong host", func() {
				req, err := http.NewRequest("https", "https://quic.clemente.io:1336/foobar.html", nil)
				Expect(err).ToNot(HaveOccurred())
				_, err = client.Do(req)
				Expect(err).To(MatchError("h2quic Client BUG: Do called for the wrong client"))
			})

			It("refuses to do plain HTTP requests", func() {
				req, err := http.NewRequest("https", "http://quic.clemente.io:1337/foobar.html", nil)
				Expect(err).ToNot(HaveOccurred())
				_, err = client.Do(req)
				Expect(err).To(MatchError("quic http2: unsupported scheme"))
			})

			It("adds the port for request URLs without one", func(done Done) {
				var err error
				client, err = NewClient(quicTransport, "quic.clemente.io")
				Expect(err).ToNot(HaveOccurred())
				req, err := http.NewRequest("https", "https://quic.clemente.io/foobar.html", nil)
				Expect(err).ToNot(HaveOccurred())

				var doErr error
				var doReturned bool
				// the client.Do will block, because the encryption level is still set to Unencrypted
				go func() {
					_, doErr = client.Do(req)
					doReturned = true
				}()

				Consistently(doReturned).Should(BeFalse())
				Expect(doErr).ToNot(HaveOccurred())
				close(done)
			})
		})

		Context("gzip compression", func() {
			var gzippedData []byte // a gzipped foobar
			var response *http.Response

			BeforeEach(func() {
				var b bytes.Buffer
				w := gzip.NewWriter(&b)
				w.Write([]byte("foobar"))
				w.Close()
				gzippedData = b.Bytes()
				response = &http.Response{
					StatusCode: 200,
					Header:     http.Header{"Content-Length": []string{"1000"}},
				}
			})

			It("adds the gzip header to requests", func() {
				var doRsp *http.Response
				var doErr error
				go func() { doRsp, doErr = client.Do(request) }()

				Eventually(func() chan *http.Response { return client.responses[5] }).ShouldNot(BeNil())
				qClient.streams[5].dataToRead.Write(gzippedData)
				response.Header.Add("Content-Encoding", "gzip")
				client.responses[5] <- response
				Eventually(func() *http.Response { return doRsp }).ShouldNot(BeNil())
				Expect(doErr).ToNot(HaveOccurred())
				headers := getHeaderFields(getRequest(headerStream.dataWritten.Bytes()))
				Expect(headers).To(HaveKeyWithValue("accept-encoding", "gzip"))
				Expect(doRsp.ContentLength).To(BeEquivalentTo(-1))
				Expect(doRsp.Header.Get("Content-Encoding")).To(BeEmpty())
				Expect(doRsp.Header.Get("Content-Length")).To(BeEmpty())
				data := make([]byte, 6)
				doRsp.Body.Read(data)
				Expect(data).To(Equal([]byte("foobar")))
			})

			It("doesn't add gzip if the header disable it", func() {
				quicTransport.DisableCompression = true
				var doRsp *http.Response
				var doErr error
				go func() { doRsp, doErr = client.Do(request) }()

				Eventually(func() chan *http.Response { return client.responses[5] }).ShouldNot(BeNil())
				Expect(doErr).ToNot(HaveOccurred())
				headers := getHeaderFields(getRequest(headerStream.dataWritten.Bytes()))
				Expect(headers).ToNot(HaveKey("accept-encoding"))
			})

			It("only decompresses the response if the response contains the right content-encoding header", func() {
				var doRsp *http.Response
				var doErr error
				go func() { doRsp, doErr = client.Do(request) }()

				Eventually(func() chan *http.Response { return client.responses[5] }).ShouldNot(BeNil())
				qClient.streams[5].dataToRead.Write([]byte("not gzipped"))
				client.responses[5] <- response
				Eventually(func() *http.Response { return doRsp }).ShouldNot(BeNil())
				Expect(doErr).ToNot(HaveOccurred())
				headers := getHeaderFields(getRequest(headerStream.dataWritten.Bytes()))
				Expect(headers).To(HaveKeyWithValue("accept-encoding", "gzip"))
				data := make([]byte, 11)
				doRsp.Body.Read(data)
				Expect(doRsp.ContentLength).ToNot(BeEquivalentTo(-1))
				Expect(data).To(Equal([]byte("not gzipped")))
			})

			It("doesn't add the gzip header for requests that have the accept-enconding set", func() {
				request.Header.Add("accept-encoding", "gzip")
				var doRsp *http.Response
				var doErr error
				go func() { doRsp, doErr = client.Do(request) }()

				Eventually(func() chan *http.Response { return client.responses[5] }).ShouldNot(BeNil())
				qClient.streams[5].dataToRead.Write([]byte("gzipped data"))
				client.responses[5] <- response
				Eventually(func() *http.Response { return doRsp }).ShouldNot(BeNil())
				Expect(doErr).ToNot(HaveOccurred())
				headers := getHeaderFields(getRequest(headerStream.dataWritten.Bytes()))
				Expect(headers).To(HaveKeyWithValue("accept-encoding", "gzip"))
				data := make([]byte, 12)
				doRsp.Body.Read(data)
				Expect(doRsp.ContentLength).ToNot(BeEquivalentTo(-1))
				Expect(data).To(Equal([]byte("gzipped data")))
			})
		})

		Context("handling the header stream", func() {
			var h2framer *http2.Framer

			BeforeEach(func() {
				h2framer = http2.NewFramer(&headerStream.dataToRead, nil)
				client.responses[23] = make(chan *http.Response)
			})

			It("reads header values from a response", func() {
				// Taken from https://http2.github.io/http2-spec/compression.html#request.examples.with.huffman.coding
				data := []byte{0x48, 0x03, 0x33, 0x30, 0x32, 0x58, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x61, 0x1d, 0x4d, 0x6f, 0x6e, 0x2c, 0x20, 0x32, 0x31, 0x20, 0x4f, 0x63, 0x74, 0x20, 0x32, 0x30, 0x31, 0x33, 0x20, 0x32, 0x30, 0x3a, 0x31, 0x33, 0x3a, 0x32, 0x31, 0x20, 0x47, 0x4d, 0x54, 0x6e, 0x17, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x77, 0x77, 0x77, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}
				headerStream.dataToRead.Write([]byte{0x0, 0x0, byte(len(data)), 0x1, 0x5, 0x0, 0x0, 0x0, 23})
				headerStream.dataToRead.Write(data)
				go client.handleHeaderStream()
				var rsp *http.Response
				Eventually(client.responses[23]).Should(Receive(&rsp))
				Expect(rsp).ToNot(BeNil())
				Expect(rsp.Proto).To(Equal("HTTP/2.0"))
				Expect(rsp.ProtoMajor).To(BeEquivalentTo(2))
				Expect(rsp.StatusCode).To(BeEquivalentTo(302))
				Expect(rsp.Status).To(Equal("302 Found"))
				Expect(rsp.Header).To(HaveKeyWithValue("Location", []string{"https://www.example.com"}))
				Expect(rsp.Header).To(HaveKeyWithValue("Cache-Control", []string{"private"}))
			})

			It("errors if the H2 frame is not a HeadersFrame", func() {
				var handlerReturned bool
				go func() {
					client.handleHeaderStream()
					handlerReturned = true
				}()

				h2framer.WritePing(true, [8]byte{0, 0, 0, 0, 0, 0, 0, 0})
				var rsp *http.Response
				Eventually(client.responses[23]).Should(Receive(&rsp))
				Expect(rsp).To(BeNil())
				Expect(client.headerErr).To(MatchError(qerr.Error(qerr.InvalidHeadersStreamData, "not a headers frame")))
				Eventually(func() bool { return handlerReturned }).Should(BeTrue())
			})

			It("errors if it can't read the HPACK encoded header fields", func() {
				var handlerReturned bool
				go func() {
					client.handleHeaderStream()
					handlerReturned = true
				}()

				h2framer.WriteHeaders(http2.HeadersFrameParam{
					StreamID:      23,
					EndHeaders:    true,
					BlockFragment: []byte("invalid HPACK data"),
				})

				var rsp *http.Response
				Eventually(client.responses[23]).Should(Receive(&rsp))
				Expect(rsp).To(BeNil())
				Expect(client.headerErr).To(MatchError(qerr.Error(qerr.InvalidHeadersStreamData, "cannot read header fields")))
				Eventually(func() bool { return handlerReturned }).Should(BeTrue())
			})
		})
	})
})
