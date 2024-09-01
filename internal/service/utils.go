package service

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"todoapp/internal/models"
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
		return models.ErrInvalidTitle
	}

	if _, err := strconv.ParseBool(isDone); err != nil {
		return models.ErrInvalidDoneStatus
	}

	return nil
}

func validateID(id string) error {
	splits := strings.Split(id, "css-")
	if len(splits) != 2 {
		return models.ErrInvalidID
	}

	uid, err := uuid.Parse(splits[1])
	if err != nil || uid == uuid.Nil {
		return models.ErrInvalidID
	}

	return nil
}
