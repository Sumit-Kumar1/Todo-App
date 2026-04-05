// Package errors provides domain error types and sentinel errors for the application.
package errors

import (
	"errors"
	"fmt"
)

// ConstError is a string-based error type that supports const declarations,
// making sentinel errors truly immutable.
type ConstError string

func (e ConstError) Error() string { return string(e) }

// Sentinel errors for domain-specific failures.
const (
	ErrUserAlreadyExists ConstError = "user already exists"
	ErrPsswdNotMatch     ConstError = "password does not match"
	ErrUserNotFound      ConstError = "user not found"
	ErrInvalidCookie     ConstError = "invalid cookie value"
	ErrTaskNotFound      ConstError = "task not found"
	ErrInvalidTaskID     ConstError = "invalid task id"
	ErrInsertFailed      ConstError = "not able to insert into database"
	ErrEmptyDBPassword   ConstError = "empty db password"
	ErrNilDB             ConstError = "db is nil"
	ErrDueDatePast       ConstError = "older due date from today"
	ErrInvalidMigMethod  ConstError = "invalid migration method"
)

// ErrRequired and ErrInvalid are sentinel errors for validation categories.
// Use errors.Is(err, ErrRequired) to check if an error is a missing-field error.
var (
	ErrRequired = errors.New("missing required field")
	ErrInvalid  = errors.New("invalid field value")
)

// ValidationError represents a field-level validation failure.
// It wraps ErrRequired or ErrInvalid so callers can match the category
// with errors.Is and inspect the field with errors.As.
type ValidationError struct {
	Field string
	kind  error
}

func (e ValidationError) Error() string {
	if errors.Is(e.kind, ErrRequired) {
		return "missing field: " + e.Field
	}

	return "invalid field: " + e.Field
}

func (e ValidationError) Unwrap() error { return e.kind }

// Required creates a validation error indicating a missing required field.
func Required(field string) ValidationError {
	return ValidationError{Field: field, kind: ErrRequired}
}

// Invalid creates a validation error indicating an invalid field value.
func Invalid(field string) ValidationError {
	return ValidationError{Field: field, kind: ErrInvalid}
}

// HTTPError represents an error that originated from an HTTP response.
// Used by the client package to propagate upstream API errors.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *HTTPError) Error() string { return e.Message }

// NewHTTPError creates an HTTPError with the given status code and message.
func NewHTTPError(code int, msg, details string) *HTTPError {
	return &HTTPError{Code: code, Message: msg, Details: details}
}

// Wrap adds context to an error using fmt.Errorf wrapping.
// The original error remains matchable via errors.Is and errors.As.
func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
