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
		return nil, fmt.Errorf("task is empty")
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
	if err := validateID(id); err != nil {
		log.Println("id error : ", err.Error())

		return err
	}

	uid, _ := uuid.Parse(id)

	_, ok := s.Data[uid]
	if !ok {
		return fmt.Errorf("not found")
	}

	delete(s.Data, uid)

	log.Print("Deleted : ", id)

	return nil
}

func (s *Service) MarkDone(id string) error {
	if err := validateID(id); err != nil {
		log.Println("id error : ", err.Error())

		return err
	}

	uid, _ := uuid.Parse(id)

	t, ok := s.Data[uid]
	if !ok {
		return fmt.Errorf("not found")
	}

	t.IsDone = true
	s.Data[uid] = t

	log.Print("Done Task: ", id)
	return nil
}

func validateID(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if uid == uuid.Nil {
		return fmt.Errorf("nil uuid")
	}

	return nil
}
