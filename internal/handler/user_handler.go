package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"go-starter/internal/apperr"
	"go-starter/internal/auth"
	"go-starter/internal/dto"
	"go-starter/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body must be valid JSON")
		return
	}

	user, err := h.svc.Create(r.Context(), actor, req)
	if err != nil {
		writeAppError(w, "create user", err)
		return
	}
	writeJSON(w, http.StatusCreated, dto.UserDetailResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    *user,
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	users, err := h.svc.List(r.Context(), actor)
	if err != nil {
		writeAppError(w, "list users", err)
		return
	}
	writeJSON(w, http.StatusOK, dto.UserListResponse{
		Status:  http.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	user, err := h.svc.GetByID(r.Context(), actor, r.PathValue("id"))
	if err != nil {
		writeAppError(w, "get user", err)
		return
	}
	writeJSON(w, http.StatusOK, dto.UserDetailResponse{
		Status:  http.StatusOK,
		Message: "User retrieved successfully",
		Data:    *user,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body must be valid JSON")
		return
	}

	user, err := h.svc.Update(r.Context(), actor, r.PathValue("id"), req)
	if err != nil {
		writeAppError(w, "update user", err)
		return
	}
	writeJSON(w, http.StatusOK, dto.UserDetailResponse{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    *user,
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(w, r)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), actor, r.PathValue("id")); err != nil {
		writeAppError(w, "delete user", err)
		return
	}
	writeJSON(w, http.StatusOK, dto.MessageResponse{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	})
}

func currentUser(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	user, ok := auth.UserFrom(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
	}
	return user, ok
}

func writeAppError(w http.ResponseWriter, action string, err error) {
	status := apperr.HTTPStatus(err)
	if status == http.StatusInternalServerError {
		log.Printf("%s: %v", action, err)
	}
	writeError(w, status, apperr.Message(err))
}