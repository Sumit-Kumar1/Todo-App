package todohttp

import (
	"context"

	"todoapp/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todohttp
type TodoServicer interface {
	GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error)
	AddTask(ctx context.Context, task string, userID *uuid.UUID) (*models.Task, error)
	DeleteTask(ctx context.Context, id string, userID *uuid.UUID) error
	UpdateTask(
		ctx context.Context,
		id, title string,
		isDone bool,
		userID *uuid.UUID,
	) (*models.Task, error)
	MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
