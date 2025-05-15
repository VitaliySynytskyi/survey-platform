package model

import (
	"fmt"
)

// AppError is a custom error type
type AppError struct {
	Message string
	Type    string
	Code    int
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message string, errType string, code int) *AppError {
	return &AppError{
		Message: message,
		Type:    errType,
		Code:    code,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *AppError {
	return NewAppError(message, "NOT_FOUND", 404)
}

// NewAccessDeniedError creates a new access denied error
func NewAccessDeniedError(message string) *AppError {
	return NewAppError(message, "ACCESS_DENIED", 403)
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError(message, "BAD_REQUEST", 400)
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(err error) *AppError {
	return NewAppError(fmt.Sprintf("Internal server error: %v", err), "INTERNAL_SERVER_ERROR", 500)
}
