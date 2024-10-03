package service

import (
	"errors"
	"testing"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	uuid.DisableRandPool()
	defer uuid.EnableRandPool()

	id := uuid.NewString()

	tests := []struct {
		wantErr error
		name    string
		id      string
	}{
		{name: "valid case", id: "css-" + id, wantErr: nil},
		{name: "nil case", id: "", wantErr: models.ErrInvalid("task id")},
		{name: "nil case", id: "css-" + uuid.Nil.String(), wantErr: models.ErrInvalid("task id")},
		{name: "invalid length", id: "abceox", wantErr: models.ErrInvalid("task id")},
		{name: "invalid case", id: "abc12", wantErr: models.ErrInvalid("task id")},
		{name: "invalid case", id: "12345", wantErr: models.ErrInvalid("task id")},
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
	uuid.DisableRandPool()
	defer uuid.EnableRandPool()

	id := uuid.NewString()

	tests := []struct {
		wantErr error
		name    string
		expID   string
	}{
		{name: "valid case", wantErr: nil, expID: "css-" + id},
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
	uuid.DisableRandPool()
	defer uuid.EnableRandPool()

	cssID := "css-" + uuid.NewString()

	tests := []struct {
		wantErr error
		name    string
		id      string
		title   string
		isDone  string
	}{
		{name: "valid case", id: cssID, title: "hello", isDone: "true", wantErr: nil},
		{name: "invalid ID", id: "css-" + uuid.Nil.String(), title: "hello", isDone: "true", wantErr: models.ErrInvalid("task id")},
		{name: "invalid Title", id: cssID, title: "", isDone: "true", wantErr: models.ErrInvalid("task title")},
		{name: "invalid done", id: cssID, title: "hello world", isDone: "not known", wantErr: models.ErrInvalid("task done")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTask(tt.id, tt.title, tt.isDone); !errors.Is(err, tt.wantErr) {
				t.Errorf("\nTEST[%d] Failed - %s\n\tExpected: %+v\n\tActual: %+v", i, tt.name, tt.wantErr, err)
			}
		})
	}
}

func Test_encryptedPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"default case", "hello word", nil},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptedPassword(tt.password)

			assert.Equal(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}
