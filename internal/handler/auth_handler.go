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
		writeError(w, http.StatusBadRequest, "body harus berupa JSON yang valid")
		return
	}
	if req.PhoneNumber == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "phone_number dan password wajib diisi")
		return
	}

	res, err := h.svc.Login(r.Context(), req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		log.Println("login:", err)
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan di server")
		return
	}

	res.Status = http.StatusOK
	writeJSON(w, res.Status, res)
}