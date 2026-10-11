package router

import (
	"net/http"

	"go-starter/internal/handler"
)

func New(health *handler.HealthHandler, auth *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Check)
	mux.HandleFunc("POST /auth/login", auth.Login)
	return mux
}