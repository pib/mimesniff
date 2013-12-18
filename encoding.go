package mimesniff

import (
	"code.google.com/p/go.text/encoding"
	"code.google.com/p/go.text/encoding/charmap"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/encoding/korean"
	"code.google.com/p/go.text/encoding/simplifiedchinese"
	"code.google.com/p/go.text/encoding/traditionalchinese"
	"code.google.com/p/go.text/encoding/unicode"
	"fmt"
	"github.com/marketvibe/chardet"
	"os"
)

var charsetsEncoding = map[string]encoding.Encoding{
	"Big5":         traditionalchinese.Big5,
	"EUC-JP":       japanese.EUCJP,
	"EUC-KR":       korean.EUCKR,
	"GB-18030":     simplifiedchinese.GB18030,
	"ISO-2022-JP":  japanese.ISO2022JP,
	"ISO-8859-1":   charmap.Windows1252,
	"ISO-8859-5":   charmap.ISO8859_5,
	"ISO-8859-6":   charmap.ISO8859_6,
	"ISO-8859-7":   charmap.ISO8859_7,
	"ISO-8859-8":   charmap.ISO8859_8,
	"KOI8-R":       charmap.KOI8R,
	"Shift_JIS":    japanese.ShiftJIS,
	"UTF-16BE":     unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM),
	"UTF-16LE":     unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM),
	"UTF-8":        encoding.Nop,
	"windows-1251": charmap.Windows1251,
	"windows-1256": charmap.Windows1256,
}

func determineEncoding(givenType *ParsedMimeType, head []byte) string {
	var detector *chardet.Detector
	if givenType.IsHtml() || givenType.IsXml() {
		detector = chardet.NewHtmlDetector()
	} else {
		detector = chardet.NewTextDetector()
	}

	givenCharset := givenType.Params["charset"]
	detectedCharsets, err := detector.DetectAll(head)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error detecting character set, using default.")
		if isValidCharset(givenCharset) {
			return givenCharset
		}
		return ""
	}
	for _, charset := range detectedCharsets {
		if charset.Charset == givenCharset && isValidCharset(givenCharset) {
			return charset.Charset
		}
	}
	if isValidCharset(givenCharset) {
		return givenCharset
	}

	for _, charset := range detectedCharsets {
		if isValidCharset(charset.Charset) {
			return charset.Charset
		}
	}
	return ""
}

func isValidCharset(charset string) bool {
	_, ok := charsetsEncoding[charset]
	return ok
}
