package mimesniff

import (
	"io"
	"net/http"
	"strings"
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
		meta.SniffedType = distinguishFeed(meta.res.Header())
		return
	}
}

// implements http://mimesniff.spec.whatwg.org/#rules-for-distinguishing-if-a-resource-is-a-feed-or-html
func distinguishFeed(suppliedType string, head string) string {
	// Steps 1-3 not needed
	// Step 4: Skip BOM
	head = strings.TrimPrefix(head, "\xef\xbb\xbf")

	// Step 5
	for len(head) > 0 {
		// Step 5.1: Skip whitespace and break after finding a "<"
		head = strings.TrimSpace(head)
		if !strings.HasPrefix(head, "<") {
			return suppliedType // bail out if there's a non-"<" after the whitespace
		}

		// Step 5.2: Find first mentioned MIME type
		for s < length { // 5.2.1
			ss := sequence[s:]
			switch {
			case strings.HasPrefix(ss, "!--"): // 5.2.2: Skip comments
				s += 3
				for s < length { // 5.2.2.1
					if strings.HasPrefix(sequence[s:], "-->") { // 5.2.2.2
						s += 3
						break
					}
					s++ // 5.2.2.3
				}
			case strings.HasPrefix(ss, "!"): // 5.2.3: Skip declarations
				s++
				for s < length { // 5.2.3.1
					if strings.HasPrefix(sequence[s:], ">") { // 5.2.3.2
						s++
						break
					}
					s++ // 5.2.3.3
				}
			case strings.HasPrefix(ss, "?"): // 5.2.4: Skip processing instructions
				s++
				for s < length {
				}
			}
		}
	}
	return suppliedType
}
