package todohttp

import (
	"todoapp/internal/models"

	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=todohttp
type TodoServicer interface {
	GetAll(ctx *gofr.Context, userID *uuid.UUID) ([]models.Task, error)
	AddTask(ctx *gofr.Context, task *models.TaskReq, userID *uuid.UUID) (*models.Task, error)
	DeleteTask(ctx *gofr.Context, id string, userID *uuid.UUID) error
	UpdateTask(ctx *gofr.Context, id string, task *models.TaskReq, isDone bool, userID *uuid.UUID) (*models.Task, error)
	MarkDone(ctx *gofr.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
