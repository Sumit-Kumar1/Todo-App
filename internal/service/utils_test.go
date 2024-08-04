package service

import (
	"errors"
	"testing"
	"todoapp/internal/models"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		wantErr error
		name    string
		id      string
	}{
		{name: "valid case", id: "abceo", wantErr: nil},
		{name: "nil case", id: "", wantErr: models.ErrInvalidID},
		{name: "invalid length", id: "abceox", wantErr: models.ErrInvalidID},
		{name: "invalid case", id: "abc12", wantErr: models.ErrInvalidID},
		{name: "invalid case", id: "12345", wantErr: models.ErrInvalidID},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateID(tt.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("\nTEST[%d] Failed - %s\n\tExpected: %+v\n\tActual: %+v", i, tt.name, tt.wantErr, err)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		wantErr error
		name    string
	}{
		{name: "valid case", wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateID()
			if err := validateID(got); !errors.Is(err, tt.wantErr) {
				t.Errorf("generateID() = %v", got)
			}
		})
	}
}

func TestValidateTask(t *testing.T) {
	tests := []struct {
		wantErr error
		name    string
		id      string
		title   string
		isDone  string
	}{
		{name: "valid case", id: "abcde", title: "hello", isDone: "true", wantErr: nil},
		{name: "invalid ID", id: "abcd1", title: "hello", isDone: "true", wantErr: models.ErrInvalidID},
		{name: "invalid Title", id: "zAcdx", title: "", isDone: "true", wantErr: models.ErrInvalidTitle},
		{name: "invalid done", id: "pAUbe", title: "hello world", isDone: "not known", wantErr: models.ErrInvalidDoneStatus},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTask(tt.id, tt.title, tt.isDone); !errors.Is(err, tt.wantErr) {
				t.Errorf("\nTEST[%d] Failed - %s\n\tExpected: %+v\n\tActual: %+v", i, tt.name, tt.wantErr, err)
			}
		})
	}
}
