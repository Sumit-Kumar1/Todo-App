// Package errors unified package for service and handler errors
package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	invalidFieldFmt = "incorrect value for field: %s"
	missingFieldFmt = "missing field: %s"
)

var (
	ErrUserAlreadyExists = NewConstError("user already exists")
	ErrPsswdNotMatch     = NewConstError("password does not match")
	ErrUserNotFound      = NewConstError("user not found")
	ErrInvalidCookie     = NewConstError("invalid cookie value")
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

// ConstError is a type that implements the error interface.
// It's used for creating constant error values for internal error use.
type ConstError string

// NewConstError creates a new constant error with the given message.
// It returns a ConstError that can be used as a constant error value.
func NewConstError(message string) ConstError {
	return ConstError(message)
}

// Error implements the error interface for ConstError.
// It returns the string representation of the error.
func (err ConstError) Error() string {
	return string(err)
}

// Is implements error comparison for ConstError.
// It allows checking if an error matches a specific ConstError value.
func (err ConstError) Is(target error) bool {
	var t ConstError

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

func HandleHTTPError(c *gin.Context, err error) {
	hErr := NewHTTPError(http.StatusServiceUnavailable, err.Error(), "")

	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		hErr.Code = http.StatusConflict
	case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrTaskNotFound):
		hErr.Code = http.StatusNotFound
	case errors.Is(err, ErrInvalidCookie):
		hErr.Code = http.StatusUnauthorized
	case errors.Is(err, NewHTTPError(http.StatusBadRequest, "", "")):
		hErr.Code = http.StatusBadRequest
	default:
		hErr.Code = http.StatusServiceUnavailable
	}

	c.AbortWithStatusJSON(hErr.Code, hErr)
}
