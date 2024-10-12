package todosvc

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"todoapp/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	Store TodoStorer
	Log   *slog.Logger
}

func New(st TodoStorer, logger *slog.Logger) *Service {
	return &Service{Store: st, Log: logger}
}

func (s *Service) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	tasks, err := s.Store.GetAll(ctx, userID)
	if err != nil {
		s.Log.Error("error in getAll", "error", err.Error())

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx context.Context, title string, userID *uuid.UUID) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalid("task title")
	}

	task := models.Task{
		ID:      generateID(),
		Title:   title,
		UserID:  *userID,
		IsDone:  false,
		AddedAt: time.Now().UTC(),
	}

	err := s.Store.Create(ctx, &task)
	if err != nil {
		s.Log.Error("error in add task", "error", err.Error())
		return nil, err
	}

	return &task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error {
	if err := validateID(id); err != nil {
		s.Log.Debug("", "error", err.Error(), "ID", id)
		return models.ErrInvalid("task id")
	}

	if err := s.Store.Delete(ctx, id, userID); err != nil {
		s.Log.Error("error in delete task", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	if err := validateID(id); err != nil {
		s.Log.ErrorContext(ctx, err.Error(), "task id", id)
		return nil, err
	}

	task, err := s.Store.MarkDone(ctx, id, userID)
	if err != nil {
		s.Log.Error("error in mark-done", "error", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, id, title string, isDone bool, userID *uuid.UUID) (*models.Task, error) {
	if err := validateTask(id, title); err != nil {
		s.Log.ErrorContext(ctx, err.Error(), "ID", id)
		return nil, err
	}

	mt := time.Now().UTC()

	task := models.Task{
		ID:         id,
		Title:      title,
		UserID:     *userID,
		IsDone:     isDone,
		ModifiedAt: &mt,
	}

	err := s.Store.Update(ctx, &task)
	if err != nil {
		s.Log.Error("error in updating task", "error", err.Error())

		return nil, err
	}

	s.Log.Info("\nUpdated task", "task: %+v", task)

	return &task, nil
}
