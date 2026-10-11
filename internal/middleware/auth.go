package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"

	"go-starter/internal/auth"
	"go-starter/internal/repository"
	"go-starter/internal/service"
)

func Auth(jwtSecret string, roles *repository.RoleRepository) Middleware {
	secret := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || tokenString == "" {
				writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			claims := &service.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims,
				func(*jwt.Token) (any, error) { return secret, nil },
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
				jwt.WithExpirationRequired(),
			)
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			roleID, err := bson.ObjectIDFromHex(claims.RoleID)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}
			role, err := roles.FindByID(r.Context(), roleID)
			if errors.Is(err, repository.ErrNotFound) {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}
			if err != nil {
				log.Println("auth middleware:", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			user := auth.User{
				ID:       claims.Subject,
				RoleID:   claims.RoleID,
				RoleName: role.Name,
				GroupID:  claims.GroupID,
			}
			next.ServeHTTP(w, r.WithContext(auth.WithUser(r.Context(), user)))
		})
	}
}