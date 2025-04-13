package todosvc

import (
	"context"

	"todoapp/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todosvc
type TodoStorer interface {
	GetAll(ctx context.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx context.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
