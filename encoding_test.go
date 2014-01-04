package mimesniff

import (
	"code.google.com/p/go.text/encoding"
	"code.google.com/p/go.text/encoding/charmap"
	"testing"
)

type encodingTest struct {
	mimeType         string
	head             string
	expectedEncoding encoding.Encoding
}

var encodingTests = []encodingTest{
	{"text/html", `<html><head><meta http-equiv="content-type" content="text/html; charset=UTF-8"></head></html>`, encoding.Nop},
	{"text/html", `<?xml version="1.0" encoding="UTF-8"?>`, encoding.Nop},
	{"text/xml", `<?xml version="1.0" encoding="iso8859-1"?>`, charmap.Windows1252},
	{"text/xml", `⚂`, encoding.Nop},
}

func TestDetermineEncoding(t *testing.T) {
	for _, test := range encodingTests {
		mimeType := ParseMimeType(test.mimeType)
		if enc := DetermineEncoding(mimeType, []byte(test.head)); enc != test.expectedEncoding {
			t.Errorf("Expected %v got %v\n%v", test.expectedEncoding, enc, test)
		}
	}
}
