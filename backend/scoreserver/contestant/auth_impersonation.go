package contestant

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type (
	ImpersonationEffect interface {
		domain.TeamMemberGetter
		domain.UserLister
	}
	startImpersonationRequest struct {
		Name     string `json:"name"`
		TeamCode int64  `json:"teamCode"`
	}
)

var errAdminPermissionDenied = errors.New("permission denied")

func (h *AuthHandler) handleStartImpersonation(w http.ResponseWriter, r *http.Request) {
	adminViewer, err := h.authorizeAdmin(r.Context(), "contestants", "impersonate")
	if err != nil {
		writeAdminAuthorizationError(w, err)
		return
	}

	var req startImpersonationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userName, err := domain.NewUserName(req.Name)
	if err != nil {
		http.Error(w, "invalid contestant name", http.StatusBadRequest)
		return
	}
	teamCode, err := domain.NewTeamCode(req.TeamCode)
	if err != nil {
		http.Error(w, "invalid team code", http.StatusBadRequest)
		return
	}

	user, err := userName.User(r.Context(), h.ImpersonationEffect)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "contestant not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	member, err := user.ID().TeamMember(r.Context(), h.ImpersonationEffect)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "contestant not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if member.Team().Code() != teamCode {
		http.Error(w, "contestant not found", http.StatusNotFound)
		return
	}

	if err := session.UserSessionStore.Write(r, w, &session.UserSession{
		UserID: uuid.UUID(member.ID()),
		Impersonation: &session.ImpersonationSession{
			AdminName: adminViewer.Name,
		},
	}, userSessionOption(h.ExternalURL, r)); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := session.SignUpSessionStore.Write(r, w, nil, signUpSessionOption(h.ExternalURL, r)); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) authorizeAdmin(ctx context.Context, obj, act string) (*adminauth.Viewer, error) {
	viewer := adminauth.GetViewer(ctx)
	ok, err := h.AdminEnforcer.Enforce(viewer, obj, act)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errAdminPermissionDenied
	}
	return &viewer, nil
}

func writeAdminAuthorizationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, adminauth.ErrUnauthenticated):
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	case errors.Is(err, errAdminPermissionDenied):
		http.Error(w, "Forbidden", http.StatusForbidden)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
