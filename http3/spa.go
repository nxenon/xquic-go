package http3

import (
	"bytes"
	"context"
	"fmt"
	quic "github.com/nxenon/h3spacex"
	"github.com/quic-go/qpack"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"sync"
)

// Constants for timing headers used by the synchronization method
const (
	ClientTimeHeaderRTTUs      = "X-Client-Time-RTT-Us"
	ClientTimeHeaderPostLastUs = "X-Client-Time-Post-Last-Us"
)

// RequestStreamPair couples a request with its initialized QUIC stream
type RequestStreamPair struct {
	Req *http.Request
	S   quic.Stream
}

func CalculateIntegerEncodingLengthValue(l int64) []byte {
	if l < 0 {
		panic("l is lower that 0!")
	}

	var encodedValue []byte

	switch {
	case l <= 63:
		encodedValue = []byte{byte(l)}
	case l <= 16383:
		encodedValue = []byte{0x40 | byte(l>>8), byte(l)}
	case l <= 1073741823:
		encodedValue = []byte{0x80 | byte(l>>24), byte(l >> 16), byte(l >> 8), byte(l)}
	default:
		encodedValue = []byte{0xC0 | byte(l>>56), byte(l >> 48), byte(l >> 40), byte(l >> 32),
			byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}
	}

	return encodedValue
}

func ParseResponseFromStream(biStream quic.Stream) (*http.Response, error) {

	frame1, err := NewXParseNextFrame(biStream, nil)
	if err != nil {
		return &http.Response{}, err
	}

	hf, ok := frame1.(*HeadersFrame)
	if !ok {
		fmt.Println("first frame is not headers frame")
		return &http.Response{}, error(nil)
		//panic("first frame is not headers frame")
	}

	headerBlock := make([]byte, hf.Length)

	if _, err2 := io.ReadFull(biStream, headerBlock); err2 != nil {
		return &http.Response{}, err2
	}
	decoder := qpack.Decoder{}

	hfs, err := decoder.DecodeFull(headerBlock)

	res, err := ResponseFromHeaders(hfs)

	var httpStr quic.Stream
	hstr := NewStream(biStream, nil)
	if _, ok := res.Header["Content-Length"]; ok && res.ContentLength >= 0 {
		httpStr = NewLengthLimitedStream(hstr, res.ContentLength)
	} else {
		httpStr = hstr
	}
	respBody := NewResponseBody(httpStr, nil, nil)

	// Rules for when to set Content-Length are defined in https://tools.ietf.org/html/rfc7230#section-3.3.2.
	_, hasTransferEncoding := res.Header["Transfer-Encoding"]
	isInformational := res.StatusCode >= 100 && res.StatusCode < 200
	isNoContent := res.StatusCode == http.StatusNoContent
	//isSuccessfulConnect := req.Method == http.MethodConnect && res.StatusCode >= 200 && res.StatusCode < 300
	//if !hasTransferEncoding && !isInformational && !isNoContent && !isSuccessfulConnect {
	if !hasTransferEncoding && !isInformational && !isNoContent {
		res.ContentLength = -1
		if clens, ok := res.Header["Content-Length"]; ok && len(clens) == 1 {
			if clen64, err := strconv.ParseInt(clens[0], 10, 64); err == nil {
				res.ContentLength = clen64
			}
		}
	}
	requestGzip := true
	if requestGzip && res.Header.Get("Content-Encoding") == "gzip" {
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Length")
		res.ContentLength = -1
		res.Body = NewGzipReader(respBody)
		res.Uncompressed = true
	} else {
		res.Body = respBody
	}

	return res, nil
}

func Print_bytes_in_hex(data []byte) {
	for _, n := range data {
		fmt.Printf("%02x", n) // Prints hexadecimal representation
	}
	fmt.Println()
}

func SendRequestBytesInStream(stream quic.Stream, requestBytes []byte) error {

	_, err := stream.Write(requestBytes)
	return err

}

func ReadFromAllStreams(allStreams map[*http.Request]quic.Stream) map[*http.Request]*http.Response {

	streamsResponseMap := make(map[*http.Request]*http.Response)
	for request, biStream := range allStreams {
		res, err := ReadOneStream(biStream)
		if err != nil {
			fmt.Printf("Stream ID: %d has error (no response)!: %s", biStream.StreamID(), err)
			continue
		}
		streamsResponseMap[request] = res
	}

	return streamsResponseMap

}

func ReadOneStream(biStream quic.Stream) (*http.Response, error) {
	res, err := ParseResponseFromStream(biStream)
	return res, err
}

func GetBidirectionalStream(quicConnection quic.Connection) quic.Stream {
	context := context.Background()
	biStream, err := quicConnection.OpenStreamSync(context)
	if err != nil {
		panic(err)
	}

	return biStream

}

func CloseAllStreams(allStreams map[*http.Request]quic.Stream) {
	for key, _ := range allStreams {
		err := allStreams[key].Close()
		if err != nil {
			//fmt.Printf("Error closing Stream ID %d -> %s\n", value.StreamID(), err)
		}
	}
}

func GetDataFrameBytes(req http.Request) []byte {

	length := CalculateIntegerEncodingLengthValue(req.ContentLength)
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	var requestBodyBuffer bytes.Buffer
	_, err := io.Copy(&requestBodyBuffer, req.Body)
	if err != nil {
		panic(err)
	}

	buf = append(buf, requestBodyBuffer.Bytes()...)

	return buf

}

func GetDataFrameBytesWithLengthMinusLastByteNum(req http.Request, lastByteNum int) []byte {

	length := CalculateIntegerEncodingLengthValue(req.ContentLength - int64(lastByteNum))
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	var requestBodyBuffer bytes.Buffer
	_, err := io.Copy(&requestBodyBuffer, req.Body)
	if err != nil {
		panic(err)
	}

	buf = append(buf, requestBodyBuffer.Bytes()...)

	return buf

}

func GetLastByteDataFrame(b []byte) []byte {

	length := CalculateIntegerEncodingLengthValue(int64(len(b)))
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	buf = append(buf, b...)

	return buf

}

func GetLastByteZeroDataFrameForGetRequests() []byte {

	length := []byte{byte(0)}
	buf := []byte{
		0x00,
	}
	buf = append(buf, length...)

	return buf

}

func GetRequestHeadersBytes(req http.Request, setContentLength bool) []byte {
	buf := &bytes.Buffer{}
	requestWriter := NewXRequestWriter(nil)
	isGzipped := true
	err := requestWriter.NewXwriteHeaders(buf, &req, isGzipped, setContentLength)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func GetRequestFinalPayload(req http.Request) []byte {
	setContentLength := true
	var requestDataBytes []byte
	if req.Body != nil {
		//requestDataBytes := newGetDataFrameBytes(req.Body, req.ContentLength)
		requestDataBytes = GetDataFrameBytes(req)
	}

	finalPayload := append(GetRequestHeadersBytes(req, setContentLength), requestDataBytes...)

	return finalPayload
}

func GetRequestObject(urlString string, method string, headersMap map[string]string, bodyData []byte) (http.Request, error) {
	method = strings.ToUpper(method)
	var req *http.Request
	if method == "GET" {
		reqx, err := http.NewRequest(method, urlString, nil)
		if err != nil {
			return http.Request{}, err
		}
		req = reqx
	} else {
		reqx, err2 := http.NewRequest(method, urlString, bytes.NewReader(bodyData))
		if err2 != nil {
			return http.Request{}, err2
		}
		req = reqx
	}
	for key, value := range headersMap {
		req.Header.Set(key, value)
	}
	return *req, nil
}

func GetRequestObjectGetWithBody(urlString string, method string, headersMap map[string]string, bodyData []byte) (http.Request, error) {
	method = strings.ToUpper(method)
	var req *http.Request

	reqx, err2 := http.NewRequest(method, urlString, bytes.NewReader(bodyData))
	if err2 != nil {
		return http.Request{}, err2
	}
	req = reqx

	for key, value := range headersMap {
		req.Header.Set(key, value)
	}
	return *req, nil
}

func SendLastBytesOfStreams(allStreamsWithLastByte map[quic.Stream][]byte) {
	for s, b := range allStreamsWithLastByte {
		err := SendRequestBytesInStream(s, b)
		if err != nil {
			fmt.Printf("Error occurred in sending last byte of Stream ID: %d -> %s\n", s.StreamID(), err)
		}
	}
}

func SendRequestsWithLastFrameSynchronizationMethod(
	quicConn quic.Connection,
	allRequests []*http.Request,
	lastByteNum int,
	sleepMillisecondsBeforeSendingLastByte int,
	setContentLength bool,
) map[*http.Request]*http.Response {

	type streamState struct {
		req       *http.Request
		s         quic.Stream
		lastChunk []byte
		hasLast   bool
	}

	out := make(map[*http.Request]*http.Response)
	var mu sync.Mutex
	var wg sync.WaitGroup

	states := make([]streamState, 0, len(allRequests))
	startPre := make(map[quic.Stream]time.Time) 
	startPost := make(map[quic.Stream]time.Time) 

	// 1) Phase 1: Open all streams and send initial payloads (Batching initial writes)
	for _, request := range allRequests {
		headersFrameByte := GetRequestHeadersBytes(*request, setContentLength)

		var firstPayload []byte
		var finalData []byte
		hasBody := request.Body != nil

		if lastByteNum > 0 && hasBody {
			dataFrameBytes := GetDataFrameBytesWithLengthMinusLastByteNum(*request, lastByteNum)
			allDataBytesExceptLast := dataFrameBytes[:len(dataFrameBytes)-lastByteNum]
			firstPayload = append(headersFrameByte, allDataBytesExceptLast...)

			finalByte := dataFrameBytes[len(dataFrameBytes)-lastByteNum:]
			finalData = GetLastByteDataFrame(finalByte)
		} else if hasBody {
			firstPayload = append(headersFrameByte, GetDataFrameBytes(*request)...)
		} else {
			firstPayload = headersFrameByte
		}

		s := GetBidirectionalStream(quicConn)

		// Record pre-write time sequentially on the main thread
		startPre[s] = time.Now()
		if err := SendRequestBytesInStream(s, firstPayload); err != nil {
			fmt.Printf("Error sending initial bytes on Stream %d: %v\n", s.StreamID(), err)
			continue
		}

		st := streamState{
			req:       request,
			s:         s,
			lastChunk: finalData,
			hasLast:   (lastByteNum > 0 && hasBody),
		}
		states = append(states, st)
	}

	// 2) Phase 2: Hold/Sleep before releasing final synchronization frames (Maintains old logic)
	if lastByteNum > 0 {
		time.Sleep(time.Duration(sleepMillisecondsBeforeSendingLastByte) * time.Millisecond)
	}

	// 3) Phase 3: Send last chunks and close streams (Recording startPost safely on main thread)
	for i := range states {
		st := &states[i]
		if st.hasLast {
			if err := SendRequestBytesInStream(st.s, st.lastChunk); err != nil {
				fmt.Printf("Error sending last byte on Stream %d: %v\n", st.s.StreamID(), err)
				continue
			}
		}
		
		// Send FIN frame to let the server know we are done writing
		if err := st.s.Close(); err != nil {
			fmt.Printf("Error closing (FIN) Stream %d: %v\n", st.s.StreamID(), err)
		}
		startPost[st.s] = time.Now()
	}

	// 4) Phase 4: Spin up concurrent background readers AFTER all network writes are complete.
	for _, st := range states {
		t0 := startPre[st.s]
		t1 := startPost[st.s]

		wg.Add(1)
		go func(req *http.Request, stream quic.Stream, t0, t1 time.Time) {
			defer wg.Done()
			res, err := ReadOneStream(stream)
			if err != nil {
				fmt.Printf("Stream %d read error: %v\n", stream.StreamID(), err)
				return
			}
			
			// Safely read from stack-local copies of timestamps
			res.Header.Set(ClientTimeHeaderRTTUs, strconv.FormatInt(time.Since(t0).Microseconds(), 10))
			res.Header.Set(ClientTimeHeaderPostLastUs, strconv.FormatInt(time.Since(t1).Microseconds(), 10))
			
			mu.Lock()
			out[req] = res
			mu.Unlock()
		}(st.req, st.s, t0, t1) // <--- Values bound safely here
	}

	// 5) Wait for all concurrent readers to finish
	wg.Wait()
	return out
}

func SendRequestsWithoutBodyWithinASinglePacket(
	quicConn quic.Connection,
	allRequests []*http.Request,
) map[*http.Request]*http.Response {

	out := make(map[*http.Request]*http.Response)
	var mu sync.Mutex
	var wg sync.WaitGroup

	allStreams := make(map[*http.Request]quic.Stream)
	allStreamsWithHeadersByte := make(map[quic.Stream][]byte)
	startPre := make(map[quic.Stream]time.Time)

	// 1) Prime all streams and record pre-write timestamps sequentially
	for _, request := range allRequests {
		headersFrameByte := GetRequestHeadersBytes(*request, true)
		biStream := GetBidirectionalStream(quicConn)
		
		// Record the start time right before data goes out
		startPre[biStream] = time.Now() 
		
		allStreamsWithHeadersByte[biStream] = headersFrameByte
		allStreams[request] = biStream 
	}

	// 2) Burst all headers out across the network within a single packet window
	SendLastBytesOfStreams(allStreamsWithHeadersByte) 

	// 3) Close streams (FIN) so the server begins processing the batch
	CloseAllStreams(allStreams)
	
	// Capture the exact moment the entire single-packet burst finished processing
	startPostTime := time.Now()

	// 4) Spawn concurrent readers to get highly precise, per-stream response timings
	for request, stream := range allStreams {
		t0 := startPre[stream]
		t1 := startPostTime // All single-packet streams share this joint completion line

		wg.Add(1)
		// Pass t0 and t1 explicitly to isolate them to the goroutine's stack
		go func(req *http.Request, st quic.Stream, t0, t1 time.Time) {
			defer wg.Done()
			
			res, err := ReadOneStream(st)
			if err != nil {
				fmt.Printf("Stream %d read error: %v\n", st.StreamID(), err)
				return
			}
			
			// Safely stamp the duration metrics onto the specific response headers
			res.Header.Set(ClientTimeHeaderRTTUs, strconv.FormatInt(time.Since(t0).Microseconds(), 10))
			res.Header.Set(ClientTimeHeaderPostLastUs, strconv.FormatInt(time.Since(t1).Microseconds(), 10))
			
			// Lock map access to avoid write collisions
			mu.Lock()
			out[req] = res
			mu.Unlock()
		}(request, stream, t0, t1)
	}

	// 5) Wait for all concurrent stream reads to wrap up
	wg.Wait()
	return out
}
