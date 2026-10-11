// Package router mendaftarkan semua route ke handler-nya, beserta middleware yang melindunginya.
package router

import (
	"net/http"

	"go-starter/internal/handler"
	"go-starter/internal/middleware"
)

type Deps struct {
	Health *handler.HealthHandler
	Auth   *handler.AuthHandler
	Role   *handler.RoleHandler

	AuthMW          middleware.Middleware
	RoleOwnershipMW middleware.Middleware
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", d.Health.Check)
	mux.HandleFunc("POST /auth/login", d.Auth.Login)

	mux.Handle("GET /roles", middleware.Chain(d.Role.List, d.AuthMW, d.RoleOwnershipMW))
	mux.Handle("GET /roles/{id}", middleware.Chain(d.Role.GetByID, d.AuthMW, d.RoleOwnershipMW))
	return mux
}