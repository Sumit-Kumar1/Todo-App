package models

import (
	"fmt"
	"strings"
)

// Errors raised by package models.
var (
	ErrNotFound          = constError("not found")
	ErrInvalidID         = constError("invalid id")
	ErrInvalidTitle      = constError("invalid task title")
	ErrPermissionDenied  = constError("permission denied")
	ErrInvalidDoneStatus = constError("invalid task done")
)

type constError string

func (err constError) Error() string {
	return string(err)
}

func (err constError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}

//nolint:unused // will be using it
func (err constError) wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

//nolint:unused // will be using it
type wrapError struct {
	err error
	msg string
}

//nolint:unused // will be using it
func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}

//nolint:unused // will be using it
func (err wrapError) Unwrap() error {
	return err.err
}

//nolint:unused // will be using it
func (err wrapError) Is(target error) bool {
	return constError(err.msg).Is(target)
}
