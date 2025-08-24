package models

import (
	"strings"
	"testing"
	"time"
	"todoapp/internal/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	testFail = "Test[%d] failed - %s"
)

func TestTask_ToTaskResp(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want *TaskResp
	}{
		// TODO: Add test cases.
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.task.ToTaskResp(), testFail, i, tt.name)
		})
	}
}

func TestValidateTask(t *testing.T) {
	uid := uuid.NewString()
	date := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)
	tests := []struct {
		name    string
		task    TaskReq
		wantErr error
	}{
		{name: "valid case", task: TaskReq{ID: "task-" + uid, Title: "test", DueDate: date}, wantErr: nil},
		{name: "invalid ID", task: TaskReq{ID: "123", Title: "test", DueDate: date}, wantErr: errors.ErrInvalid("task id")},
		{name: "empty title", task: TaskReq{ID: "task-" + uid, Title: ""}, wantErr: errors.ErrRequired("task title")},
		{name: "description is too large", task: TaskReq{
			ID: "task-" + uid, Title: "Hello world", Description: strings.Repeat("hello", 201), DueDate: date,
		}, wantErr: errors.ErrInvalid("task description, size > 1K characters")},
		{name: "missing dueDate", task: TaskReq{ID: "task-" + uid, Title: "test", DueDate: "  "},
			wantErr: errors.ErrRequired("due date")},
		{name: "invalid dueDate", task: TaskReq{ID: "task-" + uid, Title: "test", DueDate: " 1235 "},
			wantErr: errors.ErrInvalid("due date")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.wantErr, tt.task.Validate(), testFail, i, tt.name)
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
		{name: "valid ID", id: uid, wantErr: nil},
		{name: "invalid ID", id: "123", wantErr: errors.ErrInvalid("task id")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)

			assert.Equalf(t, tt.wantErr, err, testFail, i, tt.name)
		})
	}
}

func TestValidateDueDate(t *testing.T) {
	today := time.Now().AddDate(0, 0, 0).Format(time.DateOnly)
	dd := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)
	prevDay := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)

	tests := []struct {
		name    string
		val     string
		wantErr error
	}{
		{name: "tomorrow case", val: dd, wantErr: nil},
		{name: "today case", val: today, wantErr: nil},
		{name: "2 months ahead case", val: time.Now().AddDate(0, 2, 0).Format(time.DateOnly), wantErr: nil},
		{name: "yesterday case", val: prevDay, wantErr: errors.NewConstError("older due date from today")},
		{name: "2 months back case", val: time.Now().AddDate(0, -2, 0).Format(time.DateOnly),
			wantErr: errors.NewConstError("older due date from today")},
		{name: "empty due date", val: "", wantErr: errors.ErrRequired("due date")},
		{name: "invalid due date", val: "123", wantErr: errors.ErrInvalid("due date")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDueDate(tt.val)
			assert.Equalf(t, tt.wantErr, err, testFail, i, tt.name)
		})
	}
}
