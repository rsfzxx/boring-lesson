package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"go-starter/internal/dto"
	"go-starter/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "body harus berupa JSON yang valid"})
		return
	}
	if req.PhoneNumber == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "phone_number dan password wajib diisi"})
		return
	}

	res, err := h.svc.Login(r.Context(), req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err != nil {
		log.Println("login:", err)
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "terjadi kesalahan di server"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}