package models

import (
	"strings"
	"time"
	"todoapp/internal/errors"

	"github.com/google/uuid"
)

type Task struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId,omitempty"`
	UserID      uuid.UUID  `json:"userId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      *string    `json:"status"`
	DueDate     *time.Time `json:"dueDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	ChildTasks  []Task     `json:"childTasks,omitempty"`
}

type TaskResp struct {
	ID          string     `json:"id,omitempty"`
	UserID      uuid.UUID  `json:"-"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	DueDate     *string    `json:"dueDate,omitempty"`
	ChildTasks  []TaskResp `json:"childTasks,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type TaskReq struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DueDate     string  `json:"dueDate"`
	Status      string  `json:"status"`
	ParentID    *string `json:"parentId,omitempty"`
}

func (t *Task) ToTaskResp() *TaskResp {
	tr := TaskResp{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      *t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		ChildTasks:  make([]TaskResp, 0),
	}

	if t.DueDate != nil {
		dd := t.DueDate.Format(time.DateOnly)
		tr.DueDate = &dd
	}

	for _, child := range t.ChildTasks {
		tr.ChildTasks = append(tr.ChildTasks, *child.ToTaskResp())
	}

	return &tr
}

func (t *TaskReq) Validate() error {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	if t.Title == "" {
		return errors.ErrRequired("task title")
	}

	if len(t.Description) > 1000 {
		return errors.ErrInvalid("task description, size > 1K characters")
	}

	return validateDueDate(t.DueDate)
}

func validateDueDate(val string) error {
	tn := time.Now().AddDate(0, 0, -1)

	if strings.TrimSpace(val) == "" {
		return nil
	}

	tt, err := time.Parse(time.DateOnly, val)
	if err != nil {
		return errors.ErrInvalid("due date")
	}

	if !tt.After(tn) {
		return errors.NewConstError("older due date from today")
	}

	return nil
}
