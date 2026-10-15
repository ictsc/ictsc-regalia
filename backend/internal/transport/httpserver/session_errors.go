package httpserver

import (
	"errors"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

func sessionAuthenticationError(err error, invalidCode, invalidMessage string) error {
	if errors.Is(err, session.ErrNotFound) || errors.Is(err, http.ErrNoCookie) {
		return core.NewError(http.StatusUnauthorized, invalidCode, invalidMessage)
	}
	return core.WrapError(http.StatusInternalServerError, "internal_error", "Session storage is unavailable", err)
}
