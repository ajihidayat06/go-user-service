// internal/pkg/errors/errors.go
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Business Logic Errors
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"

	// System Errors
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
	ErrCodeExternal ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	ErrCodeTimeout  ErrorCode = "TIMEOUT_ERROR"
)

// AppError represents application-specific error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Cause      error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// Error constructors
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getHTTPStatusCode(code),
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getHTTPStatusCode(code),
		Cause:      err,
	}
}

func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// Helper functions
func getHTTPStatusCode(code ErrorCode) int {
	switch code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeDatabase, ErrCodeExternal, ErrCodeInternal, ErrCodeTimeout:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsErrorCode checks if error has specific code
func IsErrorCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// As is wrapper for errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is is wrapper for errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}
