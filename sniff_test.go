package mimesniff

import (
	"io"
	"strings"
	"testing"
)

type testResource struct {
	body           string
	suppliedType   string
	shouldSniffBin bool
	noSniff        bool
	expectedType   string
}

var sniffTests = []testResource{
	// No type specified
	{body: "\x00\x00\x01\x00", expectedType: "image/x-icon"},
	{body: "\x00\x00\x02\x00", expectedType: "image/x-icon"},
	{body: "BM", expectedType: "image/bmp"},
	{body: "GIF89a", expectedType: "image/gif"},
	{body: "GIF87a", expectedType: "image/gif"},
	{body: "RIFF____WEBPVP", expectedType: "image/webp"},
	{body: "\x89PNG\x0D\x0A\x1A\x0A", expectedType: "image/png"},
	{body: "\xFF\xD8\xFF", expectedType: "image/jpeg"},

	// Correct type specified
	{body: "\x00\x00\x01\x00", suppliedType: "image/x-icon", expectedType: "image/x-icon"},
	{body: "\x00\x00\x02\x00", suppliedType: "image/x-icon", expectedType: "image/x-icon"},
	{body: "BM", suppliedType: "image/bmp", expectedType: "image/bmp"},
	{body: "GIF89a", suppliedType: "image/gif", expectedType: "image/gif"},
	{body: "GIF87a", suppliedType: "image/gif", expectedType: "image/gif"},
	{body: "RIFF____WEBPVP", suppliedType: "image/webp", expectedType: "image/webp"},
	{body: "\x89PNG\x0D\x0A\x1A\x0A", suppliedType: "image/png", expectedType: "image/png"},
	{body: "\xFF\xD8\xFF", suppliedType: "image/jpeg", expectedType: "image/jpeg"},

	// Wrong type specified
	{body: "\x00\x00\x01\x00", suppliedType: "image/x-nope", expectedType: "image/x-icon"},
	{body: "\x00\x00\x02\x00", suppliedType: "image/x-nope", expectedType: "image/x-icon"},
	{body: "BM", suppliedType: "image/x-nope", expectedType: "image/bmp"},
	{body: "GIF89a", suppliedType: "image/x-nope", expectedType: "image/gif"},
	{body: "GIF87a", suppliedType: "image/x-nope", expectedType: "image/gif"},
	{body: "RIFF____WEBPVP", suppliedType: "image/x-nope", expectedType: "image/webp"},
	{body: "\x89PNG\x0D\x0A\x1A\x0A", suppliedType: "image/x-nope", expectedType: "image/png"},
	{body: "\xFF\xD8\xFF", suppliedType: "image/x-nope", expectedType: "image/jpeg"},

	// Wrong type specified, noSniff true
	{body: "\x00\x00\x01\x00", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "\x00\x00\x02\x00", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "BM", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "GIF89a", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "GIF87a", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "RIFF____WEBPVP", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "\x89PNG\x0D\x0A\x1A\x0A", suppliedType: "image/x", expectedType: "image/x", noSniff: true},
	{body: "\xFF\xD8\xFF", suppliedType: "image/x", expectedType: "image/x", noSniff: true},

	// Apache bug which marks unknown types as the wrong thing
	{body: "BM", suppliedType: "text/plain", expectedType: "image/bmp", shouldSniffBin: true},
	{body: "BM", suppliedType: "text/plain; charset=ISO-8859-1", expectedType: "image/bmp", shouldSniffBin: true},
	{body: "BM", suppliedType: "text/plain; charset=iso-8859-1", expectedType: "image/bmp", shouldSniffBin: true},
	{body: "BM", suppliedType: "text/plain; charset=UTF-8", expectedType: "image/bmp", shouldSniffBin: true},
}

func TestSniffType(t *testing.T) {
	var meta ResourceMetadata
	for _, res := range sniffTests {
		meta.res, meta.SuppliedType = res, res.SuppliedType()
		meta.sniffType()
		if meta.SniffedType != res.expectedType {
			t.Error("For", res, "expected", res.expectedType, "got", meta.SniffedType)
		}
	}
}

func (res testResource) Header() []byte {
	if len(res.body) > 512 {
		return []byte(res.body[:512])
	}
	return []byte(res.body)
}

func (res testResource) Body() io.Reader {
	return strings.NewReader(string(res.body))
}

func (res testResource) SuppliedType() string {
	return res.suppliedType
}

func (res testResource) ShouldSniffBinary(suppliedType string) bool {
	return res.shouldSniffBin
}

func (res testResource) NoSniff() bool {
	return res.noSniff
}
