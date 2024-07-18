package service

import (
	"strconv"
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
		return nil, models.ErrTaskTitle
	}

	var (
		id = uuid.New()
		t  = models.Task{
			ID:     id,
			Desc:   desc,
			Title:  title,
			IsDone: false,
		}
	)

	if strings.TrimSpace(desc) == "" {
		t.Desc = "<n/a>"
	}

	s.Data[id] = t

	return &t, nil
}

func (s *Service) DeleteTask(id string) error {
	if err := validateID(id); err != nil {
		return models.ErrInvalidId
	}

	uid, _ := uuid.Parse(id)
	if _, ok := s.Data[uid]; !ok {
		return models.ErrNotFound
	}

	delete(s.Data, uid)

	return nil
}

func (s *Service) MarkDone(id string) (*models.Task, error) {
	if err := validateID(id); err != nil {
		return nil, models.ErrInvalidId
	}

	uid, _ := uuid.Parse(id)

	task, ok := s.Data[uid]
	if !ok {
		return nil, models.ErrNotFound
	}

	task.IsDone = true
	s.Data[uid] = task

	return &task, nil
}

func (s *Service) UpdateTask(id, title, desc, isDone string) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		return nil, err
	}

	uid, _ := uuid.Parse(id)

	task, ok := s.Data[uid]
	if !ok {
		return nil, models.ErrNotFound
	}

	task.Title = title
	if strings.TrimSpace(desc) == "" {
		task.Desc = "<n/a>"
	}

	task.IsDone, _ = strconv.ParseBool(isDone)
	s.Data[uid] = task

	return &task, nil
}
