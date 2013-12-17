package mimesniff

import "strings"

// implements "rules for distinguishing if a resource is a feed or
// HTML" from http://mimesniff.spec.whatwg.org/#sniffing-a-mislabeled-feed
func distinguishFeed(suppliedType string, head string) string {
	// Steps 1-3 not needed
	// Step 4: Skip BOM
	head = strings.TrimPrefix(head, "\xef\xbb\xbf")

	// Step 5
	for len(head) > 0 {
		// Step 5.1: find next "<"
		head = strings.TrimSpace(head)
		if !maybeSkip(&head, "<", "") {
			return suppliedType // bail out if there's a non-"<" after the whitespace
		}

		// Step 5.2: Find the MIME type by what tags exist
		if len(head) > 0 { // 5.2.1
			switch {
			case maybeSkip(&head, "!--", "-->"): // 5.2.2: skip comments
			case maybeSkip(&head, "!", ">"): // 5.2.3: Skip declarations
			case maybeSkip(&head, "?", "?>"): // 5.2.4: Skip processing instructions
			case strings.HasPrefix(head, "rss"): // 5.2.5: Check for RSS
				return "application/rss+xml"
			case strings.HasPrefix(head, "feed"): // 5.2.6: Check for Atom
				return "application/atom+xml"
			case maybeSkip(&head, "rdf:RDF", ""): // 5.2.7: Check for RDF/RSS
				rssUrl := "http://purl.org/rss/1.0/"
				rdfUrl := "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				for len(head) > 0 { // 5.2.7.1
					switch { // 5.2.7.2 & 5.2.7.3
					case maybeSkip(&head, rssUrl, rdfUrl), maybeSkip(&head, rdfUrl, rssUrl):
						return "application/rss+xml"
					default: // 5.2.7.4
						head = head[1:]
					}
				}
			default:
				return suppliedType
			}
		}
	}
	return suppliedType
}

func maybeSkip(s *string, prefix string, suffix string) bool {
	if strings.HasPrefix(*s, prefix) {
		*s = (*s)[len(prefix):]
		if end := strings.Index(*s, suffix); end > -1 {
			*s = (*s)[end+len(suffix):]
		} else {
			*s = ""
		}
		return true
	}
	return false
}
