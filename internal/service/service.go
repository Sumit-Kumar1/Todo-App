package service

import (
	"context"
	"log/slog"
	"strings"
	"todoapp/internal/models"
)

type Service struct {
	Store Storer
	Log   *slog.Logger
}

func New(st Storer, logger *slog.Logger) *Service {
	return &Service{Store: st,
		Log: logger}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.Store.GetAll(ctx)
	if err != nil {
		s.Log.Error("error in getAll", "error", err.Error())

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx context.Context, title string) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalidTitle
	}

	id := generateID()

	task, err := s.Store.Create(ctx, id, title)
	if err != nil {
		s.Log.Error("error in add task", "error", err.Error())
		return nil, err
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		s.Log.Debug("", "error", err.Error(), "ID", id)
		return models.ErrInvalidID
	}

	if err := s.Store.Delete(ctx, id); err != nil {
		s.Log.Error("error in delete task", "error", err.Error())

		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx context.Context, id string) (*models.Task, error) {
	if err := validateID(id); err != nil {
		s.Log.Debug("", "error", err.Error(), "ID", id)
		return nil, models.ErrInvalidID
	}

	task, err := s.Store.MarkDone(ctx, id)
	if err != nil {
		s.Log.Error("error in mark-done", "error", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, id, title, isDone string) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		s.Log.Debug("error while validating task", "error", err.Error(), "ID", id)
		return nil, err
	}

	task, err := s.Store.Update(ctx, id, title)
	if err != nil {
		s.Log.Error("error in updating task", "error", err.Error())

		return nil, err
	}

	s.Log.Info("\nUpdated task", "task", task)

	return task, nil
}
