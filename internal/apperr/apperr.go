package apperr

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrInvalidID  = errors.New("invalid id")

	ErrUnauthorized = errors.New("authentication required")

	ErrForbidden = errors.New("you do not have permission to access this resource")

	ErrNotFound = errors.New("resource not found")

	ErrConflict = errors.New("resource already exists")

	ErrValidation = errors.New("validation failed")

	ErrTooManyRequests = errors.New("too many requests")

	ErrInternal = errors.New("internal server error")
)

func HTTPStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrInvalidID):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrValidation):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrTooManyRequests):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func Message(err error) string {
	if HTTPStatus(err) == http.StatusInternalServerError {
		return ErrInternal.Error()
	}
	return err.Error()
}