package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"go-starter/internal/apperr"
	"go-starter/internal/auth"
)

type OwnershipChecker func(ctx context.Context, user auth.User, resourceID string) error

func Ownership(param string, check OwnershipChecker) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := auth.UserFrom(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "authentication required")
				return
			}

			resourceID := r.PathValue(param)
			if resourceID == "" {
				next.ServeHTTP(w, r)
				return
			}

			err := check(r.Context(), user, resourceID)
			switch {
			case err == nil:
				next.ServeHTTP(w, r)
			case errors.Is(err, apperr.ErrForbidden):
				writeError(w, http.StatusForbidden, err.Error())
			case errors.Is(err, apperr.ErrNotFound):
				writeError(w, http.StatusNotFound, err.Error())
			case errors.Is(err, apperr.ErrInvalidID):
				writeError(w, http.StatusBadRequest, err.Error())
			default:
				log.Println("ownership middleware:", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		})
	}
}