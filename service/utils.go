package service

import (
	"strconv"
	"strings"
	"todoapp/models"

	"github.com/google/uuid"
)

func validateTask(id, title, isDone string) error {
	if err := validateID(id); err != nil {
		return err
	}

	if strings.TrimSpace(title) == "" {
		return models.ErrTaskTitle
	}

	if _, err := strconv.ParseBool(isDone); err != nil {
		return models.ErrTaskDone
	}

	return nil
}

func validateID(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return models.ErrInvalidId
	}

	if uid == uuid.Nil {
		return models.ErrInvalidId
	}

	return nil
}
