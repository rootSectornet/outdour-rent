package response

import "errors"

// Sentinel errors used by HandleError for HTTP status mapping.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrConflict          = errors.New("resource already exists")
	ErrValidation        = errors.New("validation failed")
	ErrInsufficientStock = errors.New("insufficient stock for requested dates")
)
