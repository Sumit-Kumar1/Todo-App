package models

import (
	"errors"
	"fmt"
	"net/http"
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

func HandleHTTPError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<div class="toast toast-top toast-end">
  <div class="alert alert-error"><svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
        </svg><span>%s</span></div></div>`, err.Error())
}
