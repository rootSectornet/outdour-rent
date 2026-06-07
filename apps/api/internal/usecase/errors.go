package usecase

import "errors"

// Common business errors.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrConflict           = errors.New("resource already exists")
	ErrValidation         = errors.New("validation failed")
	ErrInsufficientStock  = errors.New("insufficient stock for requested dates")
	ErrReservationExpired = errors.New("reservation has expired")
	ErrInvalidStatus      = errors.New("invalid status transition")
	ErrPaymentFailed      = errors.New("payment processing failed")
)
