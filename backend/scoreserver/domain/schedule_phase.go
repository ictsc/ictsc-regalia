package domain

import (
	"encoding"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

type Phase int32

const (
	PhaseUnspecified Phase = iota
	PhaseOutOfContest
	PhaseInContest
	PhaseBreak
	PhaseAfterContest
)

func (p Phase) String() string {
	switch p {
	case PhaseOutOfContest:
		return "OUT_OF_CONTEST"
	case PhaseInContest:
		return "IN_CONTEST"
	case PhaseBreak:
		return "BREAK"
	case PhaseAfterContest:
		return "AFTER_CONTEST"
	case PhaseUnspecified:
		fallthrough
	default:
		return "UNSPECIFIED"
	}
}

var (
	_ encoding.TextMarshaler   = Phase(0)
	_ encoding.TextUnmarshaler = (*Phase)(nil)
	_ json.Marshaler           = Phase(0)
	_ json.Unmarshaler         = (*Phase)(nil)
)

func (p Phase) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Phase) UnmarshalText(text []byte) error {
	switch string(text) {
	case "OUT_OF_CONTEST":
		*p = PhaseOutOfContest
	case "IN_CONTEST":
		*p = PhaseInContest
	case "BREAK":
		*p = PhaseBreak
	case "AFTER_CONTEST":
		*p = PhaseAfterContest
	default:
		*p = PhaseUnspecified
	}
	return nil
}

func (p Phase) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(p.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Phase")
	}
	return b, nil
}

func (p *Phase) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.Wrap(err, "failed to unmarshal Phase")
	}
	return p.UnmarshalText([]byte(s))
}
