// Package router mendaftarkan semua route ke handler-nya.
package router

import (
	"net/http"

	"go-starter/internal/handler"
)

func New(health *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Check)
	return mux
}
