package handler

import (
	"errors"
	"log"
	"net/http"

	"go-starter/internal/apperr"
	"go-starter/internal/auth"
	"go-starter/internal/dto"
	"go-starter/internal/service"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	roles, err := h.svc.List(r.Context(), user)
	if err != nil {
		log.Println("list roles:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, dto.RoleListResponse{
		Status:  http.StatusOK,
		Message: "Roles retrieved successfully",
		Data:    roles,
	})
}

func (h *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	role, err := h.svc.GetByID(r.Context(), user, r.PathValue("id"))
	switch {
	case errors.Is(err, apperr.ErrInvalidID):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperr.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, apperr.ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case err != nil:
		log.Println("get role:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	default:
		writeJSON(w, http.StatusOK, dto.RoleDetailResponse{
			Status:  http.StatusOK,
			Message: "Role retrieved successfully",
			Data:    *role,
		})
	}
}