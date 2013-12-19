package mimesniff

import (
	"code.google.com/p/go.text/transform"
	"io"
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

// http://mimesniff.spec.whatwg.org/#determining-the-sniffed-mime-type-of-a-resource
func (meta *ResourceMetadata) sniffType() {
	switch meta.SuppliedType {
	case "", "unknown/unknown", "application/unknown", "*/*":
		meta.SniffedType = DetectContentType(meta.res.Header())
		return
	}

	if meta.res.NoSniff() {
		meta.SniffedType = meta.SuppliedType
		return
	}

	if meta.res.ShouldSniffBinary(meta.SuppliedType) {
		meta.SniffedType = DetectContentType(meta.res.Header())
		return
	}

	parsedType := ParseMimeType(meta.SuppliedType)
	switch {
	case parsedType.IsXml():
		meta.SniffedType = meta.SuppliedType
	case parsedType.MediaType == "text/html":
		meta.SniffedType = distinguishFeed(meta.SuppliedType, string(meta.res.Header()))
	case parsedType.IsImage(), parsedType.IsAudioVideo():
		if m := DetectContentType(meta.res.Header()); m != "application/octet-stream" {
			meta.SniffedType = m
		}
	default:
		meta.SniffedType = meta.SuppliedType
	}
}

func (meta *ResourceMetadata) DecodeBody() io.Reader {
	encoding := determineEncoding(meta.ParsedType, meta.res.Header())
	return transform.NewReader(meta.res.Body(), encoding.NewDecoder())
}
