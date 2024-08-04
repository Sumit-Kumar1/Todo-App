package service

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"todoapp/internal/models"
)

type Service struct {
	Data map[string]*models.Task
}

func New() *Service {
	return &Service{Data: make(map[string]*models.Task)}
}

func (s *Service) GetAll() []models.Task {
	var tasks []models.Task

	for i := range s.Data {
		if s.Data[i] == nil {
			continue
		}

		tasks = append(tasks, *s.Data[i])
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].AddedAt.Before(tasks[j].AddedAt)
	})

	return tasks
}

func (s *Service) AddTask(title string) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalidTitle
	}

	var (
		id = generateID()
		t  = models.Task{
			ID:      id,
			Title:   title,
			IsDone:  false,
			AddedAt: time.Now(),
		}
	)

	s.Data[id] = &t

	return &t, nil
}

func (s *Service) DeleteTask(id string) error {
	if err := validateID(id); err != nil {
		return models.ErrInvalidID
	}

	if _, ok := s.Data[id]; !ok {
		return models.ErrNotFound
	}

	delete(s.Data, id)

	return nil
}

func (s *Service) MarkDone(id string) (*models.Task, error) {
	if err := validateID(id); err != nil {
		return nil, models.ErrInvalidID
	}

	task, ok := s.Data[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	task.IsDone = true
	task.ModifiedAt = time.Now()
	s.Data[id] = task

	return task, nil
}

func (s *Service) UpdateTask(id, title, isDone string) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		return nil, err
	}

	task, ok := s.Data[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	task.Title = title
	task.IsDone, _ = strconv.ParseBool(isDone)
	task.ModifiedAt = time.Now()

	s.Data[id] = task

	return task, nil
}
