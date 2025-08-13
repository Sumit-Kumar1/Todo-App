package todosvc

import (
	"strings"
	"time"
	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	prefixTask = "task-"
)

func generateID() string {
	return prefixTask + uuid.New().String()
}

func validateTask(id string, task *models.TaskReq) error {
	if err := validateID(id); err != nil {
		return err
	}

	task.Title = strings.TrimSpace(task.Title)
	task.Description = strings.TrimSpace(task.Description)

	if task.Title == "" {
		return errors.ErrRequired("task title")
	}

	if len(task.Description) > 1000 {
		return errors.ErrInvalid("task description, size > 1K characters")
	}

	return validateDueDate(task.DueDate)
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
	splits := strings.Split(id, prefixTask)
	if len(splits) != 2 {
		return errors.ErrInvalid("task id")
	}

	uid, err := uuid.Parse(splits[1])
	if err != nil || uid == uuid.Nil {
		return errors.ErrInvalid("task id")
	}

	return nil
}
