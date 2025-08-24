package todosvc

import (
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todosvc
type TodoStorer interface {
	GetAll(ctx *gofr.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx *gofr.Context, task *models.Task) error
	Update(ctx *gofr.Context, task *models.Task) error
	Delete(ctx *gofr.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx *gofr.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
