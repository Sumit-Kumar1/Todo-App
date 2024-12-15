package todosvc

import (
	"strings"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

const (
	idPrefix = "task-"
)

func generateID() string {
	return idPrefix + uuid.New().String()
}

func validateTask(id, title string) error {
	if err := validateID(id); err != nil {
		return err
	}

	if strings.TrimSpace(title) == "" {
		return models.ErrInvalid("task title")
	}

	return nil
}

func validateID(id string) error {
	splits := strings.Split(id, idPrefix)
	if len(splits) != 2 {
		return models.ErrInvalid("task id")
	}

	uid, err := uuid.Parse(splits[1])
	if err != nil || uid == uuid.Nil {
		return models.ErrInvalid("task id")
	}

	return nil
}
