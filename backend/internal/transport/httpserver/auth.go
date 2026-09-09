package httpserver

import (
	"context"
	"errors"
	"net/http"
	"sort"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) StartDiscordAuthentication(ctx context.Context, request api.StartDiscordAuthenticationRequestObject) (api.StartDiscordAuthenticationResponseObject, error) {
	next := "/"
	if request.Params.Next != nil {
		next = *request.Params.Next
	}
	started, err := h.service.BeginDiscord(ctx, false, next)
	if err != nil {
		return nil, err
	}
	location, cookie := started.URL, cookieValue("oauth2-session", started.SessionToken, service.OAuthTTL, h.options.SecureCookies, http.SameSiteLaxMode)
	return api.StartDiscordAuthentication302Response{Headers: api.OAuthRedirectResponseHeaders{Location: &location, SetCookie: &cookie}}, nil
}

func (h *Handler) CompleteDiscordAuthentication(ctx context.Context, request api.CompleteDiscordAuthenticationRequestObject) (api.CompleteDiscordAuthenticationResponseObject, error) {
	if request.Params.Error != nil {
		return nil, core.NewError(http.StatusUnauthorized, "oauth_state_invalid", "Discord denied authentication")
	}
	oauthToken := cookieToken(requestFrom(ctx), "oauth2-session")
	completed, err := h.service.CompleteDiscord(ctx, false, oauthToken, stringValue(request.Params.Code), stringValue(request.Params.State))
	if err != nil {
		return nil, err
	}
	cookieName, ttl := "signup-session", service.SignupTTL
	if completed.SessionKind == session.KindContestant {
		cookieName, ttl = "user-session", service.ContestantTTL
	}
	location, cookie := completed.Next, cookieValue(cookieName, completed.SessionToken, ttl, h.options.SecureCookies, http.SameSiteStrictMode)
	addCookies(ctx, clearCookie("oauth2-session", h.options.SecureCookies, http.SameSiteLaxMode), cookie)
	return api.CompleteDiscordAuthentication302Response{Headers: api.OAuthRedirectResponseHeaders{Location: &location}}, nil
}

func (h *Handler) StartAdminDiscordAuthentication(ctx context.Context, request api.StartAdminDiscordAuthenticationRequestObject) (api.StartAdminDiscordAuthenticationResponseObject, error) {
	next := "/admin/"
	if request.Params.Next != nil {
		next = *request.Params.Next
	}
	started, err := h.service.BeginDiscord(ctx, true, next)
	if err != nil {
		return nil, err
	}
	location, cookie := started.URL, cookieValue("admin-oauth2-session", started.SessionToken, service.OAuthTTL, h.options.SecureCookies, http.SameSiteLaxMode)
	return api.StartAdminDiscordAuthentication302Response{Headers: api.OAuthRedirectResponseHeaders{Location: &location, SetCookie: &cookie}}, nil
}

func (h *Handler) CompleteAdminDiscordAuthentication(ctx context.Context, request api.CompleteAdminDiscordAuthenticationRequestObject) (api.CompleteAdminDiscordAuthenticationResponseObject, error) {
	if request.Params.Error != nil {
		return nil, core.NewError(http.StatusUnauthorized, "oauth_state_invalid", "Discord denied authentication")
	}
	oauthToken := cookieToken(requestFrom(ctx), "admin-oauth2-session")
	completed, err := h.service.CompleteDiscord(ctx, true, oauthToken, stringValue(request.Params.Code), stringValue(request.Params.State))
	if err != nil {
		return nil, err
	}
	location, cookie := completed.Next, cookieValue("admin-session", completed.SessionToken, service.AdminTTL, h.options.SecureCookies, http.SameSiteStrictMode)
	addCookies(ctx, clearCookie("admin-oauth2-session", h.options.SecureCookies, http.SameSiteLaxMode), cookie)
	return api.CompleteAdminDiscordAuthentication302Response{Headers: api.OAuthRedirectResponseHeaders{Location: &location}}, nil
}

func (h *Handler) SignUpContestant(ctx context.Context, request api.SignUpContestantRequestObject) (api.SignUpContestantResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	signupToken := cookieToken(requestFrom(ctx), "signup-session")
	_, token, err := h.service.SignUp(ctx, signupToken, request.Body.Name, request.Body.DisplayName, request.Body.InvitationCode)
	if err != nil {
		return nil, err
	}
	cookie := cookieValue("user-session", token, service.ContestantTTL, h.options.SecureCookies, http.SameSiteStrictMode)
	addCookies(ctx, clearCookie("signup-session", h.options.SecureCookies, http.SameSiteStrictMode), cookie)
	return api.SignUpContestant204Response{}, nil
}

func (h *Handler) SignOutContestant(ctx context.Context, _ api.SignOutContestantRequestObject) (api.SignOutContestantResponseObject, error) {
	r := requestFrom(ctx)
	var deleteErr error
	for _, name := range []string{"user-session", "signup-session", "oauth2-session"} {
		if token := cookieToken(r, name); token != "" {
			if err := h.service.Sessions.Delete(ctx, token); err != nil && deleteErr == nil {
				deleteErr = err
			}
		}
	}
	addCookies(ctx, clearCookie("user-session", h.options.SecureCookies, http.SameSiteStrictMode),
		clearCookie("signup-session", h.options.SecureCookies, http.SameSiteStrictMode),
		clearCookie("oauth2-session", h.options.SecureCookies, http.SameSiteLaxMode))
	if deleteErr != nil {
		return nil, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not invalidate contestant sessions", deleteErr)
	}
	return api.SignOutContestant204Response{}, nil
}

func (h *Handler) SignOutAdmin(ctx context.Context, _ api.SignOutAdminRequestObject) (api.SignOutAdminResponseObject, error) {
	r := requestFrom(ctx)
	var deleteErr error
	for _, name := range []string{"admin-session", "admin-oauth2-session"} {
		if token := cookieToken(r, name); token != "" {
			if err := h.service.Sessions.Delete(ctx, token); err != nil && deleteErr == nil {
				deleteErr = err
			}
		}
	}
	addCookies(ctx, clearCookie("admin-session", h.options.SecureCookies, http.SameSiteStrictMode),
		clearCookie("admin-oauth2-session", h.options.SecureCookies, http.SameSiteLaxMode))
	if deleteErr != nil {
		return nil, core.WrapError(http.StatusInternalServerError, "internal_error", "Could not invalidate admin sessions", deleteErr)
	}
	return api.SignOutAdmin204Response{}, nil
}

func (h *Handler) CreateContestantImpersonation(ctx context.Context, request api.CreateContestantImpersonationRequestObject) (api.CreateContestantImpersonationResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	token, err := h.service.Impersonate(ctx, principal(ctx).Session, request.Body.ContestantName)
	if err != nil {
		return nil, err
	}
	cookie := cookieValue("user-session", token, service.AdminTTL, h.options.SecureCookies, http.SameSiteStrictMode)
	return api.CreateContestantImpersonation204Response{Headers: api.CreateContestantImpersonation204ResponseHeaders{SetCookie: &cookie}}, nil
}

func (h *Handler) GetViewer(ctx context.Context, _ api.GetViewerRequestObject) (api.GetViewerResponseObject, error) {
	viewer := api.Viewer{}
	r := requestFrom(ctx)
	if token := cookieToken(r, "user-session"); token != "" {
		if data, err := h.service.Sessions.Get(ctx, token, session.KindContestant); err == nil {
			contestant, getErr := h.service.Store.GetContestant(ctx, data.ContestantName)
			if getErr == nil {
				team, teamErr := h.service.Store.GetTeam(ctx, contestant.TeamCode)
				if teamErr != nil {
					return nil, teamErr
				}
				var impersonatedBy *string
				if data.ImpersonatedBy != "" {
					impersonatedBy = &data.ImpersonatedBy
				}
				_ = viewer.FromContestantViewer(api.ContestantViewer{
					State: api.CONTESTANT, Profile: toAPIProfile(contestant), Team: toAPITeam(team), ImpersonatedBy: impersonatedBy,
				})
				return api.GetViewer200JSONResponse{Viewer: viewer}, nil
			}
			if getErr != nil {
				return nil, getErr
			}
		} else if !errors.Is(err, session.ErrNotFound) {
			return nil, err
		}
	}
	if token := cookieToken(r, "signup-session"); token != "" {
		if data, err := h.service.Sessions.Get(ctx, token, session.KindSignup); err == nil {
			_ = viewer.FromDiscordAuthenticatedViewer(api.DiscordAuthenticatedViewer{
				State: api.DISCORDAUTHENTICATED, Discord: api.DiscordIdentity{
					Id: data.Discord.ID, Username: data.Discord.Username, DisplayName: data.Discord.DisplayName,
				},
			})
			return api.GetViewer200JSONResponse{Viewer: viewer}, nil
		} else if !errors.Is(err, session.ErrNotFound) {
			return nil, err
		}
	}
	_ = viewer.FromAnonymousViewer(api.AnonymousViewer{State: api.AnonymousViewerStateANONYMOUS})
	return api.GetViewer200JSONResponse{Viewer: viewer}, nil
}

func (h *Handler) GetAdminViewer(ctx context.Context, _ api.GetAdminViewerRequestObject) (api.GetAdminViewerResponseObject, error) {
	viewer := api.AdminViewer{}
	r := requestFrom(ctx)
	if token := cookieToken(r, "admin-session"); token != "" {
		if data, err := h.service.Sessions.Get(ctx, token, session.KindAdmin); err == nil {
			roles := append([]string(nil), data.RoleIDs...)
			sort.Strings(roles)
			_ = viewer.FromAuthenticatedAdminViewer(api.AuthenticatedAdminViewer{
				State: api.ADMIN,
				Admin: api.AdminIdentity{GuildId: data.GuildID, RoleIds: roles, Discord: api.DiscordIdentity{
					Id: data.Discord.ID, Username: data.Discord.Username, DisplayName: data.Discord.DisplayName,
				}},
			})
			return api.GetAdminViewer200JSONResponse{Viewer: viewer}, nil
		} else if !errors.Is(err, session.ErrNotFound) {
			return nil, err
		}
	}
	_ = viewer.FromAnonymousAdminViewer(api.AnonymousAdminViewer{State: api.AnonymousAdminViewerStateANONYMOUS})
	return api.GetAdminViewer200JSONResponse{Viewer: viewer}, nil
}

func cookieToken(r *http.Request, name string) string {
	if r == nil {
		return ""
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func addCookies(ctx context.Context, cookies ...string) {
	jar, _ := ctx.Value(cookieJarKey{}).(*[]string)
	if jar == nil {
		return
	}
	*jar = append(*jar, cookies...)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
