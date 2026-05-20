package apperror

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidState  = errors.New("invalid state")
	ErrUnauthorized  = errors.New("unauthorized")
)

type APIError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func New(code, message string) APIError {
	return APIError{Code: code, Message: message}
}
