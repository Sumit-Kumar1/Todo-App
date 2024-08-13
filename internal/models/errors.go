package models

import (
	"strings"
)

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
