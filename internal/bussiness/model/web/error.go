package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrIncorrectUserOrPassword = errors.New("user or password incorrect")
	ErrBadRequest              = errors.New("bad request")
	ErrConflict                = errors.New("conflict")
	ErrNotFound                = errors.New("not found")
	ErrInternalServerError     = errors.New("impossible to solve")
	ErrNoRows                  = errors.New("sql: no rows in result set")
)

// Error is our custom commercial agreements implementation.
type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status  int    `json:"status"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Status:  e.Status,
		Code:    e.Code,
		Message: e.Message,
	})
}

func (e *Error) StatusCode() int {
	return e.Status
}

// Error returns a string message of the error. It is a concatenation of Code and Message fields.
// This means the Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a new error with the given status code and message.
func NewError(statusCode int, message string) error {
	return NewErrorf(statusCode, message)
}

// NewErrorf creates a new error with the given status code and the message
// formatted according to args and format.
func NewErrorf(status int, format string, args ...interface{}) error {
	return &Error{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args...),
		Status:  status,
	}
}

// Headerer is checked by DefaultErrorEncoder. If an error value implements
// Headerer, the provided headers will be applied to the response writer, after
// the Content-Type is set.
type Headerer interface {
	Headers() http.Header
}

// DefaultNotFoundHandler handler for routing paths that could not be found.
var DefaultNotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	err := NewErrorf(http.StatusNotFound, "resource %s not found", r.URL.Path)
	_ = EncodeJSON(w, err, http.StatusNotFound)
}
