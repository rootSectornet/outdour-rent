package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rentoutdoor/api/internal/usecase"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse is the standard API error response.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success sends a successful response.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with pagination metadata.
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an error response.
func Error(c *gin.Context, statusCode int, errType string, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   errType,
		Message: message,
	})
}

// ValidationError sends a 422 response with field-level errors.
func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]map[string]string, len(ve))
		for i, fe := range ve {
			details[i] = map[string]string{
				"field":   fe.Field(),
				"tag":     fe.Tag(),
				"message": formatValidationMessage(fe),
			}
		}
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Success: false,
			Error:   "validation_error",
			Message: "request validation failed",
			Details: details,
		})
		return
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   "bad_request",
		Message: err.Error(),
	})
}

// HandleError maps domain errors to HTTP status codes.
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		Error(c, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, usecase.ErrUnauthorized):
		Error(c, http.StatusUnauthorized, "unauthorized", err.Error())
	case errors.Is(err, usecase.ErrForbidden):
		Error(c, http.StatusForbidden, "forbidden", err.Error())
	case errors.Is(err, usecase.ErrConflict):
		Error(c, http.StatusConflict, "conflict", err.Error())
	case errors.Is(err, usecase.ErrValidation):
		Error(c, http.StatusUnprocessableEntity, "validation_error", err.Error())
	case errors.Is(err, usecase.ErrInsufficientStock):
		Error(c, http.StatusConflict, "insufficient_stock", err.Error())
	default:
		Error(c, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

func formatValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email"
	case "min":
		return fe.Field() + " must be at least " + fe.Param()
	case "max":
		return fe.Field() + " must be at most " + fe.Param()
	case "uuid":
		return fe.Field() + " must be a valid UUID"
	case "oneof":
		return fe.Field() + " must be one of: " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}
