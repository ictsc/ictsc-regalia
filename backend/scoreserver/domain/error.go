package domain

import "github.com/cockroachdb/errors"

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
	typ ErrType
	err error
}

func NewError(typ ErrType, err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		typ: typ,
		err: err,
	}
}

func WrapAsInternal(err error, msg string) error {
	return NewError(ErrTypeInternal, errors.WrapWithDepth(1, err, msg))
}

func (e *Error) Error() string {
	return e.typ.String() + ": " + e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func ErrTypeFrom(err error) ErrType {
	if e := new(Error); errors.As(err, &e) {
		return e.typ
	}
	return ErrTypeUnknown
}
