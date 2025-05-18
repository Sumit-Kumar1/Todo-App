package todosvc

import (
	"errors"
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
		task    models.TaskReq
		wantErr error
	}{
		{name: "valid case", task: models.TaskReq{ID: "task-" + uid, Title: "test"}, wantErr: nil},
		{name: "invalid ID", task: models.TaskReq{ID: "123", Title: "test"}, wantErr: models.ErrInvalid("task id")},
		{name: "empty title", task: models.TaskReq{ID: "task-" + uid, Title: ""}, wantErr: models.ErrInvalid("task title")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTask(tt.task.ID, &tt.task); !errors.Is(err, tt.wantErr) {
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
			if err := validateID(tt.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("validateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
