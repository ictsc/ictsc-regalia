package domain

import "fmt"

// ErrTypeはドメインエラーの種類を表す
type ErrType string

const (
	ErrInvalidArgument ErrType = "invalid_argument"
	ErrNotFound        ErrType = "not_found"
	ErrAlreadyExists   ErrType = "already_exists"
	ErrInternal        ErrType = "internal"
)

// Errorはドメイン層で発生したエラーを表す
type Error struct {
	Type ErrType
	Msg  string
	Err  error
}

// Errorはエラーメッセージを返す
func (e *Error) Error() string {
	if e.Err == nil {
		return e.Msg
	}

	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

// Unwrapは元のエラーを返す
func (e *Error) Unwrap() error {
	return e.Err
}
