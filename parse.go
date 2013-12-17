package mimesniff

import (
	"mime"
	"strings"
)

type ParsedMimeType struct {
	MediaType string // string representation of the MIME type minus any params
	Type      string
	SubType   string
	Params    map[string]string
}

func ParseMimeType(mimeType string) *ParsedMimeType {
	var parsedType ParsedMimeType
	if mediaType, params, err := mime.ParseMediaType(mimeType); err == nil {
		parsedType.MediaType = mediaType
		parts := strings.SplitN(mediaType, "/", 2)
		parsedType.Type = parts[0]
		parsedType.SubType = parts[1]
		parsedType.Params = params
	}
	return &parsedType
}

func (pm *ParsedMimeType) IsXml() bool {
	if pm.SubType == "xml" || strings.HasSuffix(pm.SubType, "+xml") {
		return true
	} else {
		return false
	}
}
