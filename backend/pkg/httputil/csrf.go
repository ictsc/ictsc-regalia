package httputil

import (
	"net/http"
	"slices"
)

func CSRFMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			if origin := r.Header.Get("Origin"); !slices.Contains(allowedOrigins, origin) {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			if secFetchSite := r.Header.Get("Sec-Fetch-Site"); secFetchSite != "" && secFetchSite != "same-origin" {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
