package middleware

import (
	"encoding/json"
	"net/http"

	"go-starter/internal/dto"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.HandlerFunc, mws ...Middleware) http.Handler {
	var wrapped http.Handler = h
	for i := len(mws) - 1; i >= 0; i-- {
		wrapped = mws[i](wrapped)
	}
	return wrapped
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Status: status, Message: message})
}