package httpserver

import (
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

type Handler struct {
	service *service.Service
	options Options
}

var _ api.StrictServerInterface = (*Handler)(nil)
