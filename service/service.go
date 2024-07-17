package service

import (
	"fmt"
	"log"
	"strings"
	"todoapp/models"

	"github.com/google/uuid"
)

type Service struct {
	Data map[uuid.UUID]models.Task
}

func New() *Service {
	d := make(map[uuid.UUID]models.Task)
	return &Service{Data: d}
}

func (s *Service) AddTask(title, desc string) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("Task is empty")
	}

	var (
		id = uuid.New()
		t  = models.Task{
			ID:     id,
			Title:  title,
			IsDone: false,
		}
	)

	if strings.TrimSpace(desc) == "" {
		t.Desc = "<n/a>"
	}

	s.Data[id] = t

	log.Println("Task Added ID: ", id)

	return &t, nil
}

func (s *Service) DeleteTask(id string) error {

	return nil
}

func (s *Service) MarkDone(id string) error {
	return nil
}
