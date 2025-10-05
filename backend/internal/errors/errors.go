// Package errors unified package for service and handler errors
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	notFoundFormat  = "%s not found"
	invalidFieldFmt = "incorrect value for field: %s"
	missingFieldFmt = "missing field: %s"
)

var (
	ErrUserAlreadyExists = NewConstError("user already exists")
	ErrPsswdNotMatch     = NewConstError("password does not match")
	ErrUserNotFound      = NewConstError("user not found")
	ErrInvalidCookie     = NewConstError("invalid cookie value")
	ErrCookieValTooLong  = NewConstError("cookie value too long")
	ErrTaskNotFound      = NewConstError("task not found")
	ErrInvalidTaskID     = NewConstError("invalid task id")
)

// CustomError represents an error that can be sent in HTTP responses.
// It includes an HTTP status code and an error message.
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface for CustomError.
// It returns the error message.
func (e *CustomError) Error() string {
	return e.Message
}

// NewHTTPError created a custom error based on code, msg and detail of error,
// use only at handler layer
func NewHTTPError(code int, msg, details string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
		Details: details,
	}
}

// ErrBadRequest creates an error for bad request scenarios.
// It wraps the provided error in a CustomError with a 400 status code.
func ErrBadRequest(err error) *CustomError {
	return NewHTTPError(400, err.Error(), "")
}

// ErrNotFound creates an error for when an entity is not found.
// It formats the error message using the notFoundFormat constant.
func ErrNotFound(entity string) *CustomError {
	return NewHTTPError(404, fmt.Sprintf(notFoundFormat, entity), "")
}

// constError is a type that implements the error interface.
// It's used for creating constant error values for internal error use.
type constError string

// NewConstError creates a new constant error with the given message.
// It returns a constError that can be used as a constant error value.
func NewConstError(message string) constError {
	return constError(message)
}

// Error implements the error interface for constError.
// It returns the string representation of the error.
func (err constError) Error() string {
	return string(err)
}

// Is implements error comparison for constError.
// It allows checking if an error matches a specific constError value.
func (err constError) Is(target error) bool {
	var t constError

	ok := errors.As(target, &t)
	if !ok {
		return false
	}

	return err == t
}

// ErrInvalid creates an error for invalid entity scenarios.
// It formats the error message using the invalidFormat constant.
func ErrInvalid(entity string) error {
	return NewConstError(fmt.Sprintf(invalidFieldFmt, entity))
}

// ErrRequired creates an error for required field scenarios.
// It formats the error message using the requiredFormat constant.
func ErrRequired(entity string) error {
	return NewConstError(fmt.Sprintf(missingFieldFmt, entity))
}

func HandleHTTPError(w http.ResponseWriter, err error) {
	hErr := NewHTTPError(http.StatusServiceUnavailable, err.Error(), "")

	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		hErr.Code = http.StatusConflict
	case errors.Is(err, ErrInvalidCookie):
		hErr.Code = http.StatusUnauthorized
	case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrTaskNotFound):
		hErr.Code = http.StatusNotFound
	case errors.Is(err, ErrPsswdNotMatch):
		hErr.Code = http.StatusUnauthorized
	default:
		hErr.Code = http.StatusServiceUnavailable
	}

	respErr := struct {
		Error CustomError `json:"error"`
	}{Error: *hErr}

	data, _ := json.Marshal(respErr)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(hErr.Code)
	_, _ = w.Write(data)
}
