package models

import (
	"fmt"
	"strings"
)

const (
	notFoundFormat = "Not Found: %s"
	invalidFormat  = "Invalid attribute: %s"
	missingFormat  = "Missing attribute: %s"
)

var (
	ErrPermissionDenied  = constError("Permission denied")
	ErrUserAlreadyExists = constError("User already exists")
	ErrPsswdNotMatch     = constError("Password does not match")
)

type constError string

func NewConstError(message string) constError {
	return constError(message)
}

func (err constError) Error() string {
	return string(err)
}

func (err constError) Is(target error) bool {
	if targetErr, ok := target.(constError); ok {
		return string(err) == string(targetErr)
	}

	ts := target.Error()
	es := string(err)

	return ts == es || strings.HasPrefix(ts, es+": ")
}

func ErrNotFound(entity string) error {
	return NewConstError(fmt.Sprintf(notFoundFormat, entity))
}

func ErrInvalid(entity string) error {
	return NewConstError(fmt.Sprintf(invalidFormat, entity))
}

func ErrRequired(entity string) error {
	return NewConstError(fmt.Sprintf(missingFormat, entity))
}
