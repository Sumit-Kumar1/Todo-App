package service

import (
	"strconv"
	"strings"
	"time"
	"todoapp/models"

	"golang.org/x/exp/rand"
)

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

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
	trimmedID := strings.TrimSpace(id)

	if trimmedID == "" || strings.ContainsAny(trimmedID, "1234567890") || len(trimmedID) > 5 {
		return models.ErrInvalidID
	}

	return nil
}
