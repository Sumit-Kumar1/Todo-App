package todosvc

import (
	"context"
	"log/slog"
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

func (s *Service) AddTask(ctx context.Context, taskInp *models.TaskReq, userID *uuid.UUID) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)
	id := generateID()

	if err := validateTask(id, taskInp); err != nil {
		return nil, err
	}

	dd, _ := time.Parse(time.DateOnly, taskInp.DueDate)

	task := models.Task{
		ID:          id,
		UserID:      *userID,
		Title:       taskInp.Title,
		Description: taskInp.Description,
		IsDone:      false,
		DueDate:     &dd,
		AddedAt:     time.Now().UTC(),
	}

	if err := s.Store.Create(ctx, &task); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while creating task - store.Create",
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
		logger.LogAttrs(ctx, slog.LevelError, "error while deleting task",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return nil, err
	}

	task, err := s.Store.MarkDone(ctx, id, userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while marking task done",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, id string, taskInp *models.TaskReq, isDone bool, userID *uuid.UUID,
) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateTask(id, taskInp); err != nil {
		return nil, err
	}

	dd, _ := time.Parse(time.DateOnly, taskInp.DueDate)
	mt := time.Now().UTC()

	task := models.Task{
		ID:          id,
		UserID:      *userID,
		Title:       taskInp.Title,
		Description: taskInp.Description,
		DueDate:     &dd,
		IsDone:      isDone,
		ModifiedAt:  &mt,
	}

	err := s.Store.Update(ctx, &task)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while updating task",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return nil, err
	}

	return &task, nil
}
