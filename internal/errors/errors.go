package errors

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
	_, _ = fmt.Fprintf(w, `<div class="toast toast-top toast-end" id="notifi">
  <div class="alert alert-error"><i class="fa-solid fa-xmark" id="toastClost" onclick="removeToast()"></i>
    </svg><span>%d: %s</span></div></div>`, status, err.Error())
}
