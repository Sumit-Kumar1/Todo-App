package todosvc

import (
	"log/slog"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	Store TodoStorer
}

func New(st TodoStorer) *Service {
	return &Service{Store: st}
}

func (s *Service) GetAll(ctx *gin.Context, userID *uuid.UUID) ([]models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(userID.String()); err != nil {
		return nil, err
	}

	tasks, err := s.Store.GetAll(ctx, userID)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while fetching all tasks",
			slog.String("error", err.Error()), slog.String("user", userID.String()))

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx *gin.Context, taskInp *models.TaskReq, userID *uuid.UUID) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := taskInp.Validate(); err != nil {
		return nil, err
	}

	dd, _ := time.Parse(time.DateOnly, taskInp.DueDate)
	taskID := uuid.NewString()

	task := models.Task{
		ID:          taskID,
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

	return s.Store.GetTaskByID(ctx, taskID, userID)
}

func (s *Service) DeleteTask(ctx *gin.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return err
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

func (s *Service) MarkDone(ctx *gin.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := s.Store.MarkDone(ctx, id, userID); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while marking task done",
			slog.String("error", err.Error()),
			slog.String("task", id),
		)

		return nil, err
	}

	return s.Store.GetTaskByID(ctx, id, userID)
}

func (s *Service) UpdateTask(ctx *gin.Context, id string, taskInp *models.TaskReq, userID *uuid.UUID) (*models.Task, error) {
	logger := models.GetLoggerFromCtx(ctx)

	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := taskInp.Validate(); err != nil {
		return nil, err
	}

	dd, _ := time.Parse(time.DateOnly, taskInp.DueDate)
	mt := time.Now().UTC()
	task := models.Task{
		ID:          id,
		UserID:      *userID,
		Title:       taskInp.Title,
		Description: taskInp.Description,
		IsDone:      taskInp.IsDone,
		DueDate:     &dd,
		ModifiedAt:  &mt,
	}

	err := s.Store.Update(ctx, &task)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while updating task",
			slog.String("error", err.Error()),
			slog.String("task_id", id),
		)

		return nil, err
	}

	return s.Store.GetTaskByID(ctx, id, userID)
}

func validateID(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil || uid == uuid.Nil {
		return errors.ErrInvalidTaskID
	}

	return nil
}
