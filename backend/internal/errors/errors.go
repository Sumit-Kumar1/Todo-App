package errors

import (
	"errors"
	"fmt"
)

const (
	notFoundFormat  = "%s not found"
	invalidFieldFmt = "invalid field: %s"
	missingFieldFmt = "missing field: %s"
)

var (
	ErrUserAlreadyExists = ConstError("user already exists")
	ErrPsswdNotMatch     = ConstError("password does not match")
	ErrUserNotFound      = ConstError("user not found")
	ErrInvalidCookie     = ConstError("invalid cookie")
)

type ConstError string

func NewConstError(message string) ConstError {
	return ConstError(message)
}

func (err ConstError) Error() string {
	return string(err)
}

func (err ConstError) Is(target error) bool {
	var targetErr ConstError
	if errors.As(target, &targetErr) {
		return err.Error() == targetErr.Error()
	}

	return target.Error() == err.Error()
}

func ErrNotFound(entity string) error {
	return NewConstError(fmt.Sprintf(notFoundFormat, entity))
}

func ErrInvalid(entity string) error {
	return NewConstError(fmt.Sprintf(invalidFieldFmt, entity))
}

func ErrRequired(entity string) error {
	return NewConstError(fmt.Sprintf(missingFieldFmt, entity))
}
