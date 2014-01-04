package mimesniff

import (
	"code.google.com/p/go.net/html/charset"
	"code.google.com/p/go.text/encoding"
	"fmt"
	"github.com/marketvibe/chardet"
	"os"
)

func DetermineEncoding(givenType *ParsedMimeType, head []byte) encoding.Encoding {
	enc, _, certain := charset.DetermineEncoding(head, givenType.String())
	if certain {
		return enc
	}

	var detector *chardet.Detector
	if givenType.IsHtml() || givenType.IsXml() {
		detector = chardet.NewHtmlDetector()
	} else {
		detector = chardet.NewTextDetector()
	}

	detected, err := detector.DetectBest(head)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error detecting character set, using default.")
		return enc
	}

	if detected.Charset == "GB-18030" {
		detected.Charset = "GB18030"
	}
	if detectedEnc, _ := charset.Lookup(detected.Charset); enc != nil {
		return detectedEnc
	}

	return enc
}
