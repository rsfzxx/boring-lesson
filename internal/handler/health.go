// Package handler menerjemahkan HTTP request ke pemanggilan service,
// lalu menulis response. Tidak ada logika bisnis di sini.
package handler

import (
	"context"
	"net/http"
	"time"

	"go-starter/internal/service"
)

type HealthHandler struct {
	svc *service.HealthService
}

func NewHealthHandler(svc *service.HealthService) *HealthHandler {
	return &HealthHandler{svc: svc}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	res := h.svc.Check(ctx)
	status := http.StatusOK
	if res.Status != "ok" {
		status = http.StatusServiceUnavailable
	}
	writeJSON(w, status, res)
}
