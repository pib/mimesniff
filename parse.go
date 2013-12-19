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

func (pm *ParsedMimeType) String() string {
	return mime.FormatMediaType(pm.MediaType, pm.Params)
}

func (pm *ParsedMimeType) IsXml() bool {
	if pm.SubType == "xml" || strings.HasSuffix(pm.SubType, "+xml") {
		return true
	} else {
		return false
	}
}

func (pm *ParsedMimeType) IsImage() bool {
	return pm.Type == "image"
}

func (pm *ParsedMimeType) IsAudioVideo() bool {
	switch {
	case pm.Type == "audio", pm.Type == "video", pm.MediaType == "application/ogg":
		return true
	}
	return false
}

func (pm *ParsedMimeType) IsText() bool {
	return pm.Type == "text"
}

func (pm *ParsedMimeType) IsHtml() bool {
	return pm.MediaType == "text/html" || pm.MediaType == "application/xhtml+xml"
}
