package contestant

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/cockroachdb/errors"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func withImpersonationAuth(base *url.URL, adminEnforcer *adminauth.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth/") {
				next.ServeHTTP(w, r)
				return
			}

			userSess, err := session.UserSessionStore.Get(r.Context())
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if userSess == nil || userSess.Impersonation == nil {
				next.ServeHTTP(w, r)
				return
			}

			viewer := adminauth.GetViewer(r.Context())
			ok, err := adminEnforcer.Enforce(viewer, "contestants", "impersonate")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if ok {
				next.ServeHTTP(w, r)
				return
			}

			if err := session.UserSessionStore.Write(r, w, nil, userSessionOption(base, r)); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			status := http.StatusUnauthorized
			if slices.Contains(viewer.Groups, "system:authenticated") {
				status = http.StatusForbidden
			}
			http.Error(w, http.StatusText(status), status)
		})
	}
}
