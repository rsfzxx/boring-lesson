package router

import (
	"net/http"

	"go-starter/internal/auth"
	"go-starter/internal/handler"
	"go-starter/internal/middleware"
)

type Deps struct {
	Health *handler.HealthHandler
	Auth   *handler.AuthHandler
	Role   *handler.RoleHandler
	User   *handler.UserHandler

	AuthMW          middleware.Middleware
	RoleOwnershipMW middleware.Middleware
	UserOwnershipMW middleware.Middleware
}

func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	adminOnly := middleware.RequireRole(auth.RoleAdmin)
	adminOrLeader := middleware.RequireRole(auth.RoleAdmin, auth.RoleLeader)

	mux.HandleFunc("GET /health", d.Health.Check)
	mux.HandleFunc("POST /auth/login", d.Auth.Login)

	mux.Handle("GET /roles", middleware.Chain(d.Role.List, d.AuthMW, d.RoleOwnershipMW))
	mux.Handle("GET /roles/{id}", middleware.Chain(d.Role.GetByID, d.AuthMW, d.RoleOwnershipMW))

	mux.Handle("POST /users", middleware.Chain(d.User.Create, d.AuthMW, adminOrLeader))
	mux.Handle("GET /users", middleware.Chain(d.User.List, d.AuthMW, d.UserOwnershipMW))
	mux.Handle("GET /users/{id}", middleware.Chain(d.User.GetByID, d.AuthMW, d.UserOwnershipMW))
	mux.Handle("PUT /users/{id}", middleware.Chain(d.User.Update, d.AuthMW, adminOnly))
	mux.Handle("DELETE /users/{id}", middleware.Chain(d.User.Delete, d.AuthMW, adminOnly))

	return mux
}