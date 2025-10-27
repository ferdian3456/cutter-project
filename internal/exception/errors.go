package exception

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message string, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

// Predefined error types
var (
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "Bad request")
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewAppError(http.StatusForbidden, "Forbidden")
	ErrNotFound            = NewAppError(http.StatusNotFound, "Resource not found")
	ErrMethodNotAllowed    = NewAppError(http.StatusMethodNotAllowed, "Method not allowed")
	ErrConflict            = NewAppError(http.StatusConflict, "Conflict")
	ErrUnprocessableEntity = NewAppError(http.StatusUnprocessableEntity, "Unprocessable entity")
	ErrTooManyRequests     = NewAppError(http.StatusTooManyRequests, "Too many requests")
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrNotImplemented      = NewAppError(http.StatusNotImplemented, "Not implemented")
	ErrServiceUnavailable  = NewAppError(http.StatusServiceUnavailable, "Service unavailable")

	// Authentication errors
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid credentials")
	ErrTokenExpired       = NewAppError(http.StatusUnauthorized, "Token expired")
	ErrInvalidToken       = NewAppError(http.StatusUnauthorized, "Invalid token")
	ErrMissingToken       = NewAppError(http.StatusUnauthorized, "Missing authentication token")

	// User errors
	ErrUserNotFound       = NewAppError(http.StatusNotFound, "User not found")
	ErrUserAlreadyExists  = NewAppError(http.StatusConflict, "User already exists")
	ErrInvalidEmail       = NewAppError(http.StatusBadRequest, "Invalid email address")
	ErrInvalidPassword    = NewAppError(http.StatusBadRequest, "Invalid password")
	ErrPasswordMismatch   = NewAppError(http.StatusBadRequest, "Password mismatch")
	ErrUnauthorizedAccess = NewAppError(http.StatusForbidden, "Unauthorized access")

	// Database errors
	ErrDatabaseConnection = NewAppError(http.StatusInternalServerError, "Database connection error")
	ErrDatabaseQuery      = NewAppError(http.StatusInternalServerError, "Database query error")
	ErrDatabaseTransaction = NewAppError(http.StatusInternalServerError, "Database transaction error")

	// Validation errors
	ErrValidationFailed = NewAppError(http.StatusBadRequest, "Validation failed")
	ErrRequiredField    = NewAppError(http.StatusBadRequest, "Required field is missing")
	ErrInvalidFormat    = NewAppError(http.StatusBadRequest, "Invalid format")
	ErrInvalidLength    = NewAppError(http.StatusBadRequest, "Invalid length")
	ErrInvalidValue     = NewAppError(http.StatusBadRequest, "Invalid value")
)

// WithDetails adds details to an error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Wrap wraps an error with additional context
func (e *AppError) Wrap(err error) *AppError {
	if err != nil {
		e.Details = fmt.Sprintf("%s: %v", e.Message, err)
	}
	return e
}

// Is checks if the error is of the given type
func (e *AppError) Is(target error) bool {
	if other, ok := target.(*AppError); ok {
		return e.Code == other.Code
	}
	return false
}

// HTTPStatus returns the HTTP status code for the error
func (e *AppError) HTTPStatus() int {
	return e.Code
}

// ToFiberError converts an AppError to a Fiber error
func (e *AppError) ToFiberError() *fiber.Error {
	return fiber.NewError(e.Code, e.Message)
}

// FromError creates an AppError from a standard error
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, return it
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Otherwise, create a generic internal server error
	return ErrInternalServerError.Wrap(err)
}