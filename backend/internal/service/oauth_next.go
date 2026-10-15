package service

import (
	"net/url"
	"strings"
	"unicode"
)

// validOAuthNext accepts only an origin-relative absolute path. It repeatedly
// decodes the path before checking it so encoded and double-encoded backslashes
// cannot be interpreted as a network-path reference by WHATWG URL parsers.
func validOAuthNext(next string) bool {
	if strings.ContainsAny(next, "?#\\") {
		return false
	}
	parsed, err := url.ParseRequestURI(next)
	if err != nil || parsed.IsAbs() || parsed.Host != "" || parsed.Scheme != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	decoded := parsed.Path
	for range 8 {
		nextDecoded, err := url.PathUnescape(decoded)
		if err != nil {
			return false
		}
		if nextDecoded == decoded {
			break
		}
		decoded = nextDecoded
	}
	if decoded == "" || decoded[0] != '/' || strings.HasPrefix(decoded, "//") || strings.Contains(decoded, "\\") || strings.Contains(decoded, "%") {
		return false
	}
	for _, character := range decoded {
		if unicode.IsControl(character) {
			return false
		}
	}
	return true
}
