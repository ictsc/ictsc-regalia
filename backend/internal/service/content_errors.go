package service

import "errors"

type invalidContentMarker interface {
	error
	InvalidContent() bool
}

func IsInvalidContentError(err error) bool {
	var marker invalidContentMarker
	return errors.As(err, &marker) && marker.InvalidContent()
}
