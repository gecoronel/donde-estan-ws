package utils

import (
	"errors"
	"net/http"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
)

func GetHTTPCodeByError(err error) int {
	if errors.Is(err, web.ErrIncorrectPassword) {
		return http.StatusBadRequest
	}
	if errors.Is(err, web.ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, web.ErrConflict) {
		return http.StatusConflict
	}
	if errors.Is(err, web.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, web.ErrInternalServerError) {
		return http.StatusInternalServerError
	}

	// If not exist a status! This doesn't have to happen.
	return http.StatusInternalServerError
}
