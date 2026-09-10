package httpserver

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type Options struct {
	SecureCookies       bool
	AllowedOrigins      []string
	SStateCallbackToken string
	ReadTimeout         time.Duration
}
type cookieJarKey struct{}

type principalKey struct{}
type requestKey struct{}

type Principal struct {
	Session       session.Data
	ActionsClaims *service.ActionsClaims
	Machine       string
}

type Server struct {
	service *service.Service
	options Options
	handler *Handler
}

func New(svc *service.Service, options Options) (http.Handler, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("load embedded OpenAPI: %w", err)
	}
	spec.Servers = nil
	server := &Server{service: svc, options: options}
	server.handler = &Handler{service: svc, options: options}

	strict := api.NewStrictHandlerWithOptions(server.handler, []api.StrictMiddlewareFunc{server.attachRequest}, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeProblem(w, r, core.NewError(http.StatusUnprocessableEntity, "validation_error", err.Error()))
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeProblem(w, r, err)
		},
	})
	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID, chimiddleware.RealIP, chimiddleware.Recoverer)
	router.Use(server.verifyOrigin)
	api.HandlerFromMux(strict, router)
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: server.authenticate,
			MultiError:         true,
		},
		DoNotValidateServers: true,
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter, r *http.Request, opts nethttpmiddleware.ErrorHandlerOpts) {
			status := http.StatusUnprocessableEntity
			code := "validation_error"
			if opts.MatchedRoute == nil {
				status, code = http.StatusNotFound, "resource_not_found"
			}
			if domainErr := preferredDomainError(err); domainErr != nil {
				status, code = domainErr.Status, domainErr.Code
			}
			writeProblem(w, r, &core.Error{Status: status, Code: code, Message: err.Error(), Cause: err})
		},
	})
	responseValidated, err := validateOpenAPIResponses(spec, router)
	if err != nil {
		return nil, err
	}
	return validator(responseValidated), nil
}

func (s *Server) attachRequest(next api.StrictHandlerFunc, _ string) api.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		ctx = context.WithValue(ctx, requestKey{}, r)
		jar := make([]string, 0, 3)
		ctx = context.WithValue(ctx, cookieJarKey{}, &jar)
		response, err := next(ctx, w, r, request)
		for _, cookie := range jar {
			w.Header().Add("Set-Cookie", cookie)
		}
		return response, err
	}
}

func (s *Server) authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	r := input.RequestValidationInput.Request
	principal := Principal{}
	switch input.SecuritySchemeName {
	case "ContestantSession":
		data, err := s.cookieSession(ctx, r, "user-session", session.KindContestant)
		if err != nil {
			return input.NewError(sessionAuthenticationError(err, "invalid_session", "Contestant session is invalid or expired"))
		}
		principal.Session = data
	case "SignupSession":
		data, err := s.cookieSession(ctx, r, "signup-session", session.KindSignup)
		if err != nil {
			return input.NewError(sessionAuthenticationError(err, "invalid_session", "Signup session is invalid or expired"))
		}
		principal.Session = data
	case "OAuthSession":
		data, err := s.cookieSession(ctx, r, "oauth2-session", session.KindOAuth)
		if err != nil {
			return input.NewError(sessionAuthenticationError(err, "oauth_state_invalid", "OAuth session is invalid or expired"))
		}
		principal.Session = data
	case "AdminSession":
		data, err := s.cookieSession(ctx, r, "admin-session", session.KindAdmin)
		if err != nil {
			return input.NewError(sessionAuthenticationError(err, "invalid_session", "Admin session is invalid or expired"))
		}
		principal.Session = data
	case "AdminOAuthSession":
		data, err := s.cookieSession(ctx, r, "admin-oauth2-session", session.KindAdminOAuth)
		if err != nil {
			return input.NewError(sessionAuthenticationError(err, "oauth_state_invalid", "Admin OAuth session is invalid or expired"))
		}
		principal.Session = data
	case "GitHubActionsOIDC":
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok || s.service.Content == nil {
			return input.NewError(core.NewError(http.StatusUnauthorized, "invalid_machine_token", "GitHub Actions token is missing"))
		}
		claims, err := s.service.Content.VerifyActionsToken(ctx, token)
		if err != nil {
			var upstream interface{ MachineTokenUpstream() bool }
			if errors.As(err, &upstream) && upstream.MachineTokenUpstream() {
				return input.NewError(core.WrapError(http.StatusBadGateway, "upstream_unavailable", "GitHub Actions identity provider is unavailable", err))
			}
			return input.NewError(core.WrapError(http.StatusUnauthorized, "invalid_machine_token", "GitHub Actions token is invalid", err))
		}
		principal.ActionsClaims = &claims
		principal.Machine = "github-actions"
	case "SStateCallbackBearer":
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok || subtle.ConstantTimeCompare([]byte(token), []byte(s.options.SStateCallbackToken)) != 1 {
			return input.NewError(core.NewError(http.StatusUnauthorized, "invalid_machine_token", "SState callback token is invalid"))
		}
		principal.Machine = "sstate"
	default:
		return input.NewError(core.NewError(http.StatusUnauthorized, "authentication_required", "Unsupported authentication scheme"))
	}
	*r = *r.WithContext(context.WithValue(r.Context(), principalKey{}, principal))
	return nil
}

func (s *Server) cookieSession(ctx context.Context, r *http.Request, name string, kind session.Kind) (session.Data, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return session.Data{}, err
	}
	return s.service.Sessions.Get(ctx, cookie.Value, kind)
}

func (s *Server) verifyOrigin(next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(s.options.AllowedOrigins))
	for _, origin := range s.options.AllowedOrigins {
		if normalized, err := normalizeOrigin(origin); err == nil {
			allowed[normalized] = struct{}{}
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasBearer := bearerToken(r.Header.Get("Authorization"))
		adminCookie, adminCookieErr := r.Cookie("admin-session")
		hasAdminCookie := adminCookieErr == nil && adminCookie.Value != ""
		machineRequest := strings.HasSuffix(r.URL.Path, "/events") ||
			r.URL.Path == "/api/v1/admin/content/actions/refresh" && hasBearer && !hasAdminCookie
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions ||
			machineRequest {
			next.ServeHTTP(w, r)
			return
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			writeProblem(w, r, core.NewError(http.StatusForbidden, "origin_forbidden", "Origin header is required for browser mutations"))
			return
		}
		normalized, err := normalizeOrigin(origin)
		if err != nil {
			writeProblem(w, r, core.NewError(http.StatusForbidden, "origin_forbidden", "Origin is invalid"))
			return
		}
		if _, ok := allowed[normalized]; !ok {
			writeProblem(w, r, core.NewError(http.StatusForbidden, "origin_forbidden", "Origin is not allowed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func normalizeOrigin(value string) (string, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("invalid origin")
	}
	return strings.ToLower(parsed.Scheme + "://" + parsed.Host), nil
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	returnValue := ""
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		returnValue = parts[1]
	}
	return returnValue, returnValue != ""
}

func principal(ctx context.Context) Principal {
	value, _ := ctx.Value(principalKey{}).(Principal)
	return value
}

func requestFrom(ctx context.Context) *http.Request {
	request, _ := ctx.Value(requestKey{}).(*http.Request)
	return request
}

func writeProblem(w http.ResponseWriter, r *http.Request, err error) {
	status, code, detail := http.StatusInternalServerError, "internal_error", "Internal server error"
	var domainErr *core.Error
	if errors.As(err, &domainErr) {
		if domainErr.Status >= 400 && domainErr.Status <= 599 {
			status = domainErr.Status
		}
		code, detail = normalizeErrorCode(domainErr.Code), domainErr.Message
		if domainErr.RetryAfter > 0 {
			seconds := int((domainErr.RetryAfter + time.Second - 1) / time.Second)
			if seconds < 1 {
				seconds = 1
			}
			w.Header().Set("Retry-After", fmt.Sprint(seconds))
		}
	}
	instance := r.URL.Path
	problem := api.ProblemDetails{
		Type: "https://score.invalid/problems/" + code, Title: http.StatusText(status), Status: status,
		Detail: detail, Code: api.ErrorCode(code), Instance: &instance,
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

func normalizeErrorCode(code string) string {
	allowed := map[string]struct{}{
		"admin_role_required": {}, "answer_rate_limited": {}, "authentication_required": {}, "conflict": {},
		"content_invalid": {}, "content_not_available": {}, "content_not_found": {}, "content_refresh_in_progress": {},
		"content_rollback_rejected": {}, "contestant_already_registered": {}, "deployment_in_progress": {},
		"deployment_not_allowed": {}, "duplicate_event_mismatch": {}, "guild_membership_required": {}, "internal_error": {},
		"invalid_deployment_transition": {}, "invalid_machine_token": {}, "invalid_session": {}, "invitation_already_used": {},
		"invitation_expired": {}, "oauth_state_invalid": {}, "origin_forbidden": {}, "permission_denied": {},
		"resource_not_found": {}, "submission_closed": {}, "team_code_conflict": {}, "team_full": {},
		"upstream_unavailable": {}, "validation_error": {},
		"team_role_required": {}, "ambiguous_team_roles": {}, "team_role_mismatch": {},
	}
	if _, ok := allowed[code]; ok {
		return code
	}
	switch code {
	case "content_unavailable":
		return "content_not_available"
	case "manifest_invalid":
		return "content_invalid"
	case "content_changed":
		return "content_refresh_in_progress"
	case "team_conflict":
		return "team_code_conflict"
	case "contestant_conflict":
		return "contestant_already_registered"
	}
	if code == "team_not_found" || code == "contestant_not_found" || code == "answer_not_found" || code == "deployment_not_found" {
		return "resource_not_found"
	}
	if strings.Contains(code, "conflict") || code == "team_not_empty" || code == "invitation_unavailable" {
		return "conflict"
	}
	return "internal_error"
}

func cookieValue(name, value string, ttl time.Duration, secure bool, sameSite http.SameSite) string {
	return (&http.Cookie{
		Name: name, Value: value, Path: "/", MaxAge: int(ttl / time.Second), Expires: time.Now().Add(ttl),
		HttpOnly: true, Secure: secure, SameSite: sameSite,
	}).String()
}

func clearCookie(name string, secure bool, sameSite http.SameSite) string {
	return (&http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1, Expires: time.Unix(1, 0), HttpOnly: true, Secure: secure, SameSite: sameSite}).String()
}

var _ = requestFrom
