package errors

// Flag エラーフラグ
type Flag string

type appError struct {
	flag Flag
	err  error
}

func newAppError(flag Flag, err error) *appError {
	return &appError{flag: flag, err: err}
}

func (e *appError) Error() string {
	return e.err.Error()
}

func (e *appError) Unwrap() error {
	return e.err
}
