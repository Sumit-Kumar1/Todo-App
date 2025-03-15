package todosvc

import (
	"strings"
	"testing"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

func TestGenerateID(t *testing.T) {
	uid := uuid.NewString()
	tests := []struct {
		name string
		want string
	}{
		{name: "valid ID", want: prefixTask + uid},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateID()

			if !strings.Contains(got, prefixTask) || len(got) != len(uid)+5 {
				t.Errorf("Test[%d] Failed - %s\nGot:\t%+v\nWant:\t%+v", i, tt.name, got, tt.want)
			}
		})
	}
}

func TestValidateTask(t *testing.T) {
	uid := uuid.NewString()
	tests := []struct {
		name    string
		id      string
		title   string
		wantErr error
	}{
		{name: "valid case", id: "task-" + uid, title: "test", wantErr: nil},
		{name: "invalid ID", id: "123", title: "test", wantErr: models.ErrInvalid("task id")},
		{name: "empty title", id: "task-" + uid, title: "", wantErr: models.ErrInvalid("task title")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTask(tt.id, tt.title); err != tt.wantErr {
				t.Errorf("validateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	uid := uuid.NewString()
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{name: "valid ID", id: "task-" + uid, wantErr: nil},
		{name: "invalid ID", id: "123", wantErr: models.ErrInvalid("task id")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateID(tt.id); err != tt.wantErr {
				t.Errorf("validateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
