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
}

func New(st TodoStorer) *Service {
	return &Service{Store: st}
}

func (s *Service) GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	tasks, err := s.Store.GetAll(ctx, userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while fetcing all tasks",
			slog.String("error", err.Error()), slog.String("user", userID.String()))

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(
	ctx context.Context,
	title string,
	userID *uuid.UUID,
) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)
	title = strings.TrimSpace(title)

	if title == "" {
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
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"error while creating task - store.Create",
			slog.String("error", err.Error()),
			slog.String("task", task.ID),
		)

		return nil, err
	}

	return &task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return models.ErrInvalid("task id")
	}

	if err := s.Store.Delete(ctx, id, userID); err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"error while deleting task",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return err
	}

	return nil
}

func (s *Service) MarkDone(
	ctx context.Context,
	id string,
	userID *uuid.UUID,
) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return nil, err
	}

	task, err := s.Store.MarkDone(ctx, id, userID)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"error while marking task done",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(
	ctx context.Context,
	id, title string,
	isDone bool,
	userID *uuid.UUID,
) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	title = strings.TrimSpace(title) // trimmed the space around the task title

	if err := validateTask(id, title); err != nil {
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
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"error while updating task",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return nil, err
	}

	return &task, nil
}
