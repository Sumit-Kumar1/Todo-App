package service

import (
	"strconv"
	"strings"
	"todoapp/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func encryptedPassword(password string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwd), nil
}

func generateID() string {
	return "css-" + uuid.New().String()
}

func validateTask(id, title, isDone string) error {
	if err := validateID(id); err != nil {
		return err
	}

	if strings.TrimSpace(title) == "" {
		return models.ErrInvalid("task title")
	}

	if _, err := strconv.ParseBool(isDone); err != nil {
		return models.ErrInvalid("task done status")
	}

	return nil
}

func validateID(id string) error {
	splits := strings.Split(id, "css-")
	if len(splits) != 2 {
		return models.ErrInvalid("task id")
	}

	uid, err := uuid.Parse(splits[1])
	if err != nil || uid == uuid.Nil {
		return models.ErrInvalid("task id")
	}

	return nil
}
