package todoservice

import (
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/google/uuid"
)

type TodoStorer interface {
	GetAll(ctx server.Context, userID *uuid.UUID) ([]models.Task, error)
	Create(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Update(ctx server.Context, id, title string, userID *uuid.UUID) (*models.Task, error)
	Delete(ctx server.Context, id string, userID *uuid.UUID) error
	MarkDone(ctx server.Context, id string, userID *uuid.UUID) (*models.Task, error)
}
