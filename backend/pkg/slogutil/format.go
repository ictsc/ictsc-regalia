package slogutil

import (
	"strings"

	"github.com/cockroachdb/errors"
)

type Format int

const (
	FormatJSON Format = iota
	FormatConsole
	FormatPretty
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatConsole:
		return "console"
	case FormatPretty:
		return "pretty"
	default:
		return "unknown"
	}
}

func (f Format) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Format) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "json":
		*f = FormatJSON
	case "console":
		*f = FormatConsole
	case "pretty":
		*f = FormatPretty
	default:
		return errors.New("unknown format")
	}
	return nil
}
