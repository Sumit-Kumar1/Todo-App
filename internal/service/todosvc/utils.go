package todosvc

import (
	"strings"
	"time"
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
		return models.ErrRequired("task title")
	}

	if len(task.Description) > 1000 {
		return models.ErrInvalid("task description, size > 1K characters")
	}

	if strings.TrimSpace(task.DueDate) == "" {
		return models.ErrRequired("due date")
	}

	if _, err := time.Parse(time.DateOnly, task.DueDate); err != nil {
		return models.ErrInvalid("due date")
	}

	return nil
}

func validateID(id string) error {
	splits := strings.Split(id, prefixTask)
	if len(splits) != 2 {
		return models.ErrInvalid("task id")
	}

	uid, err := uuid.Parse(splits[1])
	if err != nil || uid == uuid.Nil {
		return models.ErrInvalid("task id")
	}

	return nil
}
