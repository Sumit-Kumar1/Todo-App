package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_constError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  constError
		want string
	}{
		{name: "err not found", err: ErrNotFound, want: "task not found"},
		{name: "invalid id", err: ErrInvalidID, want: "invalid id"},
		{name: "invalid task title", err: ErrInvalidTitle, want: "invalid task title"},
		{name: "permission denied", err: ErrPermissionDenied, want: "permission denied"},
		{name: "invalid task done", err: ErrInvalidDoneStatus, want: "invalid task done"},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()

			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_constError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    constError
		target error
		want   bool
	}{
		{name: "valid case", err: ErrNotFound, target: constError("task not found"), want: true},
		{name: "invalid case", err: ErrNotFound, target: constError("invalid id"), want: false},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Is(tt.target)

			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}
