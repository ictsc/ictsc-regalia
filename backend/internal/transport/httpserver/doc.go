// Package httpserver serves the REST API under /api/v1 defined by
// backend/openapi.json. It authenticates each request (sessions, GitHub
// Actions OIDC, SState callbacks), checks Origin on browser writes, maps
// requests to package service and turns errors into problem responses.
package httpserver
