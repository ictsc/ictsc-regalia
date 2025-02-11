package domain

import (
	"strings"

	"github.com/cockroachdb/errors"
)

type ErrType int

const (
	ErrTypeUnknown ErrType = iota
	ErrTypeInternal
	ErrTypeInvalidArgument
	ErrTypeNotFound
	ErrTypeAlreadyExists
)

func (t ErrType) String() string {
	switch t {
	case ErrTypeUnknown:
		return "unknown"
	case ErrTypeInternal:
		return "internal"
	case ErrTypeInvalidArgument:
		return "invalid_argument"
	case ErrTypeNotFound:
		return "not_found"
	case ErrTypeAlreadyExists:
		return "already_exists"
	default:
		return "unknown"
	}
}

type Error struct {
	Type ErrType
	Msg  string
	err  error
}

var (
	ErrInternal        = &Error{Type: ErrTypeInternal}
	ErrInvalidArgument = &Error{Type: ErrTypeInvalidArgument}
	ErrNotFound        = &Error{Type: ErrTypeNotFound}
	ErrAlreadyExists   = &Error{Type: ErrTypeAlreadyExists}

	errUnknown = &Error{}
)

func newInternalError(err error) error {
	return &Error{Type: ErrTypeInternal, err: errors.WithStackDepth(err, 1)}
}

func NewInvalidArgumentError(msg string, err error) error {
	return &Error{Type: ErrTypeInvalidArgument, Msg: msg, err: errors.WithStackDepth(err, 1)}
}

func NewNotFoundError(msg string, err error) error {
	return &Error{Type: ErrTypeNotFound, Msg: msg, err: errors.WithStackDepth(err, 1)}
}

func NewAlreadyExistsError(msg string, err error) error {
	return &Error{Type: ErrTypeAlreadyExists, Msg: msg, err: errors.WithStackDepth(err, 1)}
}

func WrapAsInternal(err error, msg string) error {
	if errors.Is(err, errUnknown) {
		return err
	} else {
		return &Error{Type: ErrTypeInternal, err: errors.WrapWithDepth(1, err, msg)}
	}
}

func (e *Error) Error() string {
	var builder strings.Builder
	builder.WriteString(e.Type.String())
	if e.Msg != "" {
		builder.WriteString(": ")
		builder.WriteString(e.Msg)
	}
	if e.err != nil {
		builder.WriteString(": ")
		builder.WriteString(e.err.Error())
	}
	return builder.String()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && (t.Type == ErrTypeUnknown || t.Type == e.Type)
}

func ErrorType(err error) ErrType {
	if e := new(Error); errors.As(err, &e) {
		return e.Type
	}
	return ErrTypeUnknown
}
