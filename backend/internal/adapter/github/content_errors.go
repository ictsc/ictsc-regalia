package github

import (
	"errors"
	"fmt"
)

type invalidContentError struct{ cause error }

func (e invalidContentError) Error() string        { return e.cause.Error() }
func (e invalidContentError) Unwrap() error        { return e.cause }
func (e invalidContentError) InvalidContent() bool { return true }

func wrapContentLoadError(scope string, err error) error {
	var invalid interface{ InvalidContent() bool }
	if errors.As(err, &invalid) && invalid.InvalidContent() {
		return fmt.Errorf("fetch %s: %w", scope, err)
	}
	if errors.Is(err, ErrInvalidManifest) {
		return err
	}
	return fmt.Errorf("fetch %s: %w", scope, err)
}
