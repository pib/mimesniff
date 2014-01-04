package mimesniff

import (
	"code.google.com/p/go.net/html/charset"
	"code.google.com/p/go.text/encoding"
	"fmt"
	"github.com/marketvibe/chardet"
	"os"
	"regexp"
	"strings"
)

var attrRe *regexp.Regexp

func init() {
	attrRe = regexp.MustCompile(`(\w+)="([^"]*)"`)
}

func DetermineEncoding(givenType *ParsedMimeType, head []byte) encoding.Encoding {
	// Try to determine encoding ala whatwg "encoding sniffing algorithm"
	enc, _, certain := charset.DetermineEncoding(head, givenType.String())
	if certain {
		return enc
	}

	// Before doing a byte-level encoding detection, check for xml
	// declaration
	s := strings.TrimSpace(string(head))
	if strings.HasPrefix(s, "<?xml ") {
		s = strings.TrimPrefix(s, "<?xml ")
		end := strings.Index(s, "?>")
		if end > 0 {
			s = s[:end]
			for _, match := range attrRe.FindAllStringSubmatch(s, -1) {
				if match[1] == "encoding" {
					if detectedEnc, _ := charset.Lookup(match[2]); detectedEnc != nil {
						return detectedEnc
					}
				}
			}
		}
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
	if detectedEnc, _ := charset.Lookup(detected.Charset); detectedEnc != nil {
		return detectedEnc
	}

	return enc
}
