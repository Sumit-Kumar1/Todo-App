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

	// Use a future due date so DueWarning is empty
	dueDate := time.Now().AddDate(0, 1, 0)
	dateOnly := dueDate.Format(time.DateOnly)
	status := "ACTIVE"

	tests := []struct {
		name string
		task Task
		want *TaskResp
	}{
		{
			name: "valid case with priority and category",
			task: Task{
				ID: taskID, UserID: userID, Title: title, Description: description,
				Status: &status, Priority: "HIGH", Category: "Work",
				DueDate: &dueDate, CreatedAt: addedAt, UpdatedAt: &modifiedAt,
			},
			want: &TaskResp{
				ID: taskID, UserID: userID, Title: title, Description: description,
				Status: status, Priority: "HIGH", Category: "Work",
				DueDate: &dateOnly, ChildTasks: []TaskResp{},
				CreatedAt: addedAt, UpdatedAt: &modifiedAt,
			},
		},
		{
			name: "no due date",
			task: Task{
				ID: taskID, UserID: userID, Title: title, Description: description,
				Status: &status, Priority: "MEDIUM", Category: "",
				DueDate: nil, CreatedAt: addedAt, UpdatedAt: &modifiedAt,
			},
			want: &TaskResp{
				ID: taskID, UserID: userID, Title: title, Description: description,
				Status: status, Priority: "MEDIUM",
				DueDate: nil, ChildTasks: []TaskResp{},
				CreatedAt: addedAt, UpdatedAt: &modifiedAt,
			},
		},
		{
			name: "overdue task",
			task: func() Task {
				past := time.Now().AddDate(0, 0, -3)
				return Task{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: &status, Priority: "URGENT", Category: "Personal",
					DueDate: &past, CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
			want: func() *TaskResp {
				past := time.Now().AddDate(0, 0, -3)
				pastStr := past.Format(time.DateOnly)
				return &TaskResp{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: status, Priority: "URGENT", Category: "Personal",
					DueDate: &pastStr, DueWarning: "overdue", ChildTasks: []TaskResp{},
					CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
		},
		{
			name: "due today task",
			task: func() Task {
				today := time.Now()
				return Task{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: &status, Priority: "HIGH", Category: "",
					DueDate: &today, CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
			want: func() *TaskResp {
				today := time.Now()
				todayStr := today.Format(time.DateOnly)
				return &TaskResp{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: status, Priority: "HIGH",
					DueDate: &todayStr, DueWarning: "due_today", ChildTasks: []TaskResp{},
					CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
		},
		{
			name: "due soon task (2 days)",
			task: func() Task {
				soon := time.Now().AddDate(0, 0, 2)
				return Task{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: &status, Priority: "MEDIUM", Category: "",
					DueDate: &soon, CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
			want: func() *TaskResp {
				soon := time.Now().AddDate(0, 0, 2)
				soonStr := soon.Format(time.DateOnly)
				return &TaskResp{
					ID: taskID, UserID: userID, Title: title, Description: description,
					Status: status, Priority: "MEDIUM",
					DueDate: &soonStr, DueWarning: "due_soon", ChildTasks: []TaskResp{},
					CreatedAt: addedAt, UpdatedAt: &modifiedAt,
				}
			}(),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.ToTaskResp()
			assert.Equalf(t, tt.want, got, testFail, i, tt.name)
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
		{name: "valid with priority", task: TaskReq{Title: "test", DueDate: date, Priority: "HIGH"}, wantErr: nil},
		{name: "empty priority defaults to MEDIUM", task: TaskReq{Title: "test", DueDate: date, Priority: ""}, wantErr: nil},
		{name: "invalid priority", task: TaskReq{Title: "test", DueDate: date, Priority: "CRITICAL"},
			wantErr: errors.Invalid("priority, must be one of: LOW, MEDIUM, HIGH, URGENT")},
		{name: "valid with category", task: TaskReq{Title: "test", DueDate: date, Category: "Work"}, wantErr: nil},
		{name: "empty title", task: TaskReq{Title: ""}, wantErr: errors.Required("task title")},
		{name: "description is too large", task: TaskReq{Title: "Hello world", Description: strings.Repeat("hello", 201),
			DueDate: date}, wantErr: errors.Invalid("task description, size > 1K characters")},
		{name: "missing dueDate", task: TaskReq{Title: "test", DueDate: "  "},
			wantErr: nil},
		{name: "invalid dueDate", task: TaskReq{Title: "test", DueDate: " 1235 "},
			wantErr: errors.Invalid("due date")},
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
		{name: "yesterday case", val: prevDay, wantErr: errors.ErrDueDatePast},
		{name: "2 months back case", val: time.Now().AddDate(0, -2, 0).Format(time.DateOnly),
			wantErr: errors.ErrDueDatePast},
		{name: "empty due date", val: "", wantErr: nil},
		{name: "invalid due date", val: "123", wantErr: errors.Invalid("due date")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDueDate(tt.val)
			assert.Equalf(t, tt.wantErr, err, testFail, i, tt.name)
		})
	}
}

func TestValidatePriority(t *testing.T) {
	date := time.Now().AddDate(0, 0, 1).Format(time.DateOnly)

	tests := []struct {
		name         string
		priority     string
		wantPriority string
		wantErr      bool
	}{
		{name: "LOW is valid", priority: "LOW", wantPriority: "LOW", wantErr: false},
		{name: "MEDIUM is valid", priority: "MEDIUM", wantPriority: "MEDIUM", wantErr: false},
		{name: "HIGH is valid", priority: "HIGH", wantPriority: "HIGH", wantErr: false},
		{name: "URGENT is valid", priority: "URGENT", wantPriority: "URGENT", wantErr: false},
		{name: "lowercase converted to upper", priority: "high", wantPriority: "HIGH", wantErr: false},
		{name: "empty defaults to MEDIUM", priority: "", wantPriority: "MEDIUM", wantErr: false},
		{name: "CRITICAL is invalid", priority: "CRITICAL", wantPriority: "CRITICAL", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := TaskReq{Title: "test", DueDate: date, Priority: tt.priority}
			err := req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPriority, req.Priority)
			}
		})
	}
}

func testPtr[T any](t T) *T {
	return &t
}

func TestComputeDueWarning(t *testing.T) {
	tests := []struct {
		name    string
		dueDate *time.Time
		want    string
	}{
		{name: "nil due date", dueDate: nil, want: ""},
		{name: "far future", dueDate: testPtr(time.Now().AddDate(0, 1, 0)), want: ""},
		{name: "overdue", dueDate: testPtr(time.Now().AddDate(0, 0, -2)), want: "overdue"},
		{name: "due today", dueDate: testPtr(time.Now()), want: "due_today"},
		{name: "due in 1 day", dueDate: testPtr(time.Now().AddDate(0, 0, 1)), want: "due_soon"},
		{name: "due in 3 days", dueDate: testPtr(time.Now().AddDate(0, 0, 3)), want: "due_soon"},
		{name: "due in 4 days", dueDate: testPtr(time.Now().AddDate(0, 0, 4)), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDueWarning(tt.dueDate)
			assert.Equal(t, tt.want, got)
		})
	}
}
