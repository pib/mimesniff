package mimesniff

import (
	"io"
	"net/http"
)

// Resource wraps a resource which can be fetched via some means. For
// example, http response, metadata about a file from the file system,
// a file from an FTP server, etc.
type Resource interface {
	Header() []byte
	Body() io.Reader
	SuppliedType() string
	ShouldSniffBinary(SuppliedType string) bool
	NoSniff() bool
}

// ResourceMetadata holds a Resource plus the metadata defined at
// http://mimesniff.spec.whatwg.org/#resource
type ResourceMetadata struct {
	res          Resource
	SuppliedType string          // supplied MIME type
	SniffedType  string          // sniffed MIME type
	ParsedType   *ParsedMimeType // parsed SniffedType
}

// NewResource creates a Resource from the given response. It applies
// the "supplied MIME type detection algorithm" as specified at
// http://mimesniff.spec.whatwg.org/#supplied-mime-type-detection-algorithm,
// as well as the type sniffing algorithm from that same page
func NewResourceMetadata(res Resource) *ResourceMetadata {
	meta := ResourceMetadata{res: res, SuppliedType: res.SuppliedType()}
	meta.sniffType()
	meta.ParsedType = ParseMimeType(meta.SniffedType)

	return &meta
}

func (meta *ResourceMetadata) sniffType() {
	switch meta.SuppliedType {
	case "", "unknown/unknown", "application/unknown", "*/*":
		meta.SniffedType = http.DetectContentType(meta.res.Header())
		return
	}

	if meta.res.NoSniff() {
		meta.SniffedType = meta.SuppliedType
		return
	}

	if meta.res.ShouldSniffBinary(meta.SuppliedType) {
		meta.SniffedType = http.DetectContentType(meta.res.Header())
		return
	}

	parsedType := ParseMimeType(meta.SuppliedType)
	switch {
	case parsedType.IsXml():
		meta.SniffedType = meta.SuppliedType
		return
	case parsedType.MediaType == "text/html":
		meta.SniffedType = distinguishFeed(meta.SuppliedType, string(meta.res.Header()))
		return
	}
}
