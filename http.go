package mimesniff

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HttpResponse struct {
	res        *http.Response
	bodyReader *bufio.Reader
}

func NewHttpResponse(res *http.Response) *HttpResponse {
	return &HttpResponse{res: res, bodyReader: bufio.NewReader(res.Body)}
}

func (res *HttpResponse) Header() []byte {
	head, err := res.bodyReader.Peek(512)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading 512 bytes of HttpResponse", err)
		fmt.Fprintf(os.Stderr, "Continuing with the %d bytes read\n", len(head))
	}
	return head
}

func (res *HttpResponse) Body() io.Reader {
	return res.bodyReader
}

func (res *HttpResponse) SuppliedType() string {
	if contentType := res.res.Header["Content-Type"]; len(contentType) > 0 {
		return contentType[len(contentType)-1]
	}
	return ""
}

// ShouldSniffBinary implements the "check-for-apache-bug" flag check
// from
// http://mimesniff.spec.whatwg.org/#supplied-mime-type-detection-algorithm
func (res *HttpResponse) ShouldSniffBinary(suppliedType string) bool {
	switch suppliedType {
	case "text/plain", "text/plain; charset=ISO-8859-1",
		"text/plain; charset=iso-8859-1", "text/plain; charset=UTF-8":
		return true
	}
	return false
}

func (res *HttpResponse) NoSniff() bool {
	for _, opt := range res.res.Header["X-Content-Type-Options"] {
		if opt == "nosniff" {
			return true
		}
	}
	return false
}
