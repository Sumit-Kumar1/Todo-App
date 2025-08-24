package models

import (
	"strings"
	"time"
	"todoapp/internal/errors"

	"github.com/google/uuid"
)

type Task struct {
	ID          string     `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsDone      bool       `json:"isDone"`
	DueDate     *time.Time `json:"dueDate"`
	AddedAt     time.Time  `json:"addedAt"`
	ModifiedAt  *time.Time `json:"modifiedAt"`
}

type TaskResp struct {
	ID          string     `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsDone      bool       `json:"isDone"`
	DueDate     *string    `json:"dueDate"`
	AddedAt     time.Time  `json:"addedAt"`
	ModifiedAt  *time.Time `json:"modifiedAt"`
}

type TaskReq struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	IsDone      bool   `json:"isDone"`
}

func (t *Task) ToTaskResp() *TaskResp {
	tr := TaskResp{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		IsDone:      t.IsDone,
		AddedAt:     t.AddedAt,
		ModifiedAt:  t.ModifiedAt,
	}

	dd := t.DueDate.Format(time.DateOnly)
	tr.DueDate = &dd

	return &tr
}

func (t *TaskReq) Validate() error {
	if err := validateID(t.ID); err != nil {
		return err
	}

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
		return errors.ErrRequired("due date")
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

func validateID(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil || uid == uuid.Nil {
		return errors.ErrInvalid("task id")
	}

	return nil
}
