package middleware

import (
	"net/http"
	"slices"

	"go-starter/internal/auth"
)

func RequireRole(allowed ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := auth.UserFrom(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if !slices.Contains(allowed, user.RoleName) {
				writeError(w, http.StatusForbidden, "your role is not allowed to access this endpoint")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}