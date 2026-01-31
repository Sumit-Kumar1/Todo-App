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
	taskID := uuid.NewString()
	userID := uuid.New()
	title := "random task title"
	description := "how to do a task clear"
	addedAt := time.Now()
	modifiedAt := time.Now()
	dueDate := time.Now()
	dateOnly := dueDate.Format(time.DateOnly)

	tests := []struct {
		name string
		task Task
		want *TaskResp
	}{
		{name: "valid case", task: Task{ID: taskID, UserID: userID, Title: title, Description: description, Status: false,
			DueDate: &dueDate, CreatedAt: addedAt, UpdatedAt: &modifiedAt},
			want: &TaskResp{ID: taskID, UserID: userID, Title: title, Description: description, DueDate: &dateOnly, AddedAt: addedAt, ModifiedAt: &modifiedAt}},
		{name: "no due date", task: Task{ID: taskID, UserID: userID, Title: title, Description: description, IsDone: false,
			DueDate: nil, CreatedAt: addedAt, UpdatedAt: &modifiedAt},
			want: &TaskResp{ID: taskID, UserID: userID, Title: title, Description: description, DueDate: nil, AddedAt: addedAt, ModifiedAt: &modifiedAt}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.task.ToTaskResp(), testFail, i, tt.name)
		})
	}
}

func TestValidateTask(t *testing.T) {
	date := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)
	tests := []struct {
		name    string
		task    TaskReq
		wantErr error
	}{
		{name: "valid case", task: TaskReq{Title: "test", DueDate: date}, wantErr: nil},
		{name: "empty title", task: TaskReq{Title: ""}, wantErr: errors.ErrRequired("task title")},
		{name: "description is too large", task: TaskReq{Title: "Hello world", Description: strings.Repeat("hello", 201),
			DueDate: date}, wantErr: errors.ErrInvalid("task description, size > 1K characters")},
		{name: "missing dueDate", task: TaskReq{Title: "test", DueDate: "  "},
			wantErr: nil},
		{name: "invalid dueDate", task: TaskReq{Title: "test", DueDate: " 1235 "},
			wantErr: errors.ErrInvalid("due date")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.wantErr, tt.task.Validate(), testFail, i, tt.name)
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
		{name: "empty due date", val: "", wantErr: nil},
		{name: "invalid due date", val: "123", wantErr: errors.ErrInvalid("due date")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDueDate(tt.val)
			assert.Equalf(t, tt.wantErr, err, testFail, i, tt.name)
		})
	}
}
