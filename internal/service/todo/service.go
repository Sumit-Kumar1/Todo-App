package todoservice

import (
	"strings"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type Service struct {
	Store TodoStorer
}

func New(ts TodoStorer) *Service {
	return &Service{Store: ts}
}

func (s *Service) GetAll(ctx server.Context, userID *uuid.UUID) ([]models.Task, error) {
	tasks, err := s.Store.GetAll(ctx, userID)
	if err != nil {
		ctx.Logger.Error("error in getAll", "error", err.Error())

		return nil, err
	}

	return tasks, nil
}

func (s *Service) AddTask(ctx server.Context, title string, userID *uuid.UUID) (*models.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, models.ErrInvalidTitle
	}

	id := generateID()

	task, err := s.Store.Create(ctx, id, title, userID)
	if err != nil {
		ctx.Logger.Error("error in add task", "error", err.Error())
		return nil, err
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx server.Context, id string, userID *uuid.UUID) error {
	if err := validateID(id); err != nil {
		ctx.Logger.Debug("", "error", err.Error(), "ID", id)
		return models.ErrInvalidID
	}

	if err := s.Store.Delete(ctx, id, userID); err != nil {
		ctx.Logger.Error("error in delete task", "error", err.Error())
		return err
	}

	return nil
}

func (s *Service) MarkDone(ctx server.Context, id string, userID *uuid.UUID) (*models.Task, error) {
	if err := validateID(id); err != nil {
		ctx.Logger.Debug("", "error", err.Error(), "ID", id)
		return nil, models.ErrInvalidID
	}

	task, err := s.Store.MarkDone(ctx, id, userID)
	if err != nil {
		ctx.Logger.Error("error in mark-done", "error", err.Error())

		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx server.Context, id, title, isDone string, userID *uuid.UUID) (*models.Task, error) {
	if err := validateTask(id, title, isDone); err != nil {
		ctx.Logger.Debug("error while validating task", "error", err.Error(), "ID", id)
		return nil, err
	}

	task, err := s.Store.Update(ctx, id, title, userID)
	if err != nil {
		ctx.Logger.Error("error in updating task", "error", err.Error())

		return nil, err
	}

	ctx.Logger.Info("\nUpdated task", "task", task)

	return task, nil
}
