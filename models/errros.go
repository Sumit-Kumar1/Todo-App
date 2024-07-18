package models

import (
	"fmt"
	"strings"
)

// Errors raised by package models.
var (
	ErrPermissionDenied = constError("permission denied")
	ErrInvalidId        = constError("invalid id")
	ErrNotFound         = constError("not found")
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
func (err constError) wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

type wrapError struct {
	err error
	msg string
}

func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}
func (err wrapError) Unwrap() error {
	return err.err
}
func (err wrapError) Is(target error) bool {
	return constError(err.msg).Is(target)
}
