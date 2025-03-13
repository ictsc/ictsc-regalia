package domain

import (
	"encoding"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

type DeploymentStatus uint32

const (
	DeploymentStatusUnknown DeploymentStatus = iota
	DeploymentStatusQueued
	DeploymentStatusCreating
	DeploymentStatusCompleted
	DeploymentStatusFailed
)

func (d DeploymentStatus) IsFinished() bool {
	return d == DeploymentStatusCompleted || d == DeploymentStatusFailed
}

func (d DeploymentStatus) Next() []DeploymentStatus {
	switch d {
	case DeploymentStatusQueued:
		return []DeploymentStatus{DeploymentStatusCreating}
	case DeploymentStatusCreating:
		return []DeploymentStatus{DeploymentStatusCompleted, DeploymentStatusFailed}
	case DeploymentStatusCompleted:
		fallthrough
	case DeploymentStatusFailed:
		fallthrough
	case DeploymentStatusUnknown:
		fallthrough
	default:
		return []DeploymentStatus{}
	}
}

func (d DeploymentStatus) String() string {
	switch d {
	case DeploymentStatusQueued:
		return "Queued"
	case DeploymentStatusCreating:
		return "Creating"
	case DeploymentStatusCompleted:
		return "Completed"
	case DeploymentStatusFailed:
		return "Failed"
	case DeploymentStatusUnknown:
		fallthrough
	default:
		return "Unknown" //nolint:goconst // それぞれの Unknown は別の Enum の値
	}
}

var (
	_ encoding.TextMarshaler   = DeploymentStatus(0)
	_ encoding.TextUnmarshaler = (*DeploymentStatus)(nil)
)

func (d DeploymentStatus) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *DeploymentStatus) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Queued":
		*d = DeploymentStatusQueued
	case "Creating":
		*d = DeploymentStatusCreating
	case "Completed":
		*d = DeploymentStatusCompleted
	case "Failed":
		*d = DeploymentStatusFailed
	default:
		*d = DeploymentStatusUnknown
	}
	return nil
}

var (
	_ json.Marshaler   = DeploymentStatus(0)
	_ json.Unmarshaler = (*DeploymentStatus)(nil)
)

func (d DeploymentStatus) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(d.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DeploymentPhase")
	}
	return b, nil
}

func (d *DeploymentStatus) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return errors.Wrap(err, "failed to unmarshal DeploymentPhase")
	}
	return d.UnmarshalText([]byte(text))
}
